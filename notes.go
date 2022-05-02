package main

import (
	"notes-service/grpc/notespb"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)


type notesService struct {
	notespb.UnimplementedNotesServiceServer
}

var _ notespb.NotesServiceServer = &notesService{}

func (srv *notesService) CreateNote(ctx context.Context, in *notespb.Note) 
(*emptypb.Empty, error)
{
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesService) GetNote(ctx context.Context, in *notespb.GetNoteRequest) 
(*notespb.Note, error)
{
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesService) UpdateNote(ctx context.Context, in *notespb.UpdateNoteRequest) 
(*emptypb.Empty, error)
{
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesService) DeleteNote(ctx context.Context, in *notespb.DeleteNoteRequest) 
(*emptypb.Empty, error)
{
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}