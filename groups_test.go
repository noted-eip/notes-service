package main

import (
	"context"
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestGroupsSuite(t *testing.T) {
	tu := newTestUtilsOrDie(t)
	jhon := newTestAccount(t, tu)
	dave := newTestAccount(t, tu)
	var jhonGroup *notesv1.Group
	var daveWorkspace *notesv1.Group

	t.Run("create-group", func(t *testing.T) {
		res, err := tu.groups.CreateGroup(jhon.Context, &notesv1.CreateGroupRequest{
			Name:        "My Group Name",
			Description: "My Group Description",
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		require.Equal(t, "My Group Name", res.Group.Name)
		require.Equal(t, "My Group Description", res.Group.Description)
		require.NotEmpty(t, res.Group.Id)

		// Check one admin member exists
		require.Len(t, res.Group.Members, 1)
		require.Equal(t, res.Group.Members[0].AccountId, jhon.ID)
		require.True(t, res.Group.Members[0].IsAdmin)

		// Check one conversation exists
		require.Len(t, res.Group.Conversations, 1)

		jhonGroup = res.Group
	})

	t.Run("create-workspace", func(t *testing.T) {
		res, err := tu.groups.CreateWorkspace(dave.Context, &notesv1.CreateWorkspaceRequest{})
		require.NoError(t, err)
		require.NotNil(t, res)

		require.Equal(t, res.Group.WorkspaceAccountId, dave.ID)

		// Workspace has no members, conversations or invites.
		require.Nil(t, res.Group.Members)
		require.Nil(t, res.Group.Conversations)
		require.Nil(t, res.Group.InviteLinks)
		require.Nil(t, res.Group.Invites)

		daveWorkspace = res.Group
	})

	t.Run("cannot-create-group-invalid-name", func(t *testing.T) {
		res, err := tu.groups.CreateGroup(jhon.Context, &notesv1.CreateGroupRequest{
			Name:        "",
			Description: "My Group Description",
		})
		requireErrorHasGRPCCode(t, codes.InvalidArgument, err)
		require.Nil(t, res)
	})

	t.Run("member-can-get-group", func(t *testing.T) {
		res, err := tu.groups.GetGroup(jhon.Context, &notesv1.GetGroupRequest{
			GroupId: jhonGroup.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, jhonGroup.Name, res.Group.Name)
	})

	t.Run("owner-can-get-workspace", func(t *testing.T) {
		res, err := tu.groups.GetGroup(dave.Context, &notesv1.GetGroupRequest{
			GroupId: daveWorkspace.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, daveWorkspace.Name, res.Group.Name)
	})

	t.Run("stranger-cannot-get-group", func(t *testing.T) {
		res, err := tu.groups.GetGroup(dave.Context, &notesv1.GetGroupRequest{
			GroupId: jhonGroup.Id,
		})
		require.Error(t, err)
		requireErrorHasGRPCCode(t, codes.PermissionDenied, err)
		require.Nil(t, res)
	})

	t.Run("admin-can-update-group-name", func(t *testing.T) {
		res, err := tu.groups.UpdateGroup(jhon.Context, &notesv1.UpdateGroupRequest{
			GroupId: jhonGroup.Id,
			Name:    "New Name",
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, res.Group.Name, "New Name")
		require.Equal(t, res.Group.Description, jhonGroup.Description)

		jhonGroup = res.Group
	})

	t.Run("admin-can-update-group-description", func(t *testing.T) {
		res, err := tu.groups.UpdateGroup(jhon.Context, &notesv1.UpdateGroupRequest{
			GroupId:     jhonGroup.Id,
			Description: "New Description",
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, res.Group.Description, "New Description")
		require.Equal(t, res.Group.Name, jhonGroup.Name)

		jhonGroup = res.Group
	})

	t.Run("stranger-cannot-update-group", func(t *testing.T) {
		res, err := tu.groups.UpdateGroup(dave.Context, &notesv1.UpdateGroupRequest{
			GroupId:     jhonGroup.Id,
			Description: "New Description",
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("member-can-list-group", func(t *testing.T) {
		res, err := tu.groups.ListGroups(jhon.Context, &notesv1.ListGroupsRequest{AccountId: jhon.ID})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Groups, 1)
	})

	t.Run("stranger-cannot-delete-group", func(t *testing.T) {
		res, err := tu.groups.DeleteGroup(dave.Context, &notesv1.DeleteGroupRequest{GroupId: jhonGroup.Id})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("admin-can-delete-group", func(t *testing.T) {
		res, err := tu.groups.DeleteGroup(jhon.Context, &notesv1.DeleteGroupRequest{GroupId: jhonGroup.Id})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Make sure group is deleted.
		group, err := tu.groupsRepository.GetGroupInternal(context.TODO(), &models.OneGroupFilter{GroupID: jhonGroup.Id})
		require.Error(t, err)
		require.Nil(t, group)
	})

	t.Run("owner-can-list-one-workspace", func(t *testing.T) {
		res, err := tu.groups.ListGroups(dave.Context, &notesv1.ListGroupsRequest{AccountId: dave.ID})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Groups, 1)
	})

	daveGroup := newTestGroup(t, tu, dave)
	dave.SendInvite(t, tu, jhon, daveGroup)

	t.Run("invitee-can-get-group-preview", func(t *testing.T) {
		res, err := tu.groups.GetGroup(jhon.Context, &notesv1.GetGroupRequest{
			GroupId: daveGroup.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, daveGroup.ID, res.Group.Id)
	})
}
