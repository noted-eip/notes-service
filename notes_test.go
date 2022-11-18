package main

import (
	"context"
	"crypto/ed25519"
	"notes-service/auth"
	"notes-service/memory"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/hashicorp/go-memdb"
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

func (s *NotesAPISuite) TestNotesServiceCreateNoteShouldReturnNil() {
	res, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{})

	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

func (s *NotesAPISuite) TestNotesServiceCreateNoteShouldReturnNote() {
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

func (s *NotesAPISuite) TestNotesServiceGetNoteShouldReturnError() {
	res, err := s.srv.GetNote(context.TODO(), &notespb.GetNoteRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

/*
func (s *NotesAPISuite) TestNotesServiceGetNoteShouldReturnNoError() {
	noteId, err := uuid.NewRandom()

	res, err := s.srv.GetNote(context.TODO(), &notespb.GetNoteRequest{
		Id: noteId.String(),
	})
	s.NotNil(res)
	s.Nil(err)
}
*/
func (s *NotesAPISuite) TestNotesServiceUpdateNoteShouldReturnError() {
	res, err := s.srv.UpdateNote(context.TODO(), &notespb.UpdateNoteRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

/*
func (s *NotesAPISuite) TestNotesServiceUpdateNoteShouldReturnNoError() {
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
func (s *NotesAPISuite) TestNotesServiceDeleteNoteShouldReturnError() {
	res, err := s.srv.DeleteNote(context.TODO(), &notespb.DeleteNoteRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

/*
func (s *NotesAPISuite) TestNotesServiceDeleteNoteShouldReturnNoError() {
	id, err := uuid.NewRandom()

	res, err := s.srv.DeleteNote(context.TODO(), &notespb.DeleteNoteRequest{
		Id: id.String(),
	})
	s.Nil(err)
	s.Nil(res)
}
*/
func (s *NotesAPISuite) TestNotesServiceListNotesShouldReturnError() {
	res, err := s.srv.ListNotes(context.TODO(), &notespb.ListNotesRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

/*
func (s *NotesAPISuite) TestNotesServiceListNotesShouldReturnNoError() {
	res, err := s.srv.ListNotes(context.TODO(), &notespb.ListNotesRequest{
		AuthorId: "author-id",
	})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}
*/

func newNotesDatabaseOrFail(t *testing.T, logger *zap.Logger) *memory.Database {
	db, err := memory.NewDatabase(context.Background(), newNotesDatabaseSchema(), logger)
	require.NoError(t, err, "could not instantiate in-memory database")
	return db
}

func newNotesDatabaseSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"note": {
				Name: "note",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"author_id": {
						Name:    "author_id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "AuthorId"},
					},
					"title": {
						Name:    "title",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Title"},
					},
					"blocks": {
						Name:    "blocks",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Blocks"},
					},
				},
			},
		},
	}
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
