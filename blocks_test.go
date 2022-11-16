package main

import (
	"context"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBlocksServiceInsertBlock(t *testing.T) {
	srv := blocksService{}

	res, err := srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.Unimplemented)
	require.Nil(t, res)
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
