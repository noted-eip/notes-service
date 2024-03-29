package main

import (
	"context"
	"notes-service/models"

	background "github.com/noted-eip/noted/background-service"

	notesv1 "notes-service/protorepo/noted/notes/v1"
	"notes-service/validators"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *notesAPI) InsertBlock(ctx context.Context, req *notesv1.InsertBlockRequest) (*notesv1.InsertBlockResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateInsertBlockRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	block, err := srv.notes.InsertBlock(ctx,
		&models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId},
		&models.InsertNoteBlockPayload{
			Index: uint(req.Index),
			Block: *protobufBlockToModelsBlock(req.Block),
		},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	// Launch process to generate keywords in 15minutes after the last modification
	srv.background.AddProcess(&background.Process{
		Identifier: models.NoteIdentifier{Metadata: req.NoteId, ActionType: models.NoteUpdateKeyword},
		CallBackFct: func() error {
			err := srv.UpdateKeywordsByNoteId(req.NoteId, req.GroupId, token.AccountID)
			return err
		},
		SecondsToDebounce:             5,
		CancelProcessOnSameIdentifier: true,
		RepeatProcess:                 false,
	})

	return &notesv1.InsertBlockResponse{Block: modelsBlockToProtobufBlock(block)}, nil
}

func (srv *notesAPI) UpdateBlockIndex(ctx context.Context, req *notesv1.UpdateBlockIndexRequest) (*notesv1.UpdateBlockIndexResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateUpdateBlockIndexRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	block, err := srv.notes.GetBlock(ctx,
		&models.OneBlockFilter{GroupID: req.GroupId, NoteID: req.NoteId, BlockID: req.BlockId},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}
	err = srv.notes.DeleteBlock(ctx,
		&models.OneBlockFilter{GroupID: req.GroupId, NoteID: req.NoteId, BlockID: req.BlockId},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}
	block, err = srv.notes.InsertBlock(ctx,
		&models.OneNoteFilter{GroupID: req.GroupId, NoteID: req.NoteId},
		&models.InsertNoteBlockPayload{
			Index: uint(req.Index),
			Block: *block,
		},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.UpdateBlockIndexResponse{Block: modelsBlockToProtobufBlock(block)}, nil
}

func (srv *notesAPI) UpdateBlock(ctx context.Context, req *notesv1.UpdateBlockRequest) (*notesv1.UpdateBlockResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateUpdateBlockRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	block, err := srv.notes.UpdateBlock(ctx,
		&models.OneBlockFilter{GroupID: req.GroupId, NoteID: req.NoteId, BlockID: req.BlockId},
		&models.UpdateBlockPayload{
			Block: *protobufBlockToModelsBlock(req.Block),
		},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	// Launch process to generate keywords in 15minutes after the last modification
	srv.background.AddProcess(&background.Process{
		Identifier: models.NoteIdentifier{Metadata: req.NoteId, ActionType: models.NoteUpdateKeyword},
		CallBackFct: func() error {
			err := srv.UpdateKeywordsByNoteId(req.NoteId, req.GroupId, token.AccountID)
			return err
		},
		SecondsToDebounce:             5,
		CancelProcessOnSameIdentifier: true,
		RepeatProcess:                 false,
	})

	return &notesv1.UpdateBlockResponse{Block: modelsBlockToProtobufBlock(block)}, nil
}

func (srv *notesAPI) DeleteBlock(ctx context.Context, req *notesv1.DeleteBlockRequest) (*notesv1.DeleteBlockResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateDeleteBlockRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = srv.notes.DeleteBlock(ctx,
		&models.OneBlockFilter{GroupID: req.GroupId, NoteID: req.NoteId, BlockID: req.BlockId},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	// Launch process to generate keywords in 15minutes after the last modification
	srv.background.AddProcess(&background.Process{
		Identifier: models.NoteIdentifier{Metadata: req.NoteId, ActionType: models.NoteUpdateKeyword},
		CallBackFct: func() error {
			err := srv.UpdateKeywordsByNoteId(req.NoteId, req.GroupId, token.AccountID)
			return err
		},
		SecondsToDebounce:             5,
		CancelProcessOnSameIdentifier: true,
		RepeatProcess:                 false,
	})

	return &notesv1.DeleteBlockResponse{}, nil
}
