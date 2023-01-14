package main

import (
	"context"
	"notes-service/auth"
	"notes-service/language"
	"notes-service/models/mongo"
	notespb "notes-service/protorepo/noted/notes/v1"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotesAPISuite struct {
	suite.Suite
	auth auth.TestService
	srv  *notesService
}

func TestNotesAPI(t *testing.T) {
	if os.Getenv("NOTES_SERVICE_TEST_MONGODB_URI") == "" {
		t.Skipf("Skipping NotesAPI suite, missing NOTES_SERVICE_TEST_MONGODB_URI environment variable.")
		return
	}

	suite.Run(t, new(NotesAPISuite))
}

func (s *NotesAPISuite) SetupSuite() {
	db := newDatabaseOrFail(s.T())

	s.auth = auth.TestService{}
	s.srv = &notesService{
		auth:      &s.auth,
		logger:    zap.NewNop(),
		language:  &language.NaturalAPIService{},
		repoNote:  mongo.NewNotesRepository(db.DB, zap.NewNop()),
		repoBlock: mongo.NewBlocksRepository(db.DB, zap.NewNop()),
	}
	s.Require().NoError(s.srv.language.Init())
	db.Disconnect(context.TODO())
}

func (s *NotesAPISuite) TestCreateNoteNoAuth() {
	res, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{})

	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestCreateNoteValidator() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)
	res, err := s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{})

	s.Require().Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestCreateNoteReturnNote() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)
	res, err := s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: generatedUuid.String(),
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Require().NoError(err)
	s.NotNil(res)

}

func (s *NotesAPISuite) TestGetNoteNoAuth() {
	res, err := s.srv.GetNote(context.TODO(), &notespb.GetNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestGetNoteValidator() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)
	res, err := s.srv.GetNote(ctx, &notespb.GetNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Nil(res)

}

func (s *NotesAPISuite) TestGetNoteShouldReturnNote() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)

	resExpected, err := s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: generatedUuid.String(),
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Require().NoError(err)

	res, err := s.srv.GetNote(ctx, &notespb.GetNoteRequest{
		Id: resExpected.Note.Id,
	})

	s.NotNil(res)
	s.Require().NoError(err)
	s.Equal(res.Note.Id, resExpected.Note.Id)
	s.Equal(res.Note.AuthorId, resExpected.Note.AuthorId)
	s.Equal(res.Note.Title, resExpected.Note.Title)

}

func (s *NotesAPISuite) TestUpdateNoteNoAuth() {
	res, err := s.srv.UpdateNote(context.TODO(), &notespb.UpdateNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestUpdateNoteValidator() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)
	res, err := s.srv.UpdateNote(ctx, &notespb.UpdateNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestDeleteNoteNoAuth() {
	res, err := s.srv.DeleteNote(context.TODO(), &notespb.DeleteNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestDeleteNoteValidator() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)
	res, err := s.srv.DeleteNote(ctx, &notespb.DeleteNoteRequest{})
	s.Require().Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestDeleteNoteShouldReturnNoError() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)
	//get userId

	resCreateNote, err := s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: generatedUuid.String(),
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Require().NoError(err)
	res, err := s.srv.DeleteNote(ctx, &notespb.DeleteNoteRequest{
		Id: resCreateNote.Note.Id,
	})
	s.Require().NoError(err)
	s.Nil(res)
}

func (s *NotesAPISuite) TestListNotesNoAuth() {
	res, err := s.srv.ListNotes(context.TODO(), &notespb.ListNotesRequest{})
	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestListNotesValidator() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)
	res, err := s.srv.ListNotes(ctx, &notespb.ListNotesRequest{})
	s.Require().Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestListNotesReturnNotes() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)

	authorId := generatedUuid.String()
	noteName := "ci-test-"

	_, err = s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{
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
}

func newDatabaseOrFail(t *testing.T) *mongo.Database {
	db, err := mongo.NewDatabase(context.TODO(), os.Getenv("NOTES_SERVICE_TEST_MONGODB_URI"), "notes-service-test", zap.NewNop())
	require.NoError(t, err, "could not instantiate mongo database")
	return db
}
