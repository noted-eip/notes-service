package main

import (
	"context"

	"notes-service/auth"
	"notes-service/exports"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"

	"notes-service/language"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notesService struct {
	notespb.UnimplementedNotesAPIServer

	logger   *zap.Logger
	auth     auth.Service
	language language.Service
	repoNote models.NotesRepository
}

var _ notespb.NotesAPIServer = &notesService{}

func (srv *notesService) CreateNote(ctx context.Context, in *notespb.CreateNoteRequest) (*notespb.CreateNoteResponse, error) {

	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesService) GetNote(ctx context.Context, in *notespb.GetNoteRequest) (*notespb.GetNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesService) UpdateNote(ctx context.Context, in *notespb.UpdateNoteRequest) (*notespb.UpdateNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesService) DeleteNote(ctx context.Context, in *notespb.DeleteNoteRequest) (*notespb.DeleteNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesService) ListNotes(ctx context.Context, in *notespb.ListNotesRequest) (*notespb.ListNotesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

var protobufFormatToFormatter = map[notespb.NoteExportFormat]func(*notespb.Note) ([]byte, error){
	notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_MARKDOWN: exports.NoteToMarkdown,
	notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_PDF:      exports.NoteToPDF,
}

func (srv *notesService) ExportNote(ctx context.Context, in *notespb.ExportNoteRequest) (*notespb.ExportNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}
