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

func (srv *notesAPI) CreateNote(ctx context.Context, req *notesv1.CreateNoteRequest) (*notesv1.CreateNoteResponse, error) {

	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesAPI) GetNote(ctx context.Context, req *notesv1.GetNoteRequest) (*notesv1.GetNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesAPI) UpdateNote(ctx context.Context, req *notesv1.UpdateNoteRequest) (*notesv1.UpdateNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesAPI) DeleteNote(ctx context.Context, req *notesv1.DeleteNoteRequest) (*notesv1.DeleteNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesAPI) ListNotes(ctx context.Context, req *notesv1.ListNotesRequest) (*notesv1.ListNotesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

var protobufFormatToFormatter = map[notesv1.NoteExportFormat]func(*notesv1.Note) ([]byte, error){
	notesv1.NoteExportFormat_NOTE_EXPORT_FORMAT_MARKDOWN: exports.NoteToMarkdown,
	notesv1.NoteExportFormat_NOTE_EXPORT_FORMAT_PDF:      exports.NoteToPDF,
}

func (srv *notesAPI) ExportNote(ctx context.Context, req *notesv1.ExportNoteRequest) (*notesv1.ExportNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (srv *notesAPI) authenticate(ctx context.Context) (*auth.Token, error) {
	token, err := srv.auth.TokenFromContext(ctx)
	if err != nil {
		srv.logger.Debug("could not authenticate request", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return token, nil
}
