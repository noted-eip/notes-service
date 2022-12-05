package main

import (
	"context"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"
	"notes-service/validators"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *notesService) InsertBlock(ctx context.Context, in *notespb.InsertBlockRequest) (*notespb.InsertBlockResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateInsertBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//check si la note appartient a celui qui veut la modifier
	note, err := srv.repoNote.Get(ctx, strconv.Itoa(int(in.NoteId)))
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not update note")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	}

	var block = models.Block{}
	err = fillBlockContent(&block, in.Block)
	if err != nil {
		srv.logger.Error("failed to create block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid content provided for block index : %d", in.Index)
	}

	BlockId, err := srv.repoBlock.Create(ctx, &models.Block{NoteId: strconv.Itoa(int(in.NoteId)), Type: uint32(in.Block.Type), Index: in.Index, Content: block.Content})

	return &notespb.InsertBlockResponse{
		Block: &notespb.Block{
			Id:   *BlockId,
			Type: in.Block.Type,
			Data: in.Block.Data,
		},
	}, nil
}

func (srv *notesService) UpdateBlock(ctx context.Context, in *notespb.UpdateBlockRequest) (*notespb.UpdateBlockResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateUpdateBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//check si la note appartient a celui qui veut la modifier
	note, err := srv.repoNote.Get(ctx, in.Id)
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not update note")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	}

	var block = models.Block{}
	err = fillBlockContent(&block, in.Block)
	if err != nil {
		srv.logger.Error("failed to update block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid content provided for block id : %s", in.Id)
	}

	srv.repoBlock.Update(ctx, in.Id, &models.Block{ID: in.Id, Type: uint32(in.Block.Type), Index: in.Index, Content: block.Content})
	return &notespb.UpdateBlockResponse{
		Block: &notespb.Block{
			Id:   in.Id,
			Type: in.Block.Type,
			Data: in.Block.Data,
		},
	}, nil
}

func (srv *notesService) DeleteBlock(ctx context.Context, in *notespb.DeleteBlockRequest) (*notespb.DeleteBlockResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateDeleteBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//check si la note appartient a celui qui veut la modifier
	note, err := srv.repoNote.Get(ctx, in.Id)
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not update note")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	}

	err = srv.repoBlock.DeleteBlock(ctx, in.Id)
	if err != nil {
		srv.logger.Error("block was not deleted : ", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not delete block")
	}

	return nil, nil
}
