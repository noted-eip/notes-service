package main

import (
	"context"

	"notes-service/exports"
	notespb "notes-service/protorepo/noted/notes/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *notesService) ExportNote(ctx context.Context, in *notespb.ExportNoteRequest) (*notespb.ExportNoteResponse, error) {
	m := map[notespb.NoteExportFormat]func(*notespb.Note) ([]byte, error){
		notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_MARKDOWN: exports.NoteToMarkdown,
		notespb.NoteExportFormat_NOTE_EXPORT_FORMAT_PDF:      exports.NoteToPDF,
	}

	note, err := srv.GetNote(ctx, &notespb.GetNoteRequest{Id: in.NoteId})

	if err != nil {
		return nil, err
	}

	formatter, ok := m[in.ExportFormat]

	if !ok {
		srv.logger.Error("format not recognized", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "format not recognized : %s", in.ExportFormat.String())
	}

	fileBytes, err := formatter(note.Note)

	if err != nil {
		srv.logger.Error("failed to convert note", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to convert note to: %s", in.ExportFormat.String())
	}

	return &notespb.ExportNoteResponse{File: fileBytes}, nil
}
