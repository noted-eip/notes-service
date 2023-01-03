package main

import (
	"context"
	"notes-service/auth"
	"notes-service/models/memory"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExportAPISuite struct {
	suite.Suite
	srv *notesService
}

func TestExport(t *testing.T) {
	suite.Run(t, &ExportAPISuite{})
}

func (s *ExportAPISuite) SetupSuite() {
	logger := newLoggerOrFail(s.T())
	dbNote := newDatabaseOrFail(s.T(), logger)
	dbBlock := newDatabaseOrFail(s.T(), logger)

	s.srv = &notesService{
		auth:      &auth.TestService{},
		logger:    logger,
		repoNote:  memory.NewNotesRepository(dbNote, logger),
		repoBlock: memory.NewBlocksRepository(dbBlock, logger),
	}
}

func (s *ExportAPISuite) TestExportWrongNoteIDShouldReturnAnError() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.srv.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)

	generatedUuid, err = uuid.NewRandom()
	s.Require().NoError(err)
	res, err := s.srv.ExportNote(ctx, &notespb.ExportNoteRequest{NoteId: generatedUuid.String()})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.Internal)
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
	res_note, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{Note: &notespb.Note{AuthorId: "placeholder-author-id", Title: "Placeholder Title", Blocks: []*notespb.Block{}}})

	s.Nil(err)
	s.NotNil(res_note)

	res, err := s.srv.ExportNote(context.TODO(), &notespb.ExportNoteRequest{NoteId: res_note.Note.Id, ExportFormat: notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_INVALID})

	s.Require().Error(err)
	s.Equal(status.Code(err), codes.Internal)
	s.Nil(res)
}

func (s *ExportAPISuite) TestExportMarkdownShouldBeValid() {
	res_note, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: "placeholder-author-id",
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

	print(res_note.Note.Id)

	res, err := s.srv.ExportNote(context.TODO(), &notespb.ExportNoteRequest{NoteId: res_note.Note.Id, ExportFormat: notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_MARKDOWN})

	print(err.Error())

	s.Nil(err)
	s.NotNil(res)
	print(string(res.File))
}
