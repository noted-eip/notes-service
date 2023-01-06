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

//il faudrait un index aussi ici
func (srv *notesService) InsertBlock(ctx context.Context, in *notespb.InsertBlockRequest) (*notespb.InsertBlockResponse, error) {
	token, err := Authenticate(srv, ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateInsertBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//Check if the user own the note
	note, err := srv.repoNote.Get(ctx, in.NoteId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "could not get block")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to create a block")
	}

	var block = models.Block{}
	err = convertApiBlockToModelBlock(&block, in.Block)
	if err != nil {
		srv.logger.Error("failed to create block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid data content provided for block index : %d", in.Index)
	}
	BlockId, err := srv.repoBlock.Create(ctx, &models.Block{NoteId: in.NoteId, Type: uint32(in.Block.Type), Index: in.Index, Content: block.Content})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't create block id : %s", *BlockId)
	}
	blockResponse := &notespb.Block{Id: *BlockId, Type: in.Block.Type, Data: in.Block.Data}
	return &notespb.InsertBlockResponse{Block: blockResponse}, nil
}

func (srv *notesService) UpdateBlock(ctx context.Context, in *notespb.UpdateBlockRequest) (*notespb.UpdateBlockResponse, error) {
	token, err := Authenticate(srv, ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateUpdateBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//check if the block exist
	block, err := srv.repoBlock.GetBlock(ctx, in.Id)
	if err != nil {
		srv.logger.Error("Block not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not delete block")
	}
	//Check if the user own the note
	note, err := srv.repoNote.Get(ctx, block.NoteId)
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not upate block")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	}

	var blockUpated = models.Block{}
	err = convertApiBlockToModelBlock(&blockUpated, in.Block)
	if err != nil {
		srv.logger.Error("failed to update block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid content provided for block id : %s", in.Id)
	}

	GetLastBlockIndex()
	//si in.Index > lastIndex
	// {  in.Index = lastIndex  }
	UpdateBlocksByIndex()
	//if actual index (in.Index) < old index (block.Index) :
	// -> { tout les indexs entre l'index actuel (in.Index) et l'ancien index (block.Index) : Index = +1 } (>=)
	//if old index (block.Index) < actual index (in.Index) :
	// -> { tout les indexs entre l'index actuel (in.Index) et l'ancien index (block.Index) : Index = -1 } (<=)

	print("### Block upated index = " + strconv.Itoa(int(in.Index)) + " & block content = " + blockUpated.Content + "\n")

	srv.repoBlock.Update(ctx, in.Id, &models.Block{ID: in.Id, Type: uint32(in.Block.Type), Index: in.Index, Content: blockUpated.Content})
	return &notespb.UpdateBlockResponse{
		Block: &notespb.Block{
			Id:   in.Id,
			Type: in.Block.Type,
			Data: in.Block.Data,
		},
	}, nil
}

func (srv *notesService) DeleteBlock(ctx context.Context, in *notespb.DeleteBlockRequest) (*notespb.DeleteBlockResponse, error) {
	token, err := Authenticate(srv, ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateDeleteBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//check if the block exist
	block, err := srv.repoBlock.GetBlock(ctx, in.Id)
	if err != nil {
		srv.logger.Error("Block not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not delete block")
	}
	//Check if the user own the note
	note, err := srv.repoNote.Get(ctx, block.NoteId)
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not delete block")
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
