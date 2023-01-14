package main

import (
	"context"
	"notes-service/auth"
	"notes-service/models/mongo"
	notespb "notes-service/protorepo/noted/notes/v1"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExportAPISuite struct {
	suite.Suite
	auth auth.TestService
	srv  *notesService
}

func TestExportAPI(t *testing.T) {
	if os.Getenv("NOTES_SERVICE_TEST_MONGODB_URI") == "" {
		t.Skipf("Skipping NotesAPI suite, missing NOTES_SERVICE_TEST_MONGODB_URI environment variable.")
		return
	}

	suite.Run(t, &ExportAPISuite{})
}

func (s *ExportAPISuite) SetupSuite() {
	db := newDatabaseOrFail(s.T())

	s.srv = &notesService{
		auth:      &s.auth,
		logger:    zap.NewNop(),
		repoNote:  mongo.NewNotesRepository(db.DB, zap.NewNop()),
		repoBlock: mongo.NewBlocksRepository(db.DB, zap.NewNop()),
	}
}

func (s *ExportAPISuite) TestExportWrongNoteIDShouldReturnAnError() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)

	generatedUuid, err = uuid.NewRandom()
	s.Require().NoError(err)
	res, err := s.srv.ExportNote(ctx, &notespb.ExportNoteRequest{NoteId: generatedUuid.String(), ExportFormat: notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_MARKDOWN})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.NotFound)
	s.Nil(res)
}

func (s *ExportAPISuite) TestUnauthenticatedShouldReturnAnError() {
	uuid, _ := uuid.NewRandom()
	res, err := s.srv.ExportNote(context.TODO(), &notespb.ExportNoteRequest{NoteId: uuid.String(), ExportFormat: notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_MARKDOWN})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.Unauthenticated)
	s.Nil(res)
}

func (s *ExportAPISuite) TestExportInvalidFormatShouldReturnAnError() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)

	res_note, err := s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{Note: &notespb.Note{AuthorId: generatedUuid.String(), Title: "Placeholder Title", Blocks: []*notespb.Block{}}})

	s.Nil(err)
	s.NotNil(res_note)

	res, err := s.srv.ExportNote(ctx, &notespb.ExportNoteRequest{NoteId: res_note.Note.Id, ExportFormat: notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_INVALID})

	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

func (s *ExportAPISuite) TestExportMarkdownShouldBeValid() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)

	res_note, err := s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: generatedUuid.String(),
			Title:    "Placeholder Title",
			Blocks: []*notespb.Block{
				{
					Type: notespb.Block_TYPE_HEADING_1,
					Data: &notespb.Block_Heading{Heading: "Heading"},
				},
				{
					Type: notespb.Block_TYPE_HEADING_2,
					Data: &notespb.Block_Heading{Heading: "Heading 2"},
				},
			},
		}})
	s.Nil(err)
	s.NotNil(res_note)

	res, err := s.srv.ExportNote(ctx, &notespb.ExportNoteRequest{NoteId: res_note.Note.Id, ExportFormat: notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_MARKDOWN})
	s.Nil(err)
	s.Equal(string(res.File), "# Heading\n## Heading 2\n")
}

func (s *ExportAPISuite) TestExportPdfShouldBeValid() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)

	res_note, err := s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: generatedUuid.String(),
			Title:    "Placeholder Title",
			Blocks: []*notespb.Block{
				{
					Type: notespb.Block_TYPE_HEADING_1,
					Data: &notespb.Block_Heading{Heading: "Heading"},
				},
				{
					Type: notespb.Block_TYPE_HEADING_2,
					Data: &notespb.Block_Heading{Heading: "Heading 2"},
				},
			},
		}})
	s.Nil(err)
	s.NotNil(res_note)

	res, err := s.srv.ExportNote(ctx, &notespb.ExportNoteRequest{NoteId: res_note.Note.Id, ExportFormat: notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_PDF})
	s.Nil(err)
	s.NotNil(res)
	s.Greater(len(res.File), len("%PDF"))
}
