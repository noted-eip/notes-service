package main

import (
	"context"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"
	"strconv"

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
	/*_, err := uuid.Parse(strconv.Itoa(int(in.NoteId)))
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get block")
	}*/

	//in.Block.Data;//string or Image or code
	content := "test-content"

	srv.repo.Create(
		ctx,
		&models.BlockWithIndex{ID: "test-id", NoteId: strconv.Itoa(int(in.NoteId)), Type: uint32(in.Block.Type), Index: in.Index, Content: &content})
	return &emptypb.Empty{}, nil
}

func (srv *blocksService) UpdateBlock(ctx context.Context, in *notespb.UpdateBlockRequest) (*emptypb.Empty, error) {
	/*_, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get block")
	}*/
	content := "test-content-updated"

	srv.repo.Update(
		ctx,
		&models.BlockFilter{BlockId: in.Id, NoteId: ""},
		&models.BlockWithIndex{ID: in.Id, NoteId: "", Type: uint32(in.Block.Type), Index: in.Index, Content: &content})
	return &emptypb.Empty{}, nil
}

func (srv *blocksService) DeleteBlock(ctx context.Context, in *notespb.DeleteBlockRequest) (*emptypb.Empty, error) {
	/*id, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get note")
	}*/
	srv.repo.Delete(ctx, &models.BlockFilter{BlockId: in.Id, NoteId: ""})
	return &emptypb.Empty{}, nil
}
