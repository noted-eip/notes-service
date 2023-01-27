package main

import (
	"context"

	notesv1 "notes-service/protorepo/noted/notes/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *groupsAPI) GetConversation(ctx context.Context, req *notesv1.GetConversationRequest) (*notesv1.GetConversationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) UpdateConversation(ctx context.Context, req *notesv1.UpdateConversationRequest) (*notesv1.UpdateConversationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}
