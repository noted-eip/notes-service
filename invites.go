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
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateGetInviteRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	invite, err := srv.groups.GetInvite(ctx, &models.OneInviteFilter{GroupID: req.GroupId, InviteID: req.InviteId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.GetInviteResponse{Invite: modelsInviteToProtobufInvite(invite, req.GroupId)}, nil
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

	srv.activities.CreateActivityInternal(ctx, &models.ActivityPayload{
		GroupID: req.GroupId,
		Type:    models.MemberJoined,
		Event:   "<userId:" + member.AccountID + "> joined the group <groupID:" + req.GroupId + ">.",
	})

	return &notesv1.AcceptInviteResponse{Member: modelsMemberToProtobufMember(member)}, nil
}

func (srv *groupsAPI) DenyInvite(ctx context.Context, req *notesv1.DenyInviteRequest) (*notesv1.DenyInviteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateDenyInviteRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = srv.groups.DenyInvite(ctx, &models.OneInviteFilter{GroupID: req.GroupId, InviteID: req.InviteId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.DenyInviteResponse{}, nil
}

func (srv *groupsAPI) RevokeInvite(ctx context.Context, req *notesv1.RevokeInviteRequest) (*notesv1.RevokeInviteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateRevokeInviteRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = srv.groups.RevokeGroupInvite(ctx, &models.OneInviteFilter{GroupID: req.GroupId, InviteID: req.InviteId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.RevokeInviteResponse{}, nil
}

func (srv *groupsAPI) ListInvites(ctx context.Context, req *notesv1.ListInvitesRequest) (*notesv1.ListInvitesResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateListInvitesRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if token.AccountID != req.RecipientAccountId && token.AccountID != req.SenderAccountId {
		if req.GroupId != "" {
			group, err := srv.groups.GetGroupInternal(ctx, &models.OneGroupFilter{GroupID: req.GroupId})

			if err != nil {
				return nil, err
			}
			member := group.FindMember(token.AccountID)
			if member == nil {
				return nil, status.Error(codes.PermissionDenied, "forbidden operation")
			}
		} else {
			return nil, status.Error(codes.PermissionDenied, "forbidden operation")
		}
	}

	invites, err := srv.groups.ListInvites(ctx,
		&models.ManyInvitesFilter{
			SenderAccountID:    req.SenderAccountId,
			RecipientAccountID: req.RecipientAccountId,
			GroupID:            req.GroupId,
		},
		listOptionsFromLimitOffset(req.Limit, req.Offset),
	)

	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.ListInvitesResponse{Invites: modelsListInviteResponseToProtobufInvites(invites)}, nil
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

func modelsListInviteResponseToProtobufInvites(invites []*models.ListInvitesResult) []*notesv1.GroupInvite {
	protoInvites := make([]*notesv1.GroupInvite, len(invites))

	for i := range invites {
		protoInvites[i] = modelsInviteToProtobufInvite(&invites[i].GroupInvite, invites[i].GroupID)
	}

	return protoInvites
}
