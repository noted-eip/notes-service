package main

import (
	"context"
	"notes-service/models"
	v1 "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestInvitesSuite(t *testing.T) {
	tu := newTestUtilsOrDie(t)
	stranger := newTestAccount(t, tu)
	member := newTestAccount(t, tu)
	kerchak := newTestAccount(t, tu)
	jhon := newTestAccount(t, tu)
	dave := newTestAccount(t, tu)
	kerchakGroup := newTestGroup(t, tu, kerchak, member)

	// TODO: t.Run("stranger-cannot-send-invite", func(t *testing.T) {})

	var protoInviteSlot *v1.GroupInvite

	t.Run("member-can-send-invite", func(t *testing.T) {
		res, err := tu.groups.SendInvite(kerchak.Context, &v1.SendInviteRequest{GroupId: kerchakGroup.ID, RecipientAccountId: dave.ID})
		require.NoError(t, err)
		require.Equal(t, res.Invite.GroupId, kerchakGroup.ID)
		require.Equal(t, res.Invite.SenderAccountId, kerchak.ID)
		require.Equal(t, res.Invite.RecipientAccountId, dave.ID)

		// Check invite is stored in the database.
		group, err := tu.groupsRepository.GetGroupInternal(context.Background(), &models.OneGroupFilter{GroupID: kerchakGroup.ID})
		require.NoError(t, err)
		require.NotNil(t, group.FindInvite(res.Invite.Id))
		require.NotNil(t, group.FindInviteByAccountTuple(dave.ID, kerchak.ID))

		protoInviteSlot = res.Invite
	})

	// TODO: t.Run("member-cannot-send-invite-to-self", func(t *testing.T) {})

	// TODO: t.Run("member-cannot-send-invite-to-member", func(t *testing.T) {})

	// TODO: t.Run("member-cannot-send-duplicate-invite", func(t *testing.T) {})

	t.Run("stranger-cannot-get-invite", func(t *testing.T) {
		res, err := tu.groups.GetInvite(stranger.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: protoInviteSlot.Id})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("recipient-can-get-invite", func(t *testing.T) {
		res, err := tu.groups.GetInvite(dave.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: protoInviteSlot.Id})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("sender-can-get-invite", func(t *testing.T) {
		res, err := tu.groups.GetInvite(kerchak.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: protoInviteSlot.Id})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("member-can-get-invite", func(t *testing.T) {
		res, err := tu.groups.GetInvite(member.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: protoInviteSlot.Id})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("sender-cannot-deny-invite", func(t *testing.T) {
		res, err := tu.groups.DenyInvite(kerchak.Context, &v1.DenyInviteRequest{GroupId: kerchakGroup.ID, InviteId: protoInviteSlot.Id})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	// TODO: t.Run("stranger-cannot-deny-invite", func(t *testing.T) {})

	// TODO: t.Run("member-cannot-deny-invite", func(t *testing.T) {})

	t.Run("recipient-can-deny-invite", func(t *testing.T) {
		res, err := tu.groups.DenyInvite(dave.Context, &v1.DenyInviteRequest{GroupId: kerchakGroup.ID, InviteId: protoInviteSlot.Id})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Check invite is deleted from the database.
		group, err := tu.groupsRepository.GetGroupInternal(context.Background(), &models.OneGroupFilter{GroupID: kerchakGroup.ID})
		require.NoError(t, err)
		require.Nil(t, group.FindInvite(protoInviteSlot.Id))
		require.Nil(t, group.FindInviteByAccountTuple(dave.ID, kerchak.ID))

		// Check recipient has no access to group.
		getGroupRes, err := tu.groups.GetGroup(dave.Context, &v1.GetGroupRequest{GroupId: kerchakGroup.ID})
		requireErrorHasGRPCCode(t, codes.PermissionDenied, err)
		require.Nil(t, getGroupRes)
	})

	kerchakDaveInvite := kerchak.SendInvite(t, tu, dave, kerchakGroup)

	// TODO: t.Run("stranger-cannot-accept-invite", func(t *testing.T) {})

	// TODO: t.Run("sender-cannot-accept-invite", func(t *testing.T) {})

	// TODO: t.Run("member-cannot-accept-invite", func(t *testing.T) {}})

	t.Run("recipient-can-accept-invite", func(t *testing.T) {
		res, err := tu.groups.AcceptInvite(dave.Context, &v1.AcceptInviteRequest{GroupId: kerchakGroup.ID, InviteId: kerchakDaveInvite.ID})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Check invite is deleted from the database and recipient has been added to group.
		group, err := tu.groupsRepository.GetGroupInternal(context.Background(), &models.OneGroupFilter{GroupID: kerchakGroup.ID})
		require.NoError(t, err)
		require.Nil(t, group.FindInvite(kerchakDaveInvite.ID))
		require.Nil(t, group.FindInviteByAccountTuple(dave.ID, kerchak.ID))
		require.NotNil(t, group.FindMember(dave.ID))
		require.Len(t, *group.Members, 3)
	})

	kerchakJhonInvite := kerchak.SendInvite(t, tu, jhon, kerchakGroup)
	daveJhonInvite := dave.SendInvite(t, tu, jhon, kerchakGroup)

	t.Run("accept-invite-deletes-all-invites-destined-to-recipient", func(t *testing.T) {
		res, err := tu.groups.AcceptInvite(jhon.Context, &v1.AcceptInviteRequest{GroupId: kerchakGroup.ID, InviteId: daveJhonInvite.ID})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Check all invites are deleted from the database.
		group, err := tu.groupsRepository.GetGroupInternal(context.Background(), &models.OneGroupFilter{GroupID: kerchakGroup.ID})
		require.NoError(t, err)
		require.Nil(t, group.FindInvite(kerchakJhonInvite.ID))
		require.Nil(t, group.FindInviteByAccountTuple(kerchak.ID, jhon.ID))
		require.Nil(t, group.FindInvite(daveJhonInvite.ID))
		require.Nil(t, group.FindInviteByAccountTuple(dave.ID, jhon.ID))
	})

	// TODO: t.Run("list-invites-is-correctly-paginated", func(t *testing.T) {})

	// TODO: t.Run("user-can-list-invites-they-sent", func(t *testing.T) {})

	randomUserOne := newTestAccount(t, tu)
	randomUserTwo := newTestAccount(t, tu)
	randomUserThree := newTestAccount(t, tu)
	testAccountSlots := [3]*testAccount{randomUserOne, randomUserTwo, randomUserThree}
	for _, account := range testAccountSlots {
		dave.SendInvite(t, tu, account, kerchakGroup)
	}

	t.Run("member-can-list-invites-they-sent-in-group", func(t *testing.T) {
		res, err := tu.groups.ListInvites(dave.Context, &v1.ListInvitesRequest{SenderAccountId: dave.ID, GroupId: kerchakGroup.ID})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Check all invites are returned.
		require.Equal(t, len(res.Invites), 3)
		for _, invite := range res.Invites {
			require.Equal(t, invite.SenderAccountId, dave.ID)
			require.Equal(t, kerchakGroup.ID, invite.GroupId)
		}
	})

	t.Run("user-can-list-invites-destined-to-them", func(t *testing.T) {
		res, err := tu.groups.ListInvites(randomUserOne.Context, &v1.ListInvitesRequest{RecipientAccountId: randomUserOne.ID})
		require.NoError(t, err)
		require.Equal(t, len(res.Invites), 1)
		require.Equal(t, res.Invites[0].SenderAccountId, dave.ID)
		require.Equal(t, res.Invites[0].RecipientAccountId, randomUserOne.ID)
		require.Equal(t, res.Invites[0].GroupId, kerchakGroup.ID)
	})

	// TODO: t.Run("user-can-list-invites-destined-to-them-in-group", func(t *testing.T) {})

	t.Run("invitee-cannot-list-invites-destined-to-someone-else", func(t *testing.T) {
		res, err := tu.groups.ListInvites(randomUserTwo.Context, &v1.ListInvitesRequest{RecipientAccountId: randomUserOne.ID})
		requireErrorHasGRPCCode(t, codes.PermissionDenied, err)
		require.Nil(t, res)
	})

	// TODO: t.Run("user-can-list-invites-in-group", func(t *testing.T) {})

	// TODO: t.Run("stranger-cannot-revoke-invite", func(t *testing.T) {})

	// TODO: t.Run("recipient-cannot-revoke-invite", func(t *testing.T) {})

	maxime := newTestAccount(t, tu)
	kerchakMaximeInvite := dave.SendInvite(t, tu, maxime, kerchakGroup)

	t.Run("sender-can-revoke-invite", func(t *testing.T) {
		res, err := tu.groups.RevokeInvite(dave.Context, &v1.RevokeInviteRequest{
			GroupId:  kerchakMaximeInvite.group.ID,
			InviteId: kerchakMaximeInvite.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Check invite is deleted from the database and recipient has no access to group.
		group, err := tu.groupsRepository.GetGroupInternal(context.Background(), &models.OneGroupFilter{GroupID: kerchakGroup.ID})
		require.NoError(t, err)
		require.Nil(t, group.FindInvite(kerchakMaximeInvite.ID))
		require.Nil(t, group.FindInviteByAccountTuple(dave.ID, randomUserOne.ID))
		require.Nil(t, group.FindMember(randomUserOne.ID))
	})

	// diego := newTestAccount(t, tu)
	// kerchakDiegoInvite := dave.SendInvite(t, tu, diego, kerchakGroup)

	// t.Run("admin-can-revoke-invite", func(t *testing.T) {
	// 	res, err := tu.groups.RevokeInvite(kerchak.Context, &v1.RevokeInviteRequest{
	// 		GroupId:  kerchakGroup.ID,
	// 		InviteId: kerchakDiegoInvite.ID,
	// 	})
	// 	require.NoError(t, err)
	// 	require.NotNil(t, res)
	// })
}
