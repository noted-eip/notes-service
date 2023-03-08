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
		res, err := tu.groups.CreateWorkspace(context.TODO(), &notesv1.CreateWorkspaceRequest{AccountId: dave.ID})
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

	// We tested and it works, so let's create one for jhon
	jhon.Workspace = newTestWorkspace(t, tu, jhon.ID)

	t.Run("delete-group-should-move-notes-to-workspace", func(t *testing.T) {
		jhonNewGroup := newTestGroup(t, tu, jhon)
		jhonFirstNote := newTestNote(t, tu, jhonNewGroup, jhon, []*notesv1.Block{})
		jhonSecondNote := newTestNote(t, tu, jhonNewGroup, jhon, []*notesv1.Block{})

		res, err := tu.groups.DeleteGroup(jhon.Context, &notesv1.DeleteGroupRequest{GroupId: jhonNewGroup.ID})
		require.NoError(t, err)
		require.NotNil(t, res)

		note, err := tu.notes.GetNote(jhon.Context, &notesv1.GetNoteRequest{NoteId: jhonFirstNote.ID, GroupId: jhon.Workspace.ID})
		require.NoError(t, err)
		require.NotNil(t, note)

		note, err = tu.notes.GetNote(jhon.Context, &notesv1.GetNoteRequest{NoteId: jhonSecondNote.ID, GroupId: jhon.Workspace.ID})
		require.NoError(t, err)
		require.NotNil(t, note)
	})

	t.Run("delete-group-should-delete-notes-if-no-workspace", func(t *testing.T) {
		jean := newTestAccount(t, tu)
		jeanGroup := newTestGroup(t, tu, jean, jhon, dave)
		_ = newTestNote(t, tu, jeanGroup, jean, []*notesv1.Block{})
		_ = newTestNote(t, tu, jeanGroup, jean, []*notesv1.Block{})

		res, err := tu.groups.DeleteGroup(jean.Context, &notesv1.DeleteGroupRequest{GroupId: jeanGroup.ID})
		require.NoError(t, err)
		require.NotNil(t, res)

		notes, err := tu.notesRepository.ListAllNotesInternal(context.TODO(), &models.ManyNotesFilter{
			AuthorAccountID: jean.ID,
		})
		require.NoError(t, err)
		require.Zero(t, len(notes))
	})

	t.Run("group-is-not-foundable-after-being-deleted", func(t *testing.T) {
		res, err := tu.groups.GetGroup(jhon.Context, &notesv1.GetGroupRequest{
			GroupId: jhonGroup.Id,
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("owner-can-list-one-workspace", func(t *testing.T) {
		res, err := tu.groups.ListGroups(dave.Context, &notesv1.ListGroupsRequest{AccountId: dave.ID})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Groups, 1)
	})

	// Test account on delete side effects
	daveGroup := newTestGroup(t, tu, dave)
	daveNote := newTestNote(t, tu, daveGroup, dave, nil)

	t.Run("delete-group-change-note-group-to-user-workspace", func(t *testing.T) {
		res, err := tu.groups.DeleteGroup(dave.Context, &notesv1.DeleteGroupRequest{GroupId: daveGroup.ID})
		require.NoError(t, err)
		require.NotNil(t, res)

		note, err := tu.notes.GetNote(dave.Context, &notesv1.GetNoteRequest{NoteId: daveNote.ID, GroupId: daveWorkspace.Id})
		require.NoError(t, err)
		require.NotNil(t, note)
	})

	daveGroup = newTestGroup(t, tu, dave)
	dave.SendInvite(t, tu, jhon, daveGroup)

	t.Run("invitee-can-get-group-preview", func(t *testing.T) {
		res, err := tu.groups.GetGroup(jhon.Context, &notesv1.GetGroupRequest{
			GroupId: daveGroup.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, daveGroup.ID, res.Group.Id)
	})

	// OnAccountDelete is a repository function, no auth
	// Env setting for OnAccountDelete tests:
	jean := newTestAccount(t, tu)
	bibi := newTestAccount(t, tu)
	bibiSister := newTestAccount(t, tu)
	bibiBrother := newTestAccount(t, tu)

	bibi.Workspace = newTestWorkspace(t, tu, bibi.ID)
	bibiGroup := newTestGroup(t, tu, bibi, jean)
	bibiSchoolGroup := newTestGroup(t, tu, bibi, jean)
	bibi.SendInvite(t, tu, bibiSister, bibiGroup)
	bibi.SendInvite(t, tu, bibiBrother, bibiGroup)

	t.Run("on-account-delete-should-return-no-error", func(t *testing.T) {
		// OnAccountDelete will only return mongodb errors
		err := tu.groupsRepository.OnAccountDelete(context.TODO(), bibi.ID)
		require.NoError(t, err)
	})

	t.Run("on-account-delete-should-delete-member-references", func(t *testing.T) {
		res, err := tu.groups.GetMember(jean.Context, &notesv1.GetMemberRequest{GroupId: bibiGroup.ID, AccountId: bibi.ID})
		require.NoError(t, err)
		require.Nil(t, res.Member)

		res, err = tu.groups.GetMember(jean.Context, &notesv1.GetMemberRequest{GroupId: bibiSchoolGroup.ID, AccountId: bibi.ID})
		require.NoError(t, err)
		require.Nil(t, res.Member)
	})

	t.Run("on-account-delete-should-delete-member-invites", func(t *testing.T) {
		res, err := tu.groups.ListInvites(bibi.Context, &notesv1.ListInvitesRequest{SenderAccountId: bibi.ID})
		require.NoError(t, err)
		require.Equal(t, len(res.Invites), 0)
	})

	t.Run("on-account-delete-should-delete-workspace", func(t *testing.T) {
		_, err := tu.groups.GetGroup(bibi.Context, &notesv1.GetGroupRequest{GroupId: bibi.Workspace.ID})
		require.Error(t, err)
	})

}
