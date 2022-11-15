package main

import (
	"context"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBlocksServiceInsertBlockShouldReturnNil(t *testing.T) {
	srv := blocksService{}

	res, err := srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.InvalidArgument)
	require.Nil(t, res)
}

func TestBlocksServiceInsertBlockShouldReturnBlock(t *testing.T) {
	srv := blocksService{}

	res, err := srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{
		Block: &notespb.Block{
			Type: notespb.Block_TYPE_BULLET_POINT,
			Data: &notespb.Block_BulletPoint{},
		},
		Index:  1,
		NoteId: 1,
	})
	require.Nil(t, err)
	require.NotNil(t, res)
}

func TestBlocksServiceUpdateBlock(t *testing.T) {
	srv := blocksService{}

	res, err := srv.UpdateBlock(context.TODO(), &notespb.UpdateBlockRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.Unimplemented)
	require.Nil(t, res)
}

func TestBlocksServiceDeleteBlock(t *testing.T) {
	srv := blocksService{}

	res, err := srv.DeleteBlock(context.TODO(), &notespb.DeleteBlockRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.Unimplemented)
	require.Nil(t, res)
}
