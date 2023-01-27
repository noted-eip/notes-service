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
	AccountID string
	Context   context.Context
}

func newTestAccount(t *testing.T, tu *testUtils) *testAccount {
	aid := tu.newUUID()
	ctx, err := tu.auth.ContextWithToken(context.TODO(), &auth.Token{AccountID: aid})
	require.NoError(t, err)
	return &testAccount{
		AccountID: aid,
		Context:   ctx,
	}
}

func requireErrorHasGRPCCode(t *testing.T, code codes.Code, err error) {
	s, ok := status.FromError(err)
	require.True(t, ok, "expected grpc code %v", code)
	require.Equal(t, code, s.Code(), "expected grpc code %v", code)
}
