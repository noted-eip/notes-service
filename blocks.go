package main

import (
	"context"
	"fmt"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"
	recommendationspb "notes-service/protorepo/noted/recommendations/v1"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type blocksService struct {
	notespb.UnimplementedNotesAPIServer

	logger               *zap.Logger
	repo                 models.BlocksRepository
	recommendationClient recommendationspb.RecommendationsAPIClient
}

var _ notespb.NotesAPIServer = &notesService{}

func (srv *blocksService) InsertBlock(ctx context.Context, in *notespb.InsertBlockRequest) (*notespb.InsertBlockResponse, error) {
	_, err := uuid.Parse(strconv.Itoa(int(in.NoteId)))
	if err != nil {
		srv.logger.Error("invalid uuid :", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "could not insert block")
	}

	if in.Block.Data == nil || in.Index < 1 || in.Block.Type < 1 {
		srv.logger.Error("invalid arguments")
		return nil, status.Errorf(codes.InvalidArgument, "could not insert block")
	}
	//Convert oneof Data to model content
	block := models.Block{}
	err = FillBlockContent(&block, in.Block)
	if err != nil {
		srv.logger.Error("failed to create block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid content provided for block index : %d", in.Index)
	}
	//Get recommendation tags
	blockContent, err := GetDataContent(in.Block)
	if err != nil {
		srv.logger.Error("failed to convert the content of the block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to convert the content of the block : ", in.Index)
	}
	recommendationRequest := &recommendationspb.ExtractKeywordsRequest{Content: blockContent}
	clientResponse, err := srv.recommendationClient.ExtractKeywords(ctx, recommendationRequest)
	if err != nil {
		srv.logger.Error("failed to get the recommendation from client", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get the recommendation from client")
	}

	blockId, err := srv.repo.Create(ctx, &models.BlockWithTags{NoteId: strconv.Itoa(int(in.NoteId)), Type: uint32(in.Block.Type), Index: in.Index, Content: block.Content, Tags: clientResponse.Keywords})
	if err != nil {
		srv.logger.Error("failed to create the block in DB", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create the block in DB")
	}

	return &notespb.InsertBlockResponse{
		Block: &notespb.Block{
			Id:   *blockId,
			Type: in.Block.Type,
			Data: in.Block.Data,
		},
	}, nil
}

func (srv *blocksService) UpdateBlock(ctx context.Context, in *notespb.UpdateBlockRequest) (*emptypb.Empty, error) {
	_, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Error("invalid uuid", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not update block")
	}

	var block = models.Block{}
	err = FillBlockContent(&block, in.Block)
	if err != nil {
		srv.logger.Error("failed to update block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid content provided for block id : %s", in.Id)
	}

	srv.repo.Update(ctx, &in.Id, &models.BlockWithIndex{ID: in.Id, Type: uint32(in.Block.Type), Index: in.Index, Content: block.Content})
	return nil, nil
}

func (srv *blocksService) DeleteBlock(ctx context.Context, in *notespb.DeleteBlockRequest) (*emptypb.Empty, error) {
	_, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Error("invalid uuid", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "could not delete block")
	}

	err = srv.repo.DeleteBlock(ctx, &in.Id)
	if err != nil {
		srv.logger.Error("block was not deleted : ", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not delete block")
	}

	return &emptypb.Empty{}, nil
}

func GetDataContent(blockRequest *notespb.Block) (string, error) {
	switch op := blockRequest.Data.(type) {
	case *notespb.Block_Heading:
		return op.Heading, nil
	case *notespb.Block_Paragraph:
		return op.Paragraph, nil
	case *notespb.Block_NumberPoint:
		return op.NumberPoint, nil
	case *notespb.Block_BulletPoint:
		return op.BulletPoint, nil
	case *notespb.Block_Math:
		return op.Math, nil
	}
	return "", nil
}

func FillBlockContent(block *models.Block, blockRequest *notespb.Block) error {
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
	/*
		case *notespb.Block_Image_:
			block.Image.caption = &op.Image.Caption
			block.Image.url = &op.Image.Url
		case *notespb.Block_Code_:
			block.Code.lang = &op.Code.Lang
			block.Code.Snippet = &op.Code.Snippet
	*/
	default:
		fmt.Println("No Data in this block")
		return status.Error(codes.Internal, "no data in this block")
	}
	return nil
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
