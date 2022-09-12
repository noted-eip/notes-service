package main

import (
	"context"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type blocksService struct {
	notespb.UnimplementedNotesAPIServer

	logger *zap.SugaredLogger
	repo   models.BlocksRepository
}

var _ notespb.NotesAPIServer = &notesService{}

func (srv *blocksService) InsertBlock(ctx context.Context, in *notespb.InsertBlockRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(in.Block.Id)
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get block")
	}
	//notespb.Block_Heading.Heading
	content := "contentTestZebi"
	/*for _, patch := range in.Block.Data {
		switch op := patch.Op.(type) {
		case *notespb.Block_Paragraph:
			fmt.Printf("Paragraph")
		case *notespb.Block_Heading:
			fmt.Printf("Heading")
		default:
			fmt.Println("No matching operations")
		}
	}*/

	srv.repo.Create(
		ctx,
		&models.BlockWithIndex{ID: id.String(), NoteId: in.Block.Id, Type: uint32(in.Block.Type), Index: in.Index, Content: &content})
	return &emptypb.Empty{}, nil
}

func (srv *blocksService) UpdateBlock(ctx context.Context, in *notespb.UpdateBlockRequest) (*emptypb.Empty, error) {
	_, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get block")
	}

	srv.repo.Update(
		ctx,
		&models.BlockFilter{BlockId: in.Id, NoteId: ""},
		&models.BlockWithIndex{ID: in.Id, NoteId: "", Type: uint32(in.Block.Type), Index: in.Index, Content: nil /*&in.Block.Data*/})
	return &emptypb.Empty{}, nil
}

func (srv *blocksService) DeleteBlock(ctx context.Context, in *notespb.DeleteBlockRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get note")
	}
	srv.repo.Delete(ctx, &models.BlockFilter{BlockId: id.String(), NoteId: ""})
	return &emptypb.Empty{}, nil
}
