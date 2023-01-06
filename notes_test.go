package main

import (
	"context"
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
	auth auth.TestService
	srv  *notesService
}

func TestNotesService(t *testing.T) {
	suite.Run(t, new(NotesAPISuite))
}

func (s *NotesAPISuite) SetupSuite() {
	logger := newLoggerOrFail(s.T())
	dbNote := newDatabaseOrFail(s.T(), logger)
	dbBlock := newDatabaseOrFail(s.T(), logger)

	s.auth = auth.TestService{}
	s.srv = &notesService{
		auth:      &s.auth,
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
	//il faut implem les memory de blocks.go si on veux bien get la note
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

}*/

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

	//get all notes
	//delete all notes

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
