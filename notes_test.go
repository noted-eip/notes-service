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
	dbNote := newDatabaseOrFail(s.T(), logger)
	dbBlock := newDatabaseOrFail(s.T(), logger)

	s.srv = &notesService{
		auth:      auth.NewService(genKeyOrFail(s.T())),
		logger:    logger,
		repoNote:  memory.NewNotesRepository(dbNote, logger),
		repoBlock: memory.NewBlocksRepository(dbBlock, logger),
	}
}

func (s *NotesAPISuite) TestCreateNoteShouldReturnNil() {
	res, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{})

	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

func (s *NotesAPISuite) TestCreateNoteShouldReturnNote() {
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

func (s *NotesAPISuite) TestGetNoteShouldReturnError() {
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
func (s *NotesAPISuite) TestUpdateNoteShouldReturnError() {
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
func (s *NotesAPISuite) TestDeleteNoteShouldReturnError() {
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
func (s *NotesAPISuite) TestListNotesShouldReturnError() {
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

func newDatabaseOrFail(t *testing.T, logger *zap.Logger) *memory.Database {
	db, err := memory.NewDatabase(context.Background(), logger)
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
