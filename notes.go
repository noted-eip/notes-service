package main

import (
	"context"

	"notes-service/auth"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"
	"notes-service/validators"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type notesService struct {
	notespb.UnimplementedNotesAPIServer

	auth      auth.Service
	logger    *zap.Logger
	repoNote  models.NotesRepository
	repoBlock models.BlocksRepository
}

var _ notespb.NotesAPIServer = &notesService{}

func (srv *notesService) CreateNote(ctx context.Context, in *notespb.CreateNoteRequest) (*notespb.CreateNoteResponse, error) {
	_, err := Authenticate(srv, ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateCreateNoteRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	note, err := srv.repoNote.Create(ctx, &models.NotePayload{AuthorId: in.Note.AuthorId, Title: in.Note.Title})

	if err != nil {
		srv.logger.Error("failed to create note", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not create note")
	}

	blocks := make([]models.Block, len(in.Note.Blocks))

	for index, block := range in.Note.Blocks {
		err := convertApiBlockToModelBlock(&blocks[index], block)
		if err != nil {
			srv.logger.Error("failed to create note", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "invalid content provided for block index : %d", index)
		}
		srv.repoBlock.Create(ctx, &models.Block{NoteId: note.ID, Type: uint32(in.Note.Blocks[index].Type), Index: uint32(index + 1), Content: blocks[index].Content}) // NOTE: Shouldn't we start index at 0 ?
	}
	noteResponse := notespb.Note{Id: note.ID, AuthorId: note.AuthorId, Title: note.Title, Blocks: in.Note.Blocks, CreatedAt: timestamppb.New(note.CreationDate), ModifiedAt: timestamppb.New(note.ModificationDate)}
	return &notespb.CreateNoteResponse{Note: &noteResponse}, nil
}

func (srv *notesService) GetNote(ctx context.Context, in *notespb.GetNoteRequest) (*notespb.GetNoteResponse, error) {
	_, err := Authenticate(srv, ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateGetNoteRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//Check if user is a note group member

	note, err := srv.repoNote.Get(ctx, in.Id)
	if err != nil {
		srv.logger.Error("failed to get note", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not get note.")
	}

	blocksTmp, err := srv.repoBlock.GetBlocks(ctx, note.ID)
	if err != nil {
		srv.logger.Error("failed to get blocks", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "invalid content provided for blocks form noteId : %s", note.ID)
	}

	//Convert []models.block to []notespb.Block
	blocks := make([]*notespb.Block, len(blocksTmp))
	for index, block := range blocksTmp {
		blocks[index] = &notespb.Block{}
		err := convertModelBlockToApiBlock(block, blocks[index])
		if err != nil {
			srv.logger.Error("failed to the content of a block", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "fail to get content from block Id : %s", block.ID)
		}
		blocks[index] = &notespb.Block{Id: block.ID, Type: notespb.Block_Type(block.Type), Data: blocks[index].Data}
	}
	noteResponse := notespb.Note{Id: note.ID, AuthorId: note.AuthorId, Title: note.Title, Blocks: blocks, CreatedAt: timestamppb.New(note.CreationDate), ModifiedAt: timestamppb.New(note.ModificationDate)}
	return &notespb.GetNoteResponse{Note: &noteResponse}, nil
}

func (srv *notesService) UpdateNote(ctx context.Context, in *notespb.UpdateNoteRequest) (*notespb.UpdateNoteResponse, error) {
	token, err := Authenticate(srv, ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateUpdateNoteRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//Check if the user own the note
	note, err := srv.repoNote.Get(ctx, in.Id)
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not update note")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	}

	//Delete all blocks
	err = srv.repoBlock.DeleteBlocks(ctx, in.Id)
	if err != nil {
		srv.logger.Error("blocks weren't deleted : ", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not delete blocks")
	}

	//Todo : Update only the note metInformation if the blocks are nil
	err = srv.repoNote.Update(ctx, in.Id, &models.NotePayload{AuthorId: in.Note.AuthorId, Title: in.Note.Title})
	if err != nil {
		srv.logger.Error("failed to update note", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not update note")
	}

	//Create all blocks
	blocks := make([]models.Block, len(in.Note.Blocks))
	for index, block := range in.Note.Blocks {
		err := convertApiBlockToModelBlock(&blocks[index], block)
		if err != nil {
			srv.logger.Error("failed to update blocks", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "invalid content provided for block index : %d", index)
		}
		srv.repoBlock.Create(ctx, &models.Block{NoteId: in.Id, Type: uint32(in.Note.Blocks[index].Type), Index: uint32(index + 1), Content: blocks[index].Content})
	}

	return nil, nil
}

func (srv *notesService) DeleteNote(ctx context.Context, in *notespb.DeleteNoteRequest) (*notespb.DeleteNoteResponse, error) {
	token, err := Authenticate(srv, ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateDeleteNoteRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//Check if the user own the note
	note, err := srv.repoNote.Get(ctx, in.Id)
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not update note")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	}
	//Delete all blocks related to the noteId
	err = srv.repoBlock.DeleteBlocks(ctx, in.Id)
	if err != nil {
		srv.logger.Error("blocks weren't deleted : ", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not delete blocks")
	}
	//Delete the note
	err = srv.repoNote.Delete(ctx, in.Id)
	if err != nil {
		srv.logger.Error("failed to delete note", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not delete note")
	}

	return nil, nil
}

func (srv *notesService) ListNotes(ctx context.Context, in *notespb.ListNotesRequest) (*notespb.ListNotesResponse, error) {
	_, err := Authenticate(srv, ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateListNoteRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//TODO: check if user is a note group member

	notes, err := srv.repoNote.List(ctx, in.AuthorId)
	if err != nil {
		srv.logger.Error("failed to get note", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "could not get note")
	}

	notesResponse := make([]*notespb.Note, len(notes))
	for index, note := range notes {
		notesResponse[index] = &notespb.Note{Id: note.ID, AuthorId: note.AuthorId, Title: note.Title, CreatedAt: timestamppb.New(note.CreationDate), ModifiedAt: timestamppb.New(note.ModificationDate)}
	}
	return &notespb.ListNotesResponse{Notes: notesResponse}, nil
}

func convertApiBlockToModelBlock(block *models.Block, blockRequest *notespb.Block) error {
	switch op := blockRequest.Data.(type) {
	case *notespb.Block_Heading:
		block.Content = op.Heading
	case *notespb.Block_Paragraph:
		block.Content = op.Paragraph
	case *notespb.Block_NumberPoint:
		block.Content = op.NumberPoint
	case *notespb.Block_BulletPoint:
		block.Content = op.BulletPoint
	case *notespb.Block_Math:
		block.Content = op.Math
	case *notespb.Block_Image_:
		block.Image.Caption = op.Image.Caption
		block.Image.Url = op.Image.Url
	case *notespb.Block_Code_:
		block.Code.Lang = op.Code.Lang
		block.Code.Snippet = op.Code.Snippet

	default:
		return status.Error(codes.Internal, "no data in this block")
	}
	return nil
}

func convertModelBlockToApiBlock(blockSrc *models.Block, blockDest *notespb.Block) error {
	switch blockSrc.Type {
	case uint32(notespb.Block_TYPE_HEADING_1):
		fallthrough
	case uint32(notespb.Block_TYPE_HEADING_2):
		fallthrough
	case uint32(notespb.Block_TYPE_HEADING_3):
		blockDest.Data = &notespb.Block_Heading{Heading: blockSrc.Content}
	case uint32(notespb.Block_TYPE_PARAGRAPH):
		blockDest.Data = &notespb.Block_Paragraph{Paragraph: blockSrc.Content}
	case uint32(notespb.Block_TYPE_NUMBERED_POINT):
		blockDest.Data = &notespb.Block_NumberPoint{NumberPoint: blockSrc.Content}
	case uint32(notespb.Block_TYPE_BULLET_POINT):
		blockDest.Data = &notespb.Block_BulletPoint{BulletPoint: blockSrc.Content}
	case uint32(notespb.Block_TYPE_MATH):
		blockDest.Data = &notespb.Block_Math{Math: blockSrc.Content}
	case uint32(notespb.Block_TYPE_IMAGE):
		(*blockDest).Data = &notespb.Block_Image_{Image: &notespb.Block_Image{Caption: blockSrc.Image.Caption, Url: blockSrc.Image.Url}}
	case uint32(notespb.Block_TYPE_CODE):
		(*blockDest).Data = &notespb.Block_Code_{Code: &notespb.Block_Code{Snippet: blockSrc.Code.Snippet, Lang: blockSrc.Code.Lang}}
	default:
		return status.Errorf(codes.Internal, "no such content in this block")
	}
	return nil
}

func Authenticate(srv *notesService, ctx context.Context) (*auth.Token, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return token, nil
}

func (srv *notesService) authenticate(ctx context.Context) (*auth.Token, error) {
	token, err := srv.auth.TokenFromContext(ctx)
	if err != nil {
		srv.logger.Debug("failed to authenticate request", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return token, nil
}
