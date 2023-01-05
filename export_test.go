package main

import (
	"context"
	"notes-service/auth"
	"notes-service/models/memory"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExportAPISuite struct {
	suite.Suite
	auth auth.TestService
	srv  *notesService
}

func TestExport(t *testing.T) {
	suite.Run(t, &ExportAPISuite{})
}

func (s *ExportAPISuite) SetupSuite() {
	logger := newLoggerOrFail(s.T())
	dbNote := newDatabaseOrFail(s.T(), logger)
	dbBlock := newDatabaseOrFail(s.T(), logger)

	s.srv = &notesService{
		auth:      &s.auth,
		logger:    logger,
		repoNote:  memory.NewNotesRepository(dbNote, logger),
		repoBlock: memory.NewBlocksRepository(dbBlock, logger),
	}
}

func (s *ExportAPISuite) TestExportWrongNoteIDShouldReturnAnError() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)

	generatedUuid, err = uuid.NewRandom()
	s.Require().NoError(err)
	res, err := s.srv.ExportNote(ctx, &notespb.ExportNoteRequest{NoteId: generatedUuid.String()})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.NotFound)
	s.Nil(res)
}

func (s *ExportAPISuite) TestUnauthenticatedShouldReturnAnError() {
	uuid, _ := uuid.NewRandom()
	res, err := s.srv.ExportNote(context.TODO(), &notespb.ExportNoteRequest{NoteId: uuid.String()})

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
	s.Equal(status.Code(err), codes.Internal)
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

	a := time.Now()
	res, err := s.srv.ExportNote(ctx, &notespb.ExportNoteRequest{NoteId: res_note.Note.Id, ExportFormat: notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_PDF})
	b := time.Since(a)
	print(b.Seconds())
	s.Nil(err)
	s.NotNil(res)
	s.Greater(len(res.File), len("%PDF"))
}
