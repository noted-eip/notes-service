package main

import (
	"context"
	"time"

	"notes-service/background"
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"notes-service/validators"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (srv *groupsAPI) GenerateInviteLink(ctx context.Context, req *notesv1.GenerateInviteLinkRequest) (*notesv1.GenerateInviteLinkResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateGenerateInviteLinkRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// NOTE: for now 1 week
	validUntil := time.Now().Add(time.Hour * 24 * 7)

	inviteLink, err := srv.groups.GenerateGroupInviteLink(ctx,
		&models.OneGroupFilter{GroupID: req.GroupId},
		&models.GenerateGroupInviteLinkPayload{
			GeneratedByAccountID: token.AccountID,
			ValidUntil:           validUntil,
		}, token.AccountID)
	if err != nil {
		return nil, err
	}

	// Process that will revoke the invite one week later
	err = srv.background.AddProcess(&background.Process{
		Identifier: &models.InviteLinkIdentifier{
			Code:    inviteLink.Code,
			GroupId: req.GroupId,
			Action:  models.InviteLinkRevoke,
		},
		CallBackFct: func() error {
			return srv.groups.RevokeInviteLink(
				ctx,
				&models.OneInviteLinkFilter{
					GroupID:        req.GroupId,
					InviteLinkCode: inviteLink.Code,
				},
				token.AccountID)
		},
		SecondsToDebounce:             uint32(validUntil.Second()),
		CancelProcessOnSameIdentifier: true,
		RepeatProcess:                 false,
	})
	if err != nil {
		return nil, err
	}

	return &notesv1.GenerateInviteLinkResponse{
		InviteLink: modelsInviteLinkToProtobufGroup(inviteLink),
	}, nil

}

func (srv *groupsAPI) GetInviteLink(ctx context.Context, req *notesv1.GetInviteLinkRequest) (*notesv1.GetInviteLinkResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateGetInviteLinkRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	inviteLink, err := srv.groups.GetInviteLink(ctx, &models.OneInviteLinkFilter{GroupID: req.GroupId, InviteLinkCode: req.InviteLinkCode}, token.AccountID)
	if err != nil {
		return nil, err
	}

	return &notesv1.GetInviteLinkResponse{
		InviteLink: modelsInviteLinkToProtobufGroup(inviteLink),
	}, nil
}

func (srv *groupsAPI) RevokeInviteLink(ctx context.Context, req *notesv1.RevokeInviteLinkRequest) (*notesv1.RevokeInviteLinkResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateRevokeInviteLinkRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = srv.groups.RevokeInviteLink(ctx, &models.OneInviteLinkFilter{GroupID: req.GroupId, InviteLinkCode: req.InviteLinkCode}, token.AccountID)
	if err != nil {
		return nil, err
	}

	err = srv.background.CancelProcess(
		&background.Process{
			Identifier: &models.InviteLinkIdentifier{
				Code:    req.InviteLinkCode,
				GroupId: req.GroupId,
				Action:  models.InviteLinkRevoke,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return &notesv1.RevokeInviteLinkResponse{}, nil
}

func (srv *groupsAPI) UseInviteLink(ctx context.Context, req *notesv1.UseInviteLinkRequest) (*notesv1.UseInviteLinkResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateUseInviteLinkRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = srv.groups.UseInviteLink(ctx, &models.OneInviteLinkFilter{
		GroupID:        req.GroupId,
		InviteLinkCode: req.InviteLinkCode,
	}, token.AccountID)
	if err != nil {
		return nil, err
	}

	return &notesv1.UseInviteLinkResponse{}, nil
}

func modelsInviteLinkToProtobufGroup(invite *models.GroupInviteLink) *notesv1.GroupInviteLink {
	return &notesv1.GroupInviteLink{
		Code:                 invite.Code,
		GeneratedByAccountId: invite.GeneratedByAccountID,
		CreatedAt:            timestamppb.New(invite.CreatedAt),
		ValidUntil:           timestamppb.New(invite.ValidUntil),
	}
}
