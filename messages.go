package main

import (
	"context"

	notesv1 "notes-service/protorepo/noted/notes/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *groupsAPI) SendConversationMessage(ctx context.Context, req *notesv1.SendConversationMessageRequest) (*notesv1.SendConversationMessageResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) GetConversationMessage(ctx context.Context, req *notesv1.GetConversationMessageRequest) (*notesv1.GetConversationMessageResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) UpdateConversationMessage(ctx context.Context, req *notesv1.UpdateConversationMessageRequest) (*notesv1.UpdateConversationMessageResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) DeleteConversationMessage(ctx context.Context, req *notesv1.DeleteConversationMessageRequest) (*notesv1.DeleteConversationMessageResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) ListConversationMessages(ctx context.Context, req *notesv1.ListConversationMessagesRequest) (*notesv1.ListConversationMessagesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}
