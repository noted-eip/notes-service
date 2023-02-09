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

	var inviteSlot *v1.GroupInvite

	t.Run("send-and-get-invite", func(t *testing.T) {
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

		inviteSlot = getInviteRes.Invite
	})

	t.Run("get-invite-rights-check", func(t *testing.T) {
		// Random Dude should not have the rights
		res, err := tu.groups.GetInvite(randomDude.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: inviteSlot.Id})
		require.Error(t, err)
		require.Nil(t, res)

		// Kerchak should have the right, dave already has been tested
		res, err = tu.groups.GetInvite(kerchak.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: inviteSlot.Id})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("deny-invite-and-rights", func(t *testing.T) {
		res, err := tu.groups.DenyInvite(kerchak.Context, &v1.DenyInviteRequest{GroupId: kerchakGroup.ID, InviteId: inviteSlot.Id})
		require.Error(t, err)
		require.Nil(t, res)
		res, err = tu.groups.DenyInvite(randomDude.Context, &v1.DenyInviteRequest{GroupId: kerchakGroup.ID, InviteId: inviteSlot.Id})
		require.Error(t, err)
		require.Nil(t, res)

		res, err = tu.groups.DenyInvite(dave.Context, &v1.DenyInviteRequest{GroupId: kerchakGroup.ID, InviteId: inviteSlot.Id})
		require.NoError(t, err)
		require.NotNil(t, res)

		getInviteRes, err := tu.groups.GetInvite(dave.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: inviteSlot.Id})
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

	t.Run("list-invite-extensive", func(t *testing.T) {
		randomUserOne := newTestAccount(t, tu)
		randomUserTwo := newTestAccount(t, tu)
		randomUserThree := newTestAccount(t, tu)

		accountList := [3]*testAccount{randomUserOne, randomUserTwo, randomUserThree}
		var inviteList [3]*v1.GroupInvite

		// Send invites two the three new accounts
		for idx, account := range accountList {
			sendInviteRes, err := tu.groups.SendInvite(dave.Context, &v1.SendInviteRequest{GroupId: kerchakGroup.ID, RecipientAccountId: account.ID})
			require.NoError(t, err)
			require.NotNil(t, sendInviteRes)
			inviteList[idx] = sendInviteRes.Invite
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

		// Test on an invited account, basic check on the returned list (len and ids)
		invitesList, err = tu.groups.ListInvites(randomUserOne.Context, &v1.ListInvitesRequest{RecipientAccountId: randomUserOne.ID})
		require.NoError(t, err)
		require.Equal(t, len(invitesList.Invites), 1)
		require.Equal(t, invitesList.Invites[0].SenderAccountId, dave.ID)
		require.Equal(t, invitesList.Invites[0].RecipientAccountId, randomUserOne.ID)
		require.Equal(t, invitesList.Invites[0].GroupId, kerchakGroup.ID)

		// Test that you can't fetch another person invites
		invitesList, err = tu.groups.ListInvites(randomUserTwo.Context, &v1.ListInvitesRequest{RecipientAccountId: randomUserOne.ID})
		require.Error(t, err)
		require.Nil(t, invitesList)

		// Accept every invite and show that the list of dave's invites is now 0
		for idx, invite := range inviteList {
			curAccount := accountList[idx]

			acceptInviteRes, err := tu.groups.AcceptInvite(curAccount.Context, &v1.AcceptInviteRequest{GroupId: kerchakGroup.ID, InviteId: invite.Id})
			require.NoError(t, err)
			require.NotNil(t, acceptInviteRes)
		}

		// List from dave's context
		invitesList, err = tu.groups.ListInvites(dave.Context, &v1.ListInvitesRequest{SenderAccountId: dave.ID, GroupId: kerchakGroup.ID})
		require.NoError(t, err)
		require.NotNil(t, invitesList)
		require.Equal(t, len(invitesList.Invites), 0)
	})

	t.Run("revoke-invite", func(t *testing.T) {
		randomUserOne := newTestAccount(t, tu)
		sendInviteRes, err := tu.groups.SendInvite(dave.Context, &v1.SendInviteRequest{GroupId: kerchakGroup.ID, RecipientAccountId: randomUserOne.ID})
		require.NoError(t, err)
		require.NotNil(t, sendInviteRes)

		invite := sendInviteRes.Invite

		revokeRes, err := tu.groups.RevokeInvite(dave.Context, &v1.RevokeInviteRequest{
			GroupId:  kerchakGroup.ID,
			InviteId: invite.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, revokeRes)

		getInviteRes, err := tu.groups.GetInvite(dave.Context, &v1.GetInviteRequest{GroupId: kerchakGroup.ID, InviteId: invite.Id})
		require.Error(t, err)
		require.Nil(t, getInviteRes)
	})

}
