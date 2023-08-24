package main

import (
	"context"
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestMembersSuite(t *testing.T) {
	tu := newTestUtilsOrDie(t)
	balthi := newTestAccount(t, tu)
	thomas := newTestAccount(t, tu)
	kilian := newTestAccount(t, tu)
	stranger := newTestAccount(t, tu)
	edouard := newTestAccount(t, tu)
	maxime := newTestAccount(t, tu)
	balthiGroup := newTestGroup(t, tu, balthi, thomas, edouard, kilian, maxime)

	t.Run("update-member-promote-to-admin", func(t *testing.T) {
		res, err := tu.groups.UpdateMember(balthi.Context, &notesv1.UpdateMemberRequest{
			GroupId:   balthiGroup.ID,
			AccountId: thomas.ID,
			Member: &notesv1.GroupMember{
				IsAdmin: true,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"is_admin"},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Ensure user is now admin.
		group, err := tu.groupsRepository.GetGroup(balthi.Context, &models.OneGroupFilter{GroupID: balthiGroup.ID}, balthi.ID)
		require.NoError(t, err)
		require.NotNil(t, group.FindMember(thomas.ID))
		require.True(t, group.FindMember(thomas.ID).IsAdmin)
	})

	t.Run("update-member-non-admin-cannot-promote", func(t *testing.T) {
		res, err := tu.groups.UpdateMember(edouard.Context, &notesv1.UpdateMemberRequest{
			GroupId:   balthiGroup.ID,
			AccountId: kilian.ID,
			Member: &notesv1.GroupMember{
				IsAdmin: true,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"is_admin"},
			},
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("update-member-stranger-cannot-promote", func(t *testing.T) {
		res, err := tu.groups.UpdateMember(stranger.Context, &notesv1.UpdateMemberRequest{
			GroupId:   balthiGroup.ID,
			AccountId: kilian.ID,
			Member: &notesv1.GroupMember{
				IsAdmin: true,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"is_admin"},
			},
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("update-member-stranger-cannot-promote-admin", func(t *testing.T) {
		res, err := tu.groups.UpdateMember(stranger.Context, &notesv1.UpdateMemberRequest{
			GroupId:   balthiGroup.ID,
			AccountId: balthi.ID,
			Member: &notesv1.GroupMember{
				IsAdmin: true,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"is_admin"},
			},
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("update-member-non-admin-cannot-promote-itself", func(t *testing.T) {
		res, err := tu.groups.UpdateMember(edouard.Context, &notesv1.UpdateMemberRequest{
			GroupId:   balthiGroup.ID,
			AccountId: edouard.ID,
			Member: &notesv1.GroupMember{
				IsAdmin: true,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"is_admin"},
			},
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("update-member-admin-cannot-promote-itself", func(t *testing.T) {
		res, err := tu.groups.UpdateMember(balthi.Context, &notesv1.UpdateMemberRequest{
			GroupId:   balthiGroup.ID,
			AccountId: balthi.ID,
			Member: &notesv1.GroupMember{
				IsAdmin: true,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"is_admin"},
			},
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("remove-member-admin-can-leave-group", func(t *testing.T) {
		res, err := tu.groups.RemoveMember(balthi.Context, &notesv1.RemoveMemberRequest{
			GroupId:   balthiGroup.ID,
			AccountId: balthi.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Ensure user has no access to group anymore.
		_, err = tu.groups.GetGroup(balthi.Context, &notesv1.GetGroupRequest{
			GroupId: balthiGroup.ID,
		})
		requireErrorHasGRPCCode(t, codes.PermissionDenied, err)
	})

	t.Run("remove-member-admin-can-remove-regular-user", func(t *testing.T) {
		res, err := tu.groups.RemoveMember(edouard.Context, &notesv1.RemoveMemberRequest{
			GroupId:   balthiGroup.ID,
			AccountId: edouard.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Ensure user has no access to group anymore.
		_, err = tu.groups.GetGroup(edouard.Context, &notesv1.GetGroupRequest{
			GroupId: balthiGroup.ID,
		})
		requireErrorHasGRPCCode(t, codes.PermissionDenied, err)
	})

	t.Run("remove-member-regular-member-cannot-remove-regular-user", func(t *testing.T) {
		res, err := tu.groups.RemoveMember(maxime.Context, &notesv1.RemoveMemberRequest{
			GroupId:   balthiGroup.ID,
			AccountId: kilian.ID,
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)

		// Ensure user has still access to group.
		_, err = tu.groups.GetGroup(kilian.Context, &notesv1.GetGroupRequest{
			GroupId: balthiGroup.ID,
		})
		require.NoError(t, err)
	})

	t.Run("remove-member-regular-member-can-remove-itself", func(t *testing.T) {
		res, err := tu.groups.RemoveMember(maxime.Context, &notesv1.RemoveMemberRequest{
			GroupId:   balthiGroup.ID,
			AccountId: maxime.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Ensure user has no access to group anymore.
		_, err = tu.groups.GetGroup(maxime.Context, &notesv1.GetGroupRequest{
			GroupId: balthiGroup.ID,
		})
		requireErrorHasGRPCCode(t, codes.PermissionDenied, err)
	})

	jhon := newTestAccount(t, tu)
	jhonGroup := newTestGroup(t, tu, jhon, balthi)

	t.Run("delete-group-should-move-notes-to-workspace", func(t *testing.T) {
		balthi.Workspace = newTestWorkspace(t, tu, balthi.ID)
		balthiFirstNote := newTestNote(t, tu, jhonGroup, balthi, []*notesv1.Block{})
		balthiSecondNote := newTestNote(t, tu, jhonGroup, balthi, []*notesv1.Block{})

		res, err := tu.groups.RemoveMember(jhon.Context, &notesv1.RemoveMemberRequest{
			GroupId:   jhonGroup.ID,
			AccountId: balthi.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		note, err := tu.notes.GetNote(balthi.Context, &notesv1.GetNoteRequest{NoteId: balthiFirstNote.ID, GroupId: balthi.Workspace.ID})
		require.NoError(t, err)
		require.NotNil(t, note)

		note, err = tu.notes.GetNote(balthi.Context, &notesv1.GetNoteRequest{NoteId: balthiSecondNote.ID, GroupId: balthi.Workspace.ID})
		require.NoError(t, err)
		require.NotNil(t, note)
	})

	t.Run("delete-group-should-delete-notes-if-no-workspace", func(t *testing.T) {
		jean := newTestAccount(t, tu)
		jean.AcceptInvite(t, tu, jhon.SendInvite(t, tu, jean, jhonGroup))

		_ = newTestNote(t, tu, jhonGroup, jean, []*notesv1.Block{})
		_ = newTestNote(t, tu, jhonGroup, jean, []*notesv1.Block{})

		res, err := tu.groups.RemoveMember(jhon.Context, &notesv1.RemoveMemberRequest{
			GroupId:   jhonGroup.ID,
			AccountId: jean.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		notes, err := tu.notesRepository.ListAllNotesInternal(context.TODO(), &models.ManyNotesFilter{
			AuthorAccountID: jean.ID,
		})
		require.NoError(t, err)
		require.Zero(t, len(notes))
	})

}
