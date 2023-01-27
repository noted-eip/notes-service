package main

import (
	"context"

	notesv1 "notes-service/protorepo/noted/notes/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *groupsAPI) GenerateInviteLink(ctx context.Context, req *notesv1.GenerateInviteLinkRequest) (*notesv1.GenerateInviteLinkResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) GetInviteLink(ctx context.Context, req *notesv1.GetInviteLinkRequest) (*notesv1.GetInviteLinkResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) RevokeInviteLink(ctx context.Context, req *notesv1.RevokeInviteLinkRequest) (*notesv1.RevokeInviteLinkResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) UseInviteLink(ctx context.Context, req *notesv1.UseInviteLinkRequest) (*notesv1.UseInviteLinkResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}
