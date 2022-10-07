package main

import (
	"context"
	"fmt"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"

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
	return &emptypb.Empty{}, nil
}

func (srv *blocksService) UpdateBlock(ctx context.Context, in *notespb.UpdateBlockRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (srv *blocksService) DeleteBlock(ctx context.Context, in *notespb.DeleteBlockRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func FillContentFromModelToApi(blockRequest *models.BlockWithIndex, contentType uint32, blockApi *notespb.Block) error {
	switch contentType {
	case 1:
		blockApi.Data = &notespb.Block_Heading{Heading: blockRequest.Content}
	case 2:
		blockApi.Data = &notespb.Block_Paragraph{Paragraph: blockRequest.Content}
	case 3:
		blockApi.Data = &notespb.Block_NumberPoint{NumberPoint: blockRequest.Content}
	case 4:
		blockApi.Data = &notespb.Block_BulletPoint{BulletPoint: blockRequest.Content}
	case 5:
		blockApi.Data = &notespb.Block_Math{Math: blockRequest.Content}
	/*
		case 6:
			(*blockApi).Data = &notespb.Block_Image_{Image: {caption: blockRequest.Image.caption, url: blockRequest.Image.url}}
		case 7:
			(*blockApi).Data = &notespb.Block_Code_{Code: {sinppet: blockRequest.Code.Snippet, lang: blockRequest.Code.Lang}}
	*/
	default:
		fmt.Println("No such content in this block")
		return status.Errorf(codes.Internal, "no such content in this block")
	}
	return nil
}
