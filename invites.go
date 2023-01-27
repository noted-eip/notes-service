package main

import (
	"context"

	notesv1 "notes-service/protorepo/noted/notes/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *groupsAPI) SendInvite(ctx context.Context, req *notesv1.SendInviteRequest) (*notesv1.SendInviteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) GetInvite(ctx context.Context, req *notesv1.GetInviteRequest) (*notesv1.GetInviteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) AcceptInvite(ctx context.Context, req *notesv1.AcceptInviteRequest) (*notesv1.AcceptInviteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) DenyInvite(ctx context.Context, req *notesv1.DenyInviteRequest) (*notesv1.DenyInviteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) RevokeInvite(ctx context.Context, req *notesv1.RevokeInviteRequest) (*notesv1.RevokeInviteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) ListInvites(ctx context.Context, req *notesv1.ListInvitesRequest) (*notesv1.ListInvitesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}
