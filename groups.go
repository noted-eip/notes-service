package main

import (
	"context"
	"notes-service/auth"
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type groupsAPI struct {
	notesv1.UnimplementedGroupsAPIServer

	logger *zap.Logger

	auth auth.Service

	groups models.GroupsRepository
	notes  models.NotesRepository
}

func (srv *groupsAPI) CreateGroup(ctx context.Context, req *notesv1.CreateGroupRequest) (*notesv1.CreateGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) CreateWorkspace(ctx context.Context, req *notesv1.CreateWorkspaceRequest) (*notesv1.CreateWorkspaceResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) GetGroup(ctx context.Context, req *notesv1.GetGroupRequest) (*notesv1.GetGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) UpdateGroup(ctx context.Context, req *notesv1.UpdateGroupRequest) (*notesv1.UpdateGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) DeleteGroup(ctx context.Context, req *notesv1.DeleteGroupRequest) (*notesv1.DeleteGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *groupsAPI) ListGroups(ctx context.Context, req *notesv1.ListGroupsRequest) (*notesv1.ListGroupsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}
