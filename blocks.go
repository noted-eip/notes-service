package main

import (
	"context"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type blocksService struct {
	notespb.UnimplementedNotesAPIServer

	logger *zap.SugaredLogger
	repo   models.BlocksRepository
}

var _ notespb.NotesAPIServer = &notesService{}

func (srv *blocksService) InsertBlock(ctx context.Context, in *notespb.InsertBlockRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (srv *blocksService) UpdateBlock(ctx context.Context, in *notespb.UpdateBlockRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (srv *blocksService) DeleteBlock(ctx context.Context, in *notespb.DeleteBlockRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
