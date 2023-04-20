package main

import (
	"context"
	"notes-service/models"
	v1 "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinkInvitesSuite(t *testing.T) {
	tu := newTestUtilsOrDie(t)
	stranger := newTestAccount(t, tu)
	member := newTestAccount(t, tu)
	kerchak := newTestAccount(t, tu)
	jhon := newTestAccount(t, tu)
	kerchakGroup := newTestGroup(t, tu, kerchak, member)

	t.Run("stranger-cannot-generate-link-invite", func(t *testing.T) {
		res, err := tu.groups.GenerateInviteLink(stranger.Context, &v1.GenerateInviteLinkRequest{
			GroupId: kerchakGroup.ID,
		})

		require.Error(t, err)
		require.Nil(t, res)
	})

	var protoInviteLinkSlot *v1.GroupInviteLink

	t.Run("member-can-generate-invite", func(t *testing.T) {
		res, err := tu.groups.GenerateInviteLink(member.Context, &v1.GenerateInviteLinkRequest{
			GroupId: kerchakGroup.ID,
		})
		require.NoError(t, err)

		// Check invite is stored in the database.
		group, err := tu.groupsRepository.GetGroupInternal(context.Background(), &models.OneGroupFilter{GroupID: kerchakGroup.ID})
		require.NoError(t, err)
		require.NotNil(t, group.FindInviteLink(res.InviteLink.Code))

		protoInviteLinkSlot = res.InviteLink
	})

	t.Run("member-cannot-generate-duplicate-invite-link", func(t *testing.T) {
		res, err := tu.groups.GenerateInviteLink(member.Context, &v1.GenerateInviteLinkRequest{
			GroupId: kerchakGroup.ID,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("non-creator-cannot-get-invite-link", func(t *testing.T) {
		res, err := tu.groups.GetInviteLink(kerchak.Context, &v1.GetInviteLinkRequest{
			GroupId:        kerchakGroup.ID,
			InviteLinkCode: protoInviteLinkSlot.Code,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("creator-can-get-invite-link", func(t *testing.T) {
		res, err := tu.groups.GetInviteLink(member.Context, &v1.GetInviteLinkRequest{
			GroupId:        kerchakGroup.ID,
			InviteLinkCode: protoInviteLinkSlot.Code,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("member-cant-use-invite-link", func(t *testing.T) {
		res, err := tu.groups.UseInviteLink(kerchak.Context, &v1.UseInviteLinkRequest{
			GroupId:        kerchakGroup.ID,
			InviteLinkCode: protoInviteLinkSlot.Code,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("non-member-can-use-invite-link", func(t *testing.T) {
		res, err := tu.groups.UseInviteLink(jhon.Context, &v1.UseInviteLinkRequest{
			GroupId:        kerchakGroup.ID,
			InviteLinkCode: protoInviteLinkSlot.Code,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		member, err := tu.groups.GetMember(kerchak.Context, &v1.GetMemberRequest{
			GroupId:   kerchakGroup.ID,
			AccountId: jhon.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, member)
	})

	t.Run("non-creator-can-t-revoke-invite-link", func(t *testing.T) {
		res, err := tu.groups.RevokeInviteLink(kerchak.Context, &v1.RevokeInviteLinkRequest{
			GroupId:        kerchakGroup.ID,
			InviteLinkCode: protoInviteLinkSlot.Code,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("creator-can-revoke-invite-link", func(t *testing.T) {
		res, err := tu.groups.RevokeInviteLink(member.Context, &v1.RevokeInviteLinkRequest{
			GroupId:        kerchakGroup.ID,
			InviteLinkCode: protoInviteLinkSlot.Code,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		iLink, err := tu.groups.GetInviteLink(member.Context,
			&v1.GetInviteLinkRequest{
				GroupId:        kerchakGroup.ID,
				InviteLinkCode: protoInviteLinkSlot.Code,
			})
		require.Error(t, err)
		require.Nil(t, iLink)

	})

}
