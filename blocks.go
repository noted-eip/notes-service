package main

import (
	"context"
	"fmt"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"
	"strconv"

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
	_, err := uuid.Parse(strconv.Itoa(int(in.NoteId)))
	if err != nil {
		srv.logger.Errorw("invalid uuid", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not insert block")
	}

	if in.Block.Data == nil || in.Index < 1 || in.Block.Type < 1 {
		srv.logger.Errorw("invalid arguments", err.Error())
		return nil, status.Errorf(codes.Internal, "could not insert block")
	}

	var block = models.Block{}
	err = FillBlockContent(&block, in.Block)
	if err != nil {
		srv.logger.Errorw("failed to create block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid content provided for block index : ", in.Index)
	}

	srv.repo.Create(ctx, &models.BlockWithIndex{NoteId: strconv.Itoa(int(in.NoteId)), Type: uint32(in.Block.Type), Index: in.Index, Content: block.Content})
	return &emptypb.Empty{}, nil
}

func (srv *blocksService) UpdateBlock(ctx context.Context, in *notespb.UpdateBlockRequest) (*emptypb.Empty, error) {
	_, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Errorw("invalid uuid", err.Error())
		return nil, status.Errorf(codes.Internal, "could not update block")
	}

	var block = models.Block{}
	err = FillBlockContent(&block, in.Block)
	if err != nil {
		srv.logger.Errorw("failed to update block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid content provided for block id : ", in.Id)
	}

	srv.repo.Update(ctx,
		&models.BlockFilter{BlockId: in.Id},
		&models.BlockWithIndex{ID: in.Id, Type: uint32(in.Block.Type), Index: in.Index, Content: block.Content})
	return &emptypb.Empty{}, nil
}

func (srv *blocksService) DeleteBlock(ctx context.Context, in *notespb.DeleteBlockRequest) (*emptypb.Empty, error) {
	_, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Errorw("invalid uuid", err.Error())
		return nil, status.Errorf(codes.Internal, "could not delete block")
	}

	err = srv.repo.Delete(ctx, &models.BlockFilter{BlockId: in.Id, NoteId: ""})
	if err != nil {
		srv.logger.Errorw("block was not deleted : ", err.Error())
		return nil, status.Errorf(codes.Internal, "could not delete block")
	}

	return &emptypb.Empty{}, nil
}

func FillBlockContent(block *models.Block, blockRequest *notespb.Block) error {
	switch op := blockRequest.Data.(type) {
	case *notespb.Block_Heading:
		block.Content = &op.Heading
	case *notespb.Block_Paragraph:
		block.Content = &op.Paragraph
	case *notespb.Block_NumberPoint:
		block.Content = &op.NumberPoint
	case *notespb.Block_BulletPoint:
		block.Content = &op.BulletPoint
	case *notespb.Block_Math:
		block.Content = &op.Math
	/*
		case *notespb.Block_Image_:
			block.Image.caption = &op.Image.Caption
			block.Content = &op.Image.Url
		case *notespb.Block_Code_:
			block.Content = &op.Code.Lang
			block.Content = &op.Code.Snippet
	*/
	default:
		fmt.Println("No Data in this block")
		return status.Errorf(codes.Internal, "no data in this block")
	}
	return nil
}
