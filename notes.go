package main

import (
	"context"

	"notes-service/auth"
	"notes-service/exports"
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"

	"notes-service/language"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notesAPI struct {
	notesv1.UnimplementedNotesAPIServer

	logger *zap.Logger

	auth     auth.Service
	language language.Service

	notes  models.NotesRepository
	groups models.GroupsRepository
}

var _ notesv1.NotesAPIServer = &notesAPI{}

func (srv *notesAPI) CreateNote(ctx context.Context, in *notesv1.CreateNoteRequest) (*notesv1.CreateNoteResponse, error) {

	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesAPI) GetNote(ctx context.Context, in *notesv1.GetNoteRequest) (*notesv1.GetNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesAPI) UpdateNote(ctx context.Context, in *notesv1.UpdateNoteRequest) (*notesv1.UpdateNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesAPI) DeleteNote(ctx context.Context, in *notesv1.DeleteNoteRequest) (*notesv1.DeleteNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesAPI) ListNotes(ctx context.Context, in *notesv1.ListNotesRequest) (*notesv1.ListNotesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

var protobufFormatToFormatter = map[notesv1.NoteExportFormat]func(*notesv1.Note) ([]byte, error){
	notesv1.NoteExportFormat_NOTE_EXPORT_FORMAT_MARKDOWN: exports.NoteToMarkdown,
	notesv1.NoteExportFormat_NOTE_EXPORT_FORMAT_PDF:      exports.NoteToPDF,
}

func (srv *notesAPI) ExportNote(ctx context.Context, in *notesv1.ExportNoteRequest) (*notesv1.ExportNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}
