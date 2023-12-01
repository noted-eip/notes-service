package main

import (
	"context"

	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"notes-service/validators"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (srv *groupsAPI) GetMember(ctx context.Context, req *notesv1.GetMemberRequest) (*notesv1.GetMemberResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateGetMemberRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	group, err := srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.GetMemberResponse{Member: modelsMemberToProtobufMember(group.FindMember(req.AccountId))}, nil
}

func (srv *groupsAPI) UpdateMember(ctx context.Context, req *notesv1.UpdateMemberRequest) (*notesv1.UpdateMemberResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateUpdateMemberRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: We're not making use of the req.UpdateMask because for now there's only
	// one field you can update.

	member, err := srv.groups.UpdateGroupMember(ctx,
		&models.OneMemberFilter{GroupID: req.GroupId, AccountID: req.AccountId},
		&models.UpdateMemberPayload{IsAdmin: &req.Member.IsAdmin},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.UpdateMemberResponse{Member: modelsMemberToProtobufMember(member)}, nil
}

func (srv *groupsAPI) RemoveMember(ctx context.Context, req *notesv1.RemoveMemberRequest) (*notesv1.RemoveMemberResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateRemoveMemberRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = srv.moveNotesToUserWorkspaceOrDeleteThem(ctx,
		&models.ManyNotesFilter{AuthorAccountID: req.AccountId, GroupID: req.GroupId},
	)
	if err != nil {
		return nil, err
	}

	err = srv.groups.RemoveGroupMember(ctx,
		&models.OneMemberFilter{GroupID: req.GroupId, AccountID: req.AccountId},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	err = srv.notes.RemoveEditPermissions(ctx, &models.OneNoteFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	srv.activities.CreateActivityInternal(ctx, &models.ActivityPayload{
		GroupID: req.GroupId,
		Type:    models.MemberRemoved,
		Event:   "<userID:" + req.AccountId + "> leaved the group <groupID:" + req.GroupId + ">.",
	})

	return &notesv1.RemoveMemberResponse{}, nil
}

func (srv *groupsAPI) TrackScore(ctx context.Context, req *notesv1.TrackScoreRequest) (*notesv1.TrackScoreResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateTrackScoreRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	group, err := srv.groups.GetGroupInternal(ctx, &models.OneGroupFilter{GroupID: req.GroupId})
	if err != nil {
		return nil, statusFromModelError(err)
	}

	member := group.FindMember(token.AccountID)
	if member == nil {
		return nil, status.Error(codes.Unauthenticated, "member not found")
	}

	memberScore := member.Score + int(req.Score)
	memberResponses := member.QuizTotal + int(req.Responses)

	_, err = srv.groups.UpdateGroupMemberScore(ctx, &models.OneMemberFilter{GroupID: req.GroupId, AccountID: token.AccountID}, &models.UpdateMemberScorePayload{
		Score:      &memberScore,
		ScoreTotal: &memberResponses,
	}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.TrackScoreResponse{
		AccountId:  token.AccountID,
		GroupId:    req.GroupId,
		ScoreTotal: int32(memberScore),
		Responses:  int32(memberResponses),
	}, nil
}

func modelsMemberToProtobufMember(member *models.GroupMember) *notesv1.GroupMember {
	if member == nil {
		return nil
	}
	return &notesv1.GroupMember{
		AccountId: member.AccountID,
		IsAdmin:   member.IsAdmin,
		JoinedAt:  timestamppb.New(member.JoinedAt),
	}
}
