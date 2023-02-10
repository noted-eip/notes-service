package main

import (
	v1 "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvitesSuite(t *testing.T) {
	tu := newTestUtilsOrDie(t)
	kerchak := newTestAccount(t, tu)
	randomDude := newTestAccount(t, tu)
	dave := newTestAccount(t, tu)
	kerchakGroup := newTestGroup(t, tu, kerchak)

	var protoInviteSlot *v1.GroupInvite

	var testAccountSlots [3]*testAccount
	var testInviteSlots [3]*testInvite

	t.Run("send-invite", func(t *testing.T) {
		sendInviteRes, err := tu.groups.SendInvite(kerchak.Context, &v1.SendInviteRequest{GroupId: kerchakGroup.ID, RecipientAccountId: dave.ID})
		require.NoError(t, err)
		require.Equal(t, sendInviteRes.Invite.GroupId, kerchakGroup.ID)
		require.Equal(t, sendInviteRes.Invite.SenderAccountId, kerchak.ID)
		require.Equal(t, sendInviteRes.Invite.RecipientAccountId, dave.ID)

		getInviteRes, err := tu.groups.GetInvite(dave.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: sendInviteRes.Invite.Id})
		require.NoError(t, err)
		require.Equal(t, getInviteRes.Invite.GroupId, kerchakGroup.ID)
		require.Equal(t, getInviteRes.Invite.SenderAccountId, kerchak.ID)
		require.Equal(t, getInviteRes.Invite.RecipientAccountId, dave.ID)

		protoInviteSlot = getInviteRes.Invite
	})

	t.Run("get-invite-forbidden", func(t *testing.T) {
		// Random Dude should not have the rights
		res, err := tu.groups.GetInvite(randomDude.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: protoInviteSlot.Id})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("get-invite", func(t *testing.T) {
		// Kerchak should have the right, dave already has been tested
		res, err := tu.groups.GetInvite(kerchak.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: protoInviteSlot.Id})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("deny-invite-not-recipient", func(t *testing.T) {
		res, err := tu.groups.DenyInvite(kerchak.Context, &v1.DenyInviteRequest{GroupId: kerchakGroup.ID, InviteId: protoInviteSlot.Id})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("deny-invite", func(t *testing.T) {
		res, err := tu.groups.DenyInvite(dave.Context, &v1.DenyInviteRequest{GroupId: kerchakGroup.ID, InviteId: protoInviteSlot.Id})
		require.NoError(t, err)
		require.NotNil(t, res)

		getInviteRes, err := tu.groups.GetInvite(dave.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: protoInviteSlot.Id})
		require.Error(t, err)
		require.Nil(t, getInviteRes)

		getGroupRes, err := tu.groups.GetGroup(dave.Context, &v1.GetGroupRequest{GroupId: kerchak.ID})
		require.Error(t, err)
		require.Nil(t, getGroupRes)
	})

	t.Run("accept-invite", func(t *testing.T) {
		sendInviteRes, err := tu.groups.SendInvite(kerchak.Context, &v1.SendInviteRequest{GroupId: kerchakGroup.ID, RecipientAccountId: dave.ID})
		require.NoError(t, err)
		require.NotNil(t, sendInviteRes)

		acceptRes, err := tu.groups.AcceptInvite(dave.Context, &v1.AcceptInviteRequest{GroupId: kerchakGroup.ID, InviteId: sendInviteRes.Invite.Id})
		require.NoError(t, err)
		require.NotNil(t, acceptRes)

		getGroupRes, err := tu.groups.GetGroup(dave.Context, &v1.GetGroupRequest{GroupId: kerchakGroup.ID})
		require.NoError(t, err)
		require.NotNil(t, getGroupRes)
		require.Equal(t, len(getGroupRes.Group.Members), 2)
		require.Equal(t, getGroupRes.Group.Members[1].AccountId, dave.ID) // idk man not gonna do loops and shit
	})

	t.Run("list-invite-from-sender", func(t *testing.T) {
		randomUserOne := newTestAccount(t, tu)
		randomUserTwo := newTestAccount(t, tu)
		randomUserThree := newTestAccount(t, tu)

		testAccountSlots = [3]*testAccount{randomUserOne, randomUserTwo, randomUserThree}

		// Send invites two the three new accounts
		for idx, account := range testAccountSlots {
			testInviteSlots[idx] = dave.SendInvite(t, tu, account, kerchakGroup)
		}

		// List from dave's context
		invitesList, err := tu.groups.ListInvites(dave.Context, &v1.ListInvitesRequest{SenderAccountId: dave.ID, GroupId: kerchakGroup.ID})
		require.NoError(t, err)
		require.NotNil(t, invitesList)

		// Len and senderID check
		require.Equal(t, len(invitesList.Invites), 3)
		for _, invite := range invitesList.Invites {
			require.Equal(t, invite.SenderAccountId, dave.ID)
		}
	})

	t.Run("list-invite-from-recipient", func(t *testing.T) {
		randomUserOne := testAccountSlots[0]
		// Test on an invited account, basic check on the returned list (len and ids)
		invitesList, err := tu.groups.ListInvites(randomUserOne.Context, &v1.ListInvitesRequest{RecipientAccountId: randomUserOne.ID})
		require.NoError(t, err)
		require.Equal(t, len(invitesList.Invites), 1)
		require.Equal(t, invitesList.Invites[0].SenderAccountId, dave.ID)
		require.Equal(t, invitesList.Invites[0].RecipientAccountId, randomUserOne.ID)
		require.Equal(t, invitesList.Invites[0].GroupId, kerchakGroup.ID)

	})

	t.Run("list-invite-forbidden", func(t *testing.T) {
		randomUserOne := testAccountSlots[0]
		randomUserTwo := testAccountSlots[1]

		// Test that you can't fetch another person invites
		invitesList, err := tu.groups.ListInvites(randomUserTwo.Context, &v1.ListInvitesRequest{RecipientAccountId: randomUserOne.ID})
		require.Error(t, err)
		require.Nil(t, invitesList)
	})

	t.Run("list-invite-after-everybody-accepted", func(t *testing.T) {
		// Accept every invite and show that the list of dave's invites is now 0
		for idx, invite := range testInviteSlots {
			account := testAccountSlots[idx]
			account.AcceptInvite(t, tu, invite)
		}

		// List from dave's sender context
		invitesList, err := tu.groups.ListInvites(dave.Context, &v1.ListInvitesRequest{SenderAccountId: dave.ID, GroupId: kerchakGroup.ID})
		require.NoError(t, err)
		require.NotNil(t, invitesList)
		require.Equal(t, len(invitesList.Invites), 0)
	})

	t.Run("revoke-invite", func(t *testing.T) {
		randomUserOne := newTestAccount(t, tu)

		invite := dave.SendInvite(t, tu, randomUserOne, kerchakGroup)

		revokeRes, err := tu.groups.RevokeInvite(dave.Context, &v1.RevokeInviteRequest{
			GroupId:  invite.group.ID,
			InviteId: invite.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, revokeRes)

		getInviteRes, err := tu.groups.GetInvite(dave.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: invite.ID})
		require.Error(t, err)
		require.Nil(t, getInviteRes)
	})

}
