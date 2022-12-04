package main

/*
import (
	"context"
	"fmt"
	"notes-service/auth"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"
	"notes-service/validators"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)


type blocksService struct {
	notespb.UnimplementedNotesAPIServer

	auth   auth.Service
	logger *zap.Logger
	repo   models.BlocksRepository
}

var _ notespb.NotesAPIServer = &notesService{}


//si c ta note seulement
func (srv *notesService) InsertBlock(ctx context.Context, in *notespb.InsertBlockRequest) (*notespb.InsertBlockResponse, error) {
	_, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// La il va falloir faire le get d'une note, donc avoir le contexte de noteService
	// Donc mieux vaut tout regrouper dans une service
	//note, err := srv.repoNote.Get(ctx, in.Id)
	//if err != nil {
	//	srv.logger.Error("Note not found in database", zap.Error(err))
	//	return nil, status.Error(codes.NotFound, "could not update note")
	//}
	//if token.UserID.String() != note.AuthorId {
	//	return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	//}

	err = validators.ValidateInsertBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//token blablabla

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

//si c ta note seulement
func (srv *notesService) UpdateBlock(ctx context.Context, in *notespb.UpdateBlockRequest) (*emptypb.Empty, error) {
	err := validators.ValidateUpdateBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var block = models.Block{}
	err = fillBlockContent(&block, in.Block)
	if err != nil {
		srv.logger.Error("failed to update block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid content provided for block id : %s", in.Id)
	}

	srv.repoBlock.Update(ctx, in.Id, &models.Block{ID: in.Id, Type: uint32(in.Block.Type), Index: in.Index, Content: block.Content})
	return nil, nil
}

func (srv *notesService) DeleteBlock(ctx context.Context, in *notespb.DeleteBlockRequest) (*emptypb.Empty, error) {
	err := validators.ValidateDeleteBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = srv.repoBlock.DeleteBlock(ctx, in.Id)
	if err != nil {
		srv.logger.Error("block was not deleted : ", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not delete block")
	}

	return &emptypb.Empty{}, nil
}

func fillBlockContent(block *models.Block, blockRequest *notespb.Block) error {
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
	//case *notespb.Block_Image_:
	//	block.Image.caption = &op.Image.Caption
	//	block.Image.url = &op.Image.Url
	//case *notespb.Block_Code_:
	//	block.Code.lang = &op.Code.Lang
	//	block.Code.Snippet = &op.Code.Snippet

	default:
		fmt.Println("No Data in this block")
		return status.Error(codes.Internal, "no data in this block")
	}
	return nil
}

func fillContentFromModelToApi(blockRequest *models.Block, contentType uint32, blockApi *notespb.Block) error {
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
	//case 6:
	//	(*blockApi).Data = &notespb.Block_Image_{Image: {caption: blockRequest.Image.caption, url: blockRequest.Image.url}}
	//case 7:
	//	(*blockApi).Data = &notespb.Block_Code_{Code: {sinppet: blockRequest.Code.Snippet, lang: blockRequest.Code.Lang}}
	default:
		fmt.Println("No such content in this block")
		return status.Errorf(codes.Internal, "no such content in this block")
	}
	return nil
}

func (srv *blocksService) authenticate(ctx context.Context) (*auth.Token, error) {
	token, err := srv.auth.TokenFromContext(ctx)
	if err != nil {
		srv.logger.Debug("failed to authenticate request", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return token, nil
}
*/
