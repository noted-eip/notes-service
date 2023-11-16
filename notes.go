package main

import (
	"context"
	"time"

	"notes-service/auth"

	background "github.com/noted-eip/noted/background-service"

	"notes-service/exports"
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"notes-service/validators"

	"notes-service/language"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type notesAPI struct {
	notesv1.UnimplementedNotesAPIServer

	logger *zap.Logger

	auth       auth.Service
	language   language.Service
	background background.Service

	notes      models.NotesRepository
	groups     models.GroupsRepository
	activities models.ActivitiesRepository
}

var _ notesv1.NotesAPIServer = &notesAPI{}

func (srv *notesAPI) CreateNote(ctx context.Context, req *notesv1.CreateNoteRequest) (*notesv1.CreateNoteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateCreateNoteRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check user is part of the group.
	_, err = srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	// Check user can edit the note
	note, err := srv.notes.CreateNote(ctx, &models.CreateNotePayload{
		GroupID:         req.GroupId,
		Title:           req.Title,
		AuthorAccountID: token.AccountID,
		FolderID:        "",
		Lang:            req.Lang,
		Blocks:          protobufBlocksToModelsBlocks(req.Blocks),
	}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	err = srv.background.AddProcess(&background.Process{
		Identifier: models.NoteIdentifier{Metadata: note.ID, ActionType: models.NoteUpdateKeyword},
		CallBackFct: func() error {
			err := srv.UpdateKeywordsByNoteId(note.ID, req.GroupId, token.AccountID)
			return err
		},
		SecondsToDebounce:             5,
		CancelProcessOnSameIdentifier: true,
		RepeatProcess:                 false,
	})
	if err != nil {
		return nil, err
	}

	_, err = srv.activities.CreateActivityInternal(ctx, &models.ActivityPayload{
		GroupID: note.GroupID,
		Type:    models.NoteAdded,
		Event:   "<userID:" + note.AuthorAccountID + "> has added the note <noteID:" + note.ID + "> in the folder <folderID:" + "" + ">.",
	})
	if err != nil {
		return nil, err
	}

	return &notesv1.CreateNoteResponse{Note: modelsNoteToProtobufNote(note)}, nil
}

func (srv *notesAPI) GetNote(ctx context.Context, req *notesv1.GetNoteRequest) (*notesv1.GetNoteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateGetNoteRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check user is part of the group.
	_, err = srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	note, err := srv.notes.GetNote(ctx, &models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.GetNoteResponse{Note: modelsNoteToProtobufNote(note)}, nil
}

func (srv *notesAPI) UpdateNote(ctx context.Context, req *notesv1.UpdateNoteRequest) (*notesv1.UpdateNoteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateUpdateNoteRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check user is part of the group.
	_, err = srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	note, err := srv.notes.GetNote(ctx, &models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	// Check if the user has edit access (author or in the list)
	if !hasEditPermission(note.AccountsWithEditPermissions, token.AccountID) {
		return nil, status.Error(codes.PermissionDenied, "you do not have edit permissions on this note")
	}

	updatedNote, err := srv.notes.UpdateNote(ctx,
		&models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId},
		updateNotePayloadFromUpdateNoteRequest(req),
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	err = srv.background.AddProcess(&background.Process{
		Identifier: models.NoteIdentifier{Metadata: updatedNote.ID, ActionType: models.NoteUpdateKeyword},
		CallBackFct: func() error {
			err := srv.UpdateKeywordsByNoteId(updatedNote.ID, req.GroupId, note.AuthorAccountID)
			return err
		},
		SecondsToDebounce:             5,
		CancelProcessOnSameIdentifier: true,
		RepeatProcess:                 false,
	})
	if err != nil {
		return nil, err
	}

	return &notesv1.UpdateNoteResponse{Note: modelsNoteToProtobufNote(updatedNote)}, nil
}

func updateNotePayloadFromUpdateNoteRequest(req *notesv1.UpdateNoteRequest) *models.UpdateNotePayload {
	payload := &models.UpdateNotePayload{}

	for _, path := range req.UpdateMask.Paths {
		switch path {
		case "title":
			payload.Title = req.Note.Title
		case "blocks":
			blocks := protobufBlocksToModelsBlocks(req.Note.Blocks)
			payload.Blocks = &blocks
		}
	}

	return payload
}

func (srv *notesAPI) DeleteNote(ctx context.Context, req *notesv1.DeleteNoteRequest) (*notesv1.DeleteNoteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateDeleteNoteRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = srv.notes.DeleteNote(ctx,
		&models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.DeleteNoteResponse{}, nil
}

func (srv *notesAPI) ListNotes(ctx context.Context, req *notesv1.ListNotesRequest) (*notesv1.ListNotesResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateListNoteRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check user is part of the group
	if req.GroupId == "" {
		if token.AccountID != req.AuthorAccountId {
			return nil, status.Error(codes.PermissionDenied, "could get note of another account")
		}
	} else {
		_, err = srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
		if err != nil {
			return nil, statusFromModelError(err)
		}
	}

	notes, err := srv.notes.ListNotesInternal(ctx,
		&models.ManyNotesFilter{GroupID: req.GroupId, AuthorAccountID: req.AuthorAccountId},
		&models.ListOptions{Limit: int32(req.Limit), Offset: int32(req.Offset)})
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.ListNotesResponse{Notes: modelsNotesToProtobufNotes(notes)}, nil
}

func modelsNotesToProtobufNotes(notes []*models.Note) []*notesv1.Note {
	protobufNotes := make([]*notesv1.Note, len(notes))
	for i := range notes {
		protobufNotes[i] = modelsNoteToProtobufNote(notes[i])
		// NOTE: List notes doesn't return the notes blocks but we
		// must explicitely set it to nil to avoid sending an empty
		// array.
		protobufNotes[i].Blocks = nil
	}
	return protobufNotes
}

func (srv *notesAPI) ExportNote(ctx context.Context, req *notesv1.ExportNoteRequest) (*notesv1.ExportNoteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = validators.ValidateExportNoteRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	note, err := srv.GetNote(ctx, &notesv1.GetNoteRequest{GroupId: req.GroupId, NoteId: req.NoteId})
	if err != nil {
		return nil, statusFromModelError(err)
	}

	// Check user is part of the group
	_, err = srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	formatter, ok := protobufFormatToFormatter[req.ExportFormat]
	if !ok {
		srv.logger.Error("format not recognized", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "format not recognized : %s", req.ExportFormat.String())
	}

	fileBytes, err := formatter(note.Note)

	if err != nil {
		srv.logger.Error("failed to convert note", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to convert note to: %s", req.ExportFormat.String())
	}

	return &notesv1.ExportNoteResponse{File: fileBytes}, nil
}

func (srv *notesAPI) OnAccountDelete(ctx context.Context, req *notesv1.OnAccountDeleteRequest) (*notesv1.OnAccountDeleteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = srv.notes.DeleteNotes(ctx, &models.ManyNotesFilter{AuthorAccountID: token.AccountID})
	if err != nil {
		srv.logger.Warn("Could not delete notes of " + token.AccountID + " reason " + err.Error())
	}

	err = srv.notes.RemoveEditPermissions(ctx, nil, token.AccountID)
	if err != nil {
		return nil, err
	}

	err = srv.groups.OnAccountDelete(ctx, token.AccountID)
	if err != nil {
		return nil, err
	}

	return &notesv1.OnAccountDeleteResponse{}, nil
}

func (srv *notesAPI) ListQuizs(ctx context.Context, req *notesv1.ListQuizsRequest) (*notesv1.ListQuizsResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateListQuizs(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := srv.groups.GetGroup(ctx, &models.OneGroupFilter{
		GroupID: req.GroupId,
	}, token.AccountID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "can't get group: "+err.Error())
	}

	note, err := srv.notes.GetNote(ctx, &models.OneNoteFilter{
		GroupID: req.GroupId,
		NoteID:  req.NoteId,
	}, token.AccountID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "can't get note: "+err.Error())
	}

	if res.FindMember(token.AccountID) == nil && note.AuthorAccountID != token.AccountID {
		return nil, status.Error(codes.PermissionDenied, "you don't have the access rights")
	}

	quizs, err := srv.notes.ListQuizs(ctx, &models.OneNoteFilter{
		GroupID: req.GroupId,
		NoteID:  req.NoteId,
	}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.ListQuizsResponse{
		Quizs: modelsQuizsToProtobufQuizs(quizs),
	}, nil
}

func (srv *notesAPI) GenerateQuiz(ctx context.Context, req *notesv1.GenerateQuizRequest) (*notesv1.GenerateQuizResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateGenerateQuizzRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check user is part of the group.
	_, err = srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	noteFilter := &models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId}
	note, err := srv.notes.GetNote(ctx, noteFilter, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	fullNote := noteModelToString(note)

	quiz, err := srv.language.GenerateQuizFromTextInput(fullNote, note.Lang)
	if err != nil {
		srv.logger.Error("failed to generate quiz", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to generate quiz for noteId : %s", note.ID)
	}

	_, err = srv.notes.StoreNewQuiz(ctx, noteFilter, quiz, token.AccountID)
	if err != nil {
		srv.logger.Error("failed to store quizs", zap.Error(err))
		return nil, statusFromModelError(err)
	}

	err = srv.background.AddProcess(&background.Process{
		Identifier: models.NoteIdentifier{
			ActionType: models.NoteDeleteQuiz,
			Metadata: models.Quiz{
				ID: quiz.ID,
			},
		},

		SecondsToDebounce:             uint32((time.Hour * 24 * 31) / time.Second), // TODO: Clean-up
		CancelProcessOnSameIdentifier: true,
		RepeatProcess:                 false,

		CallBackFct: func() error {
			return srv.notes.DeleteQuiz(ctx, &models.OneNoteFilter{
				GroupID: req.GroupId,
				NoteID:  req.NoteId,
			}, quiz.ID, token.AccountID)
		},
	})
	if err != nil {
		return nil, err
	}

	return &notesv1.GenerateQuizResponse{Quiz: modelsQuizToProtobufQuiz(quiz)}, nil
}

func (srv *notesAPI) GenerateSummary(ctx context.Context, req *notesv1.GenerateSummaryRequest) (*notesv1.GenerateSummaryResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateGenerateSummaryRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check user is part of the group.
	_, err = srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	note, err := srv.notes.GetNote(ctx, &models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	fullNote := noteModelToString(note)
	summary, err := srv.language.GenerateSummaryFromTextInput(fullNote, note.Lang)
	if err != nil {
		srv.logger.Error("failed to generate summarry", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to generate summarry for noteId : %s", "ok" /*note.ID*/)
	}

	return &notesv1.GenerateSummaryResponse{Summary: summary.Content}, nil
}

func (srv *notesAPI) UpdateKeywordsByNoteId(noteId string, groupId string, accountID string) error {
	note, err := srv.notes.GetNote(context.TODO(), &models.OneNoteFilter{GroupID: groupId, NoteID: noteId}, accountID)
	if err != nil {
		return statusFromModelError(err)
	}

	// TODO : mettre un timeout sur le call google
	fullNote := noteModelToString(note)

	keywords, err := srv.language.GetKeywordsFromTextInput(fullNote, note.Lang)
	if err != nil {
		srv.logger.Error("failed to gen keywords", zap.Error(err))
		return status.Errorf(codes.Internal, "failed to gen keywords for noteId : %s", note.ID)
	}

	note.Keywords = keywords

	_, err = srv.notes.UpdateNote(context.TODO(),
		&models.OneNoteFilter{GroupID: note.GroupID, NoteID: note.ID},
		&models.UpdateNotePayload{Keywords: note.Keywords},
		accountID)
	if err != nil {
		return statusFromModelError(err)
	}

	return nil
}

func (srv *notesAPI) ChangeNoteEditPermission(ctx context.Context, req *notesv1.ChangeNoteEditPermissionRequest) (*notesv1.ChangeNoteEditPermissionResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateChangeEditPermissionsRequest(req)
	if err != nil {
		return nil, err
	}

	// Check if recipient is part of the group
	group, err := srv.groups.GetGroupInternal(ctx, &models.OneGroupFilter{GroupID: req.GroupId})
	if err != nil {
		return nil, statusFromModelError(err)
	}
	if group.FindMember(req.RecipientAccountId) == nil {
		return nil, status.Error(codes.PermissionDenied, "you cannot grant permission to someone who is not part of the group")
	}

	// Store note to do later checks
	note, err := srv.notes.GetNote(ctx, &models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	requesterIsAuthor := note.AuthorAccountID == token.AccountID
	requesterIsRecipient := req.RecipientAccountId == token.AccountID

	switch req.Type {
	case notesv1.ChangeNoteEditPermissionRequest_ACTION_GRANT:
		// Requester has to be note author to grant permissions
		if !requesterIsAuthor {
			return nil, status.Error(codes.PermissionDenied, "you have to be the owner of the note to grant edit permissions")
		}

		// Grant permissions to target
		err = srv.notes.GrantNoteEditPermission(ctx, &models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId}, token.AccountID, req.RecipientAccountId)
		if err != nil {
			return nil, statusFromModelError(err)
		}
	case notesv1.ChangeNoteEditPermissionRequest_ACTION_REMOVE:

		// NOTE: I am very sad cause technically we could catch both errors with `requesterIsAuthor == requesterIsRecipient`
		// But we can't do ternary so to have 2 different error msg we have to do `if-else` :(

		if requesterIsAuthor && requesterIsRecipient {
			// Author can't remove his own rights
			return nil, status.Error(codes.PermissionDenied, "owner cannot remove his own editing rights")
		} else if !requesterIsAuthor && !requesterIsRecipient {
			// User can only remove his own rights
			return nil, status.Error(codes.PermissionDenied, "as a non-author you can only remove your own editing rights")
		}

		err = srv.notes.RemoveEditPermissions(ctx, &models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId}, req.RecipientAccountId)
		if err != nil {
			return nil, statusFromModelError(err)
		}
	}

	return &notesv1.ChangeNoteEditPermissionResponse{}, nil
}

func (srv *notesAPI) CreateBlockComment(ctx context.Context, req *notesv1.CreateBlockCommentRequest) (*notesv1.CreateBlockCommentResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateCreateBlockCommentRequest(req)
	if err != nil {
		return nil, err
	}

	res, err := srv.notes.CreateBlockComment(ctx, &models.OneBlockFilter{
		GroupID: req.GroupId,
		NoteID:  req.NoteId,
		BlockID: req.BlockId,
	}, &models.BlockComment{
		AuthorAccountID: token.AccountID,
		Content:         req.Comment.Content,
	}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.CreateBlockCommentResponse{
		Comment: modelsCommentToProtobufComment(res),
	}, nil
}

func (srv *notesAPI) DeleteBlockComment(ctx context.Context, req *notesv1.DeleteBlockCommentRequest) (*notesv1.DeleteBlockCommentResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateDeleteBlockCommentRequest(req)
	if err != nil {
		return nil, err
	}

	_, err = srv.notes.DeleteBlockComment(ctx, &models.OneBlockFilter{
		GroupID: req.GroupId,
		NoteID:  req.NoteId,
		BlockID: req.BlockId,
	}, &models.BlockComment{
		ID: req.CommentId,
	}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.DeleteBlockCommentResponse{}, nil
}

func (srv *notesAPI) ListBlockComments(ctx context.Context, req *notesv1.ListBlockCommentsRequest) (*notesv1.ListBlockCommentsResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateListBlockCommentRequest(req)
	if err != nil {
		return nil, err
	}

	res, err := srv.notes.ListBlockComments(ctx, &models.OneBlockFilter{
		GroupID: req.GroupId,
		NoteID:  req.NoteId,
		BlockID: req.BlockId,
	}, &models.ListOptions{}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	comments := make([]*notesv1.Block_Comment, 0)
	for _, modelComment := range *res {
		comments = append(comments, modelsCommentToProtobufComment(&modelComment))
	}

	return &notesv1.ListBlockCommentsResponse{
		Comments: comments,
	}, nil
}

func hasEditPermission(AccountsWithEditPermissions []string, recipientAccountID string) bool {
	for _, accountID := range AccountsWithEditPermissions {
		if accountID == recipientAccountID {
			return true
		}
	}
	return false
}

func (srv *notesAPI) authenticate(ctx context.Context) (*auth.Token, error) {
	token, err := srv.auth.TokenFromContext(ctx)
	if err != nil {
		srv.logger.Debug("could not authenticate request", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return token, nil
}

func noteModelToString(note *models.Note) string {
	var fullNote string

	for _, block := range *note.Blocks {
		if block.Type != "TYPE_CODE" && block.Type != "TYPE_IMAGE" {
			content, ok := GetBlockContent(&block)
			if ok {
				fullNote += content + "\n"
			}
		}

	}
	return fullNote
}

var protobufFormatToFormatter = map[notesv1.NoteExportFormat]func(*notesv1.Note) ([]byte, error){
	notesv1.NoteExportFormat_NOTE_EXPORT_FORMAT_MARKDOWN: exports.NoteToMarkdown,
	notesv1.NoteExportFormat_NOTE_EXPORT_FORMAT_PDF:      exports.NoteToPDF,
}

func protobufBlocksToModelsBlocks(blocks []*notesv1.Block) []models.NoteBlock {
	if blocks == nil {
		return nil
	}

	modelsBlocks := make([]models.NoteBlock, len(blocks))

	for i := range blocks {
		modelsBlocks[i] = *protobufBlockToModelsBlock(blocks[i])
	}

	return modelsBlocks
}

func protobufBlockToModelsBlock(block *notesv1.Block) *models.NoteBlock {
	modelsBlock := &models.NoteBlock{
		Type: block.Type.String(),
	}
	switch block.Type {
	case notesv1.Block_TYPE_HEADING_1:
		val := block.GetHeading()
		modelsBlock.Heading = &val
	case notesv1.Block_TYPE_HEADING_2:
		val := block.GetHeading()
		modelsBlock.Heading = &val
	case notesv1.Block_TYPE_HEADING_3:
		val := block.GetHeading()
		modelsBlock.Heading = &val
	case notesv1.Block_TYPE_PARAGRAPH:
		val := block.GetParagraph()
		modelsBlock.Paragraph = &val
	case notesv1.Block_TYPE_MATH:
		val := block.GetMath()
		modelsBlock.Math = &val
	case notesv1.Block_TYPE_CODE:
		modelsBlock.Code = &models.NoteBlockCode{}
		if code := block.GetCode(); code != nil {
			modelsBlock.Code = &models.NoteBlockCode{
				Snippet: code.Snippet,
				Lang:    code.Lang,
			}
		}
	case notesv1.Block_TYPE_IMAGE:
		modelsBlock.Image = &models.NoteBlockImage{}
		if image := block.GetImage(); image != nil {
			modelsBlock.Image = &models.NoteBlockImage{
				Caption: image.Caption,
				Url:     image.Url,
			}
		}
	case notesv1.Block_TYPE_BULLET_POINT:
		val := block.GetBulletPoint()
		modelsBlock.BulletPoint = &val
	case notesv1.Block_TYPE_NUMBER_POINT:
		val := block.GetNumberPoint()
		modelsBlock.NumberPoint = &val
	}

	for _, style := range block.Styles {
		modelStyle := models.TextStyle{
			Style: style.Style.String(),
			Position: models.Position{
				Start:  style.Pos.Start,
				Length: style.Pos.Length,
			},
		}
		if style.GetColor() != nil {
			modelStyle.Color = &models.Color{
				R: style.Color.R,
				G: style.Color.G,
				B: style.Color.B,
			}
		}
		modelsBlock.Styles = append(modelsBlock.Styles, modelStyle)
	}
	return modelsBlock
}

func modelsQuizsToProtobufQuizs(quizs *[]models.Quiz) []*notesv1.Quiz {
	res := []*notesv1.Quiz{}

	for _, quiz := range *quizs {
		res = append(res, modelsQuizToProtobufQuiz(&quiz))
	}

	return res
}

func modelsQuizToProtobufQuiz(quiz *models.Quiz) *notesv1.Quiz {
	res := &notesv1.Quiz{
		Id: quiz.ID,
	}

	for _, question := range quiz.QuizQuestions {
		res.Questions = append(res.Questions, &notesv1.QuizQuestion{
			Question:  question.Question,
			Answers:   question.Answers,
			Solutions: question.Solutions,
		})
	}
	return res
}

func modelsNoteToProtobufNote(note *models.Note) *notesv1.Note {
	var lenBlocks int

	if note.Blocks == nil {
		lenBlocks = 0
	} else {
		lenBlocks = len(*note.Blocks)
	}

	protobufNote := &notesv1.Note{
		Id:              note.ID,
		GroupId:         note.GroupID,
		AuthorAccountId: note.AuthorAccountID,
		Title:           note.Title,
		CreatedAt:       timestamppb.New(note.CreatedAt),
		ModifiedAt:      protobufTimestampOrNil(note.ModifiedAt),
		AnalyzedAt:      protobufTimestampOrNil(note.AnalyzedAt),
		Blocks:          make([]*notesv1.Block, lenBlocks),
	}

	if note.Blocks == nil {
		return protobufNote
	}

	for i := range *note.Blocks {
		protobufNote.Blocks[i] = modelsBlockToProtobufBlock(&(*note.Blocks)[i])
	}

	return protobufNote
}

func modelsBlockToProtobufBlock(block *models.NoteBlock) *notesv1.Block {
	blockType, ok := notesv1.Block_Type_value[block.Type]
	if !ok {
		blockType = int32(notesv1.Block_TYPE_INVALID)
	}
	ret := &notesv1.Block{
		Id:   block.ID,
		Type: notesv1.Block_Type(blockType),
	}

	switch notesv1.Block_Type(blockType) {
	case notesv1.Block_TYPE_HEADING_1:
		ret.Data = &notesv1.Block_Heading{
			Heading: stringPtrValueOrFallback(block.Heading, ""),
		}
	case notesv1.Block_TYPE_HEADING_2:
		ret.Data = &notesv1.Block_Heading{
			Heading: stringPtrValueOrFallback(block.Heading, ""),
		}
	case notesv1.Block_TYPE_HEADING_3:
		ret.Data = &notesv1.Block_Heading{
			Heading: stringPtrValueOrFallback(block.Heading, ""),
		}
	case notesv1.Block_TYPE_PARAGRAPH:
		ret.Data = &notesv1.Block_Paragraph{
			Paragraph: stringPtrValueOrFallback(block.Paragraph, ""),
		}
	case notesv1.Block_TYPE_CODE:
		if block.Code == nil {
			break
		}
		ret.Data = &notesv1.Block_Code_{
			Code: &notesv1.Block_Code{
				Snippet: block.Code.Snippet,
				Lang:    block.Code.Lang,
			},
		}
	case notesv1.Block_TYPE_MATH:
		ret.Data = &notesv1.Block_Math{
			Math: stringPtrValueOrFallback(block.Math, ""),
		}
	case notesv1.Block_TYPE_IMAGE:
		if block.Image == nil {
			break
		}
		ret.Data = &notesv1.Block_Image_{
			Image: &notesv1.Block_Image{
				Caption: block.Image.Caption,
				Url:     block.Image.Url,
			},
		}
	case notesv1.Block_TYPE_BULLET_POINT:
		ret.Data = &notesv1.Block_BulletPoint{
			BulletPoint: stringPtrValueOrFallback(block.BulletPoint, ""),
		}
	case notesv1.Block_TYPE_NUMBER_POINT:
		ret.Data = &notesv1.Block_NumberPoint{
			NumberPoint: stringPtrValueOrFallback(block.NumberPoint, ""),
		}
	}

	for _, style := range block.Styles {
		modelStyle := &notesv1.Block_TextStyle{
			Style: notesv1.Block_TextStyle_Style(notesv1.Block_TextStyle_Style_value[style.Style]),
			Pos: &notesv1.Block_TextStyle_Position{
				Start:  int64(style.Position.Start),
				Length: int64(style.Position.Length),
			},
		}
		if style.Color != nil {
			modelStyle.Color = &notesv1.Block_TextStyle_Color{
				R: int32(style.Color.R),
				G: int32(style.Color.G),
				B: int32(style.Color.B),
			}
		}

		ret.Styles = append(ret.Styles, modelStyle)
	}

	return ret
}

func stringPtrValueOrFallback(ptr *string, fallback string) string {
	if ptr != nil {
		return *ptr
	}
	return fallback
}

func modelsCommentToProtobufComment(cmt *models.BlockComment) *notesv1.Block_Comment {
	protobufComment := &notesv1.Block_Comment{
		Id:       cmt.ID,
		AuthorId: cmt.AuthorAccountID,
		Content:  cmt.Content,
	}
	return protobufComment
}
