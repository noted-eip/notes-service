package main

import (
	"context"
	"errors"
	"notes-service/auth"
	"notes-service/language"
	"notes-service/models"
	"notes-service/models/mongo"
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"time"

	"testing"

	"github.com/jaevor/go-nanoid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func statusFromModelError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, models.ErrNotFound) {
		return status.Error(codes.NotFound, "not found")
	}
	if errors.Is(err, models.ErrAlreadyExists) {
		return status.Error(codes.AlreadyExists, "already exists")
	}
	if errors.Is(err, models.ErrUnknown) {
		return status.Error(codes.Unknown, "unknown error")
	}
	if errors.Is(err, models.ErrForbidden) {
		return status.Error(codes.PermissionDenied, "forbidden operation")
	}
	return status.Error(codes.Internal, "internal error")
}

type testUtils struct {
	logger           *zap.Logger
	auth             *auth.TestService
	db               *mongo.Database
	notesRepository  models.NotesRepository
	groupsRepository models.GroupsRepository
	notes            notesv1.NotesAPIServer
	groups           notesv1.GroupsAPIServer
	newUUID          func() string
}

func newTestUtilsOrDie(t *testing.T) *testUtils {
	// logger, err := zap.NewDevelopment()
	// require.NoError(t, err)
	logger := zap.NewNop()
	auth := &auth.TestService{}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	db, err := mongo.NewDatabase(ctx, "mongodb://localhost:27017", "notes-service-unit-test", logger)
	if err != nil {
		t.Skip("skipping test, unable to connect to mongodb")
	}
	notesRepository := mongo.NewNotesRepository(db.DB, logger)
	groupsRepository := mongo.NewGroupsRepository(db.DB, logger)
	language := &language.NaturalAPIService{}
	require.NoError(t, language.Init())
	newUUID, err := nanoid.Standard(21)
	require.NoError(t, err)

	return &testUtils{
		logger:           logger,
		auth:             auth,
		db:               db,
		newUUID:          newUUID,
		notesRepository:  notesRepository,
		groupsRepository: groupsRepository,
		notes: &notesAPI{
			logger:   logger,
			auth:     auth,
			language: language,
			notes:    notesRepository,
			groups:   groupsRepository,
		},
		groups: &groupsAPI{
			logger: logger,
			auth:   auth,
			notes:  notesRepository,
			groups: groupsRepository,
		},
	}
}

type testAccount struct {
	ID        string
	Context   context.Context
	Workspace *testGroup
}

type testGroup struct {
	ID string
}

type testNote struct {
	ID     string
	Group  *testGroup
	Author *testAccount
}

type testBlock struct {
	note *testNote
	ID   string
}

type testInvite struct {
	recipient *testAccount
	sender    *testAccount
	group     *testGroup
	ID        string
}

func (account *testAccount) SendInvite(t *testing.T, tu *testUtils, recipient *testAccount, group *testGroup) *testInvite {
	sendInviteRes, err := tu.groups.SendInvite(account.Context, &notesv1.SendInviteRequest{GroupId: group.ID, RecipientAccountId: recipient.ID})
	require.NoError(t, err)
	require.NotNil(t, sendInviteRes)

	return &testInvite{
		recipient: recipient,
		sender:    account,
		group:     group,
		ID:        sendInviteRes.Invite.Id,
	}
}

func (account *testAccount) AcceptInvite(t *testing.T, tu *testUtils, invite *testInvite) {
	sendInviteRes, err := tu.groups.AcceptInvite(account.Context, &notesv1.AcceptInviteRequest{InviteId: invite.ID, GroupId: invite.group.ID})
	require.NoError(t, err)
	require.NotNil(t, sendInviteRes)
}

func (note *testNote) InsertBlock(t *testing.T, tu *testUtils, block *notesv1.Block, index uint32) *testBlock {
	res, err := tu.notes.InsertBlock(note.Author.Context, &notesv1.InsertBlockRequest{
		GroupId: note.Group.ID,
		NoteId:  note.ID,
		Index:   index,
		Block:   block,
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	return &testBlock{
		note: note,
		ID:   res.Block.Id,
	}
}

func newTestAccount(t *testing.T, tu *testUtils) *testAccount {
	aid := tu.newUUID()
	ctx, err := tu.auth.ContextWithToken(context.TODO(), &auth.Token{AccountID: aid})
	require.NoError(t, err)
	return &testAccount{
		ID:        aid,
		Context:   ctx,
		Workspace: nil,
	}
}

func newTestGroup(t *testing.T, tu *testUtils, owner *testAccount, members ...*testAccount) *testGroup {
	res, err := tu.groups.CreateGroup(owner.Context, &notesv1.CreateGroupRequest{
		Name:        "Some Random Name",
		Description: "Some Random Description",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	for i := range members {
		sendInvite, err := tu.groups.SendInvite(owner.Context, &notesv1.SendInviteRequest{
			GroupId:            res.Group.Id,
			RecipientAccountId: members[i].ID,
		})
		require.NoError(t, err)
		require.NotNil(t, sendInvite)
		acceptInvite, err := tu.groups.AcceptInvite(members[i].Context, &notesv1.AcceptInviteRequest{
			GroupId:  res.Group.Id,
			InviteId: sendInvite.Invite.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, acceptInvite)
	}

	return &testGroup{
		ID: res.Group.Id,
	}
}

func newTestWorkspace(t *testing.T, tu *testUtils, accountId string) *testGroup {
	res, err := tu.groups.CreateWorkspace(context.TODO(), &notesv1.CreateWorkspaceRequest{AccountId: accountId})
	require.NoError(t, err)
	require.NotNil(t, res)

	return &testGroup{
		ID: res.Group.Id,
	}
}

func newTestNote(t *testing.T, tu *testUtils, group *testGroup, author *testAccount, blocks []*notesv1.Block) *testNote {
	res, err := tu.notes.CreateNote(author.Context, &notesv1.CreateNoteRequest{
		GroupId: group.ID,
		Title:   "Default Title",
		Blocks:  blocks,
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	return &testNote{
		ID:     res.Note.Id,
		Author: author,
		Group:  group,
	}
}

func requireErrorHasGRPCCode(t *testing.T, code codes.Code, err error) {
	s, ok := status.FromError(err)
	require.True(t, ok, "expected grpc code %v got non-grpc error code", code)
	require.Equal(t, code, s.Code(), "expected grpc code %v got %v: %v", code, s.Code(), err)
}

func listOptionsFromLimitOffset(limit int32, offset int32) *models.ListOptions {
	if limit == 0 {
		limit = 20
	}
	return &models.ListOptions{
		Limit:  int64(limit),
		Offset: int64(offset),
	}
}
