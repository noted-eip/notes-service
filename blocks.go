package main

import (
	"context"
	notesv1 "notes-service/protorepo/noted/notes/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *notesService) InsertBlock(ctx context.Context, in *notesv1.InsertBlockRequest) (*notesv1.InsertBlockResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesService) UpdateBlock(ctx context.Context, in *notesv1.UpdateBlockRequest) (*notesv1.UpdateBlockResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesService) DeleteBlock(ctx context.Context, in *notesv1.DeleteBlockRequest) (*notesv1.DeleteBlockResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}
