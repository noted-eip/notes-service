package main

import (
	"context"

	"notes-service/auth"
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

	auth     auth.Service
	language language.Service

	notes  models.NotesRepository
	groups models.GroupsRepository
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

	note, err := srv.notes.CreateNote(ctx, &models.CreateNotePayload{
		GroupID:         req.GroupId,
		Title:           req.Title,
		AuthorAccountID: token.AccountID,
		FolderID:        "",
		Blocks:          protobufBlocksToModelsBlocks(req.Blocks),
	}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
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

	note, err := srv.notes.UpdateNote(ctx,
		&models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId},
		&models.UpdateNotePayload{Title: req.Note.Title},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.UpdateNoteResponse{Note: modelsNoteToProtobufNote(note)}, nil
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
	_, err = srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	notes, err := srv.notes.ListNotesInternal(ctx, &models.ManyNotesFilter{AuthorAccountID: req.AuthorAccountId}, &models.ListOptions{Limit: 20, Offset: 0})
	if err != nil {
		return nil, statusFromModelError(err)
	}

	protobufNotes := make([]*notesv1.Note, len(notes))
	for i := range notes {
		protobufNotes[i] = modelsNoteToProtobufNote(notes[i])
	}

	return &notesv1.ListNotesResponse{Notes: protobufNotes}, nil
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
		return nil, err
	}

	err = srv.groups.OnAccountDelete(ctx, token.AccountID)
	if err != nil {
		return nil, err
	}

	return &notesv1.OnAccountDeleteResponse{}, nil
}

func (srv *notesAPI) authenticate(ctx context.Context) (*auth.Token, error) {
	token, err := srv.auth.TokenFromContext(ctx)
	if err != nil {
		srv.logger.Debug("could not authenticate request", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return token, nil
}

var protobufFormatToFormatter = map[notesv1.NoteExportFormat]func(*notesv1.Note) ([]byte, error){
	notesv1.NoteExportFormat_NOTE_EXPORT_FORMAT_MARKDOWN: exports.NoteToMarkdown,
	notesv1.NoteExportFormat_NOTE_EXPORT_FORMAT_PDF:      exports.NoteToPDF,
}

func protobufBlocksToModelsBlocks(blocks []*notesv1.Block) []models.NoteBlock {
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
	return modelsBlock
}

func modelsNoteToProtobufNote(note *models.Note) *notesv1.Note {
	protobufNote := &notesv1.Note{
		Id:              note.ID,
		GroupId:         note.GroupID,
		AuthorAccountId: note.AuthorAccountID,
		Title:           note.Title,
		CreatedAt:       timestamppb.New(note.CreatedAt),
		ModifiedAt:      timestamppb.New(note.ModifiedAt),
		AnalyzedAt:      timestamppb.New(note.AnalyzedAt),
		Blocks:          make([]*notesv1.Block, len(note.Blocks)),
	}

	for i := range note.Blocks {
		protobufNote.Blocks[i] = modelsBlockToProtobufBlock(&note.Blocks[i])
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

	return ret
}

func stringPtrValueOrFallback(ptr *string, fallback string) string {
	if ptr != nil {
		return *ptr
	}
	return fallback
}
