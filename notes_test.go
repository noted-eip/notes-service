package main

import (
	"context"
	"crypto/ed25519"
	"notes-service/auth"
	"notes-service/models/memory"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"

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

func TestNotesService(t *testing.T) {
	suite.Run(t, &NotesAPISuite{})
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

func (s *NotesAPISuite) CreateNoteShouldReturnNil() {
	res, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{})

	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

func (s *NotesAPISuite) CreateNoteShouldReturnNote() {
	res, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: "CI-TEST",
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Nil(err)
	s.NotNil(res)
}

func (s *NotesAPISuite) GetNoteShouldReturnError() {
	res, err := s.srv.GetNote(context.TODO(), &notespb.GetNoteRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

/*
func (s *NotesAPISuite) GetNoteShouldReturnNoError() {
	noteId, err := uuid.NewRandom()

	res, err := s.srv.GetNote(context.TODO(), &notespb.GetNoteRequest{
		Id: noteId.String(),
	})
	s.NotNil(res)
	s.Nil(err)
}
*/
func (s *NotesAPISuite) UpdateNoteShouldReturnError() {
	res, err := s.srv.UpdateNote(context.TODO(), &notespb.UpdateNoteRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

/*
func (s *NotesAPISuite) UpdateNoteShouldReturnNoError() {
	noteId, err := uuid.NewRandom()

	res, err := s.srv.UpdateNote(context.TODO(), &notespb.UpdateNoteRequest{
		Id: noteId.String(),
		Note: &notespb.Note{
			AuthorId: "CI-TEST",
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Nil(err)
	s.Nil(res)
}
*/
func (s *NotesAPISuite) DeleteNoteShouldReturnError() {
	res, err := s.srv.DeleteNote(context.TODO(), &notespb.DeleteNoteRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

/*
func (s *NotesAPISuite) DeleteNoteShouldReturnNoError() {
	id, err := uuid.NewRandom()

	res, err := s.srv.DeleteNote(context.TODO(), &notespb.DeleteNoteRequest{
		Id: id.String(),
	})
	s.Nil(err)
	s.Nil(res)
}
*/
func (s *NotesAPISuite) ListNotesShouldReturnError() {
	res, err := s.srv.ListNotes(context.TODO(), &notespb.ListNotesRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

/*
func (s *NotesAPISuite) ListNotesShouldReturnNoError() {
	res, err := s.srv.ListNotes(context.TODO(), &notespb.ListNotesRequest{
		AuthorId: "author-id",
	})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}
*/

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
