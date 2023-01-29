package main

import (
	"context"

	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (srv *groupsAPI) GetMember(ctx context.Context, req *notesv1.GetMemberRequest) (*notesv1.GetMemberResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) UpdateMember(ctx context.Context, req *notesv1.UpdateMemberRequest) (*notesv1.UpdateMemberResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) RemoveMember(ctx context.Context, req *notesv1.RemoveMemberRequest) (*notesv1.RemoveMemberResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func modelsMemberToProtobufMember(member *models.GroupMember) *notesv1.GroupMember {
	return &notesv1.GroupMember{
		AccountId: member.AccountID,
		IsAdmin:   member.IsAdmin,
		JoinedAt:  timestamppb.New(member.JoinedAt),
	}
}
