package main

import (
	"context"
	"time"

	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"notes-service/validators"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (srv *groupsAPI) SendInvite(ctx context.Context, req *notesv1.SendInviteRequest) (*notesv1.SendInviteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateSendInviteRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	invite, err := srv.groups.SendInvite(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, &models.SendInvitePayload{
		RecipientAccountID: req.RecipientAccountId,
		ValidUntil:         time.Now().Add(time.Hour * 24 * 7),
	}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.SendInviteResponse{Invite: modelsInviteToProtobufInvite(invite, req.GroupId)}, nil
}

func (srv *groupsAPI) GetInvite(ctx context.Context, req *notesv1.GetInviteRequest) (*notesv1.GetInviteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) AcceptInvite(ctx context.Context, req *notesv1.AcceptInviteRequest) (*notesv1.AcceptInviteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateAcceptInviteRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	member, err := srv.groups.AcceptInvite(ctx, &models.OneInviteFilter{GroupID: req.GroupId, InviteID: req.InviteId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.AcceptInviteResponse{Member: modelsMemberToProtobufMember(member)}, nil
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

func modelsInviteToProtobufInvite(invite *models.GroupInvite, groupID string) *notesv1.GroupInvite {
	return &notesv1.GroupInvite{
		Id:                 invite.ID,
		GroupId:            groupID,
		SenderAccountId:    invite.SenderAccountID,
		RecipientAccountId: invite.RecipientAccountID,
		CreatedAt:          timestamppb.New(invite.CreatedAt),
		ValidUntil:         timestamppb.New(invite.ValidUntil),
	}
}
