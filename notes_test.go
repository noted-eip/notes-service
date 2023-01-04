package main

import (
	"context"
	"crypto/ed25519"
	"notes-service/auth"
	"notes-service/models/memory"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotesAPISuite struct {
	suite.Suite
	srv *notesService
}

// mock auth package
type service struct {
	service Service
	token   auth.Token
}
type Service interface {
	TokenFromContext(ctx context.Context) (*auth.Token, error)
}

func NewMockService() Service {
	userUuid, err := uuid.NewRandom()
	if err != nil {
		return nil
	}
	return &service{
		service: &service{},
		token: auth.Token{
			Role:   auth.RoleUser,
			UserID: userUuid,
		},
	}
}

func (srv *service) TokenFromContext(ctx context.Context) (*auth.Token, error) {
	return &srv.token, nil
}

// !mock auth package

func TestNotesService(t *testing.T) {
	suite.Run(t, new(NotesAPISuite))
}

func (s *NotesAPISuite) SetupSuite() {
	logger := newLoggerOrFail(s.T())
	dbNote := newNotesDatabaseOrFail(s.T(), logger)
	dbBlock := newBlocksDatabaseOrFail(s.T(), logger)

	s.srv = &notesService{
		auth:      auth.NewService(genKeyOrFail(s.T())),
		logger:    logger,
		repoNote:  memory.NewNotesRepository(dbNote, logger),
		repoBlock: memory.NewBlocksRepository(dbBlock, logger),
	}
}

func (s *NotesAPISuite) TestCreateNoteNoAuth() {
	res, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{})

	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestCreateNoteValidator() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	res, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{})

	s.Require().Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Nil(res)
	s.srv.auth = saveAuthPackage
}

func (s *NotesAPISuite) TestCreateNoteReturnNote() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	res, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: "CI-TEST",
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Require().NoError(err)
	s.NotNil(res)
	s.srv.auth = saveAuthPackage
}

func (s *NotesAPISuite) TestGetNoteNoAuth() {
	res, err := s.srv.GetNote(context.TODO(), &notespb.GetNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestGetNoteValidator() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	res, err := s.srv.GetNote(context.TODO(), &notespb.GetNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Nil(res)
	s.srv.auth = saveAuthPackage
}

func (s *NotesAPISuite) TestGetNoteShouldReturnNote() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()

	resExpected, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: "CI-TEST",
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	//il faut implem les memory de blocks.go si on veux bien get la note
	res, err := s.srv.GetNote(context.TODO(), &notespb.GetNoteRequest{
		Id: resExpected.Note.Id,
	})

	s.NotNil(res)
	s.Require().NoError(err)
	s.Equal(res.Note.Id, resExpected.Note.Id)
	s.Equal(res.Note.AuthorId, resExpected.Note.AuthorId)
	s.Equal(res.Note.Title, resExpected.Note.Title)
	s.srv.auth = saveAuthPackage
}

func (s *NotesAPISuite) TestUpdateNoteNoAuth() {
	res, err := s.srv.UpdateNote(context.TODO(), &notespb.UpdateNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestUpdateNoteValidator() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	res, err := s.srv.UpdateNote(context.TODO(), &notespb.UpdateNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Nil(res)
	s.srv.auth = saveAuthPackage
}

/*
func (s *NotesAPISuite) TestUpdateNoteShouldReturnNoError() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	//get userId
	token, err := s.srv.auth.TokenFromContext(context.TODO())
	s.Require().NoError(err)
	userId := token.UserID.String()

	resCreateNote, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: userId,
			Title:    "ci-test",
			Blocks:   nil,
		},
	})

	_, err := s.srv.UpdateNote(context.TODO(), &notespb.UpdateNoteRequest{
		Id: resCreateNote.Note.Id,
		Note: &notespb.Note{
			AuthorId: userId,
			Title:    "ci-test-uptated",
			Blocks:   nil,
		},
	})

	res, err := s.srv.GetNote(context.TODO(), &notespb.GetNoteRequest{
		Id: resCreateNote.Note.Id,
	})
	s.Require().NoError(err)

	s.Nil(res)
	s.Equal(res.Note.Title, "CI-TEST-UPDATED")
	s.srv.auth = saveAuthPackage
}*/

func (s *NotesAPISuite) TestDeleteNoteNoAuth() {
	res, err := s.srv.DeleteNote(context.TODO(), &notespb.DeleteNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestDeleteNoteValidator() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	res, err := s.srv.DeleteNote(context.TODO(), &notespb.DeleteNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Nil(res)
	s.srv.auth = saveAuthPackage
}

func (s *NotesAPISuite) TestDeleteNoteShouldReturnNoError() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	//get userId
	token, err := s.srv.auth.TokenFromContext(context.TODO())
	s.Require().NoError(err)
	userId := token.UserID.String()

	resCreateNote, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: userId,
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Require().NoError(err)
	res, err := s.srv.DeleteNote(context.TODO(), &notespb.DeleteNoteRequest{
		Id: resCreateNote.Note.Id,
	})
	s.Require().NoError(err)
	s.Nil(res)
	s.srv.auth = saveAuthPackage
}

func (s *NotesAPISuite) TestListNotesNoAuth() {
	res, err := s.srv.ListNotes(context.TODO(), &notespb.ListNotesRequest{})
	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestListNotesValidator() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	res, err := s.srv.ListNotes(context.TODO(), &notespb.ListNotesRequest{})
	s.Require().Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Nil(res)
	s.srv.auth = saveAuthPackage
}

func (s *NotesAPISuite) TestListNotesReturnNotes() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()

	authorId := "CI-TEST"
	noteName := "ci-test-"

	ctx := context.TODO()

	_, err := s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{
		Note: &notespb.Note{AuthorId: authorId, Title: (noteName + "1"), Blocks: nil},
	})
	s.Require().NoError(err)

	_, err = s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{
		Note: &notespb.Note{AuthorId: authorId, Title: (noteName + "2"), Blocks: nil},
	})
	s.Require().NoError(err)

	_, err = s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{
		Note: &notespb.Note{AuthorId: authorId, Title: (noteName + "3"), Blocks: nil},
	})
	s.Require().NoError(err)

	res, err := s.srv.ListNotes(ctx, &notespb.ListNotesRequest{
		AuthorId: authorId,
	})

	s.Require().NoError(err)
	s.NotNil(res)
	s.Equal(3, len(res.Notes))
	s.srv.auth = saveAuthPackage
}

func newNotesDatabaseOrFail(t *testing.T, logger *zap.Logger) *memory.Database {
	db, err := memory.NewDatabase(context.Background(), memory.NewNotesDatabaseSchema(), logger)
	require.NoError(t, err, "could not instantiate in-memory database")
	return db
}

func newLoggerOrFail(t *testing.T) *zap.Logger {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(zapcore.FatalLevel), zap.WithCaller(false))
	require.NoError(t, err, "could not instantiate zap logger")
	return logger
}

func genKeyOrFail(t *testing.T) ed25519.PublicKey {
	publ, _, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	return publ
}
