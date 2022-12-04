package main

import (
	"context"
	"strings"

	"notes-service/exports"
	notespb "notes-service/protorepo/noted/notes/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *notesService) ExportNote(ctx context.Context, in *notespb.ExportNoteRequest) (*notespb.ExportNoteResponse, error) {
	m := map[string]func(*notespb.Note) ([]byte, error){
		"markdown": exports.NoteToMarkdown,
		"pdf":      exports.NoteToPDF,
	}

	note, err := srv.GetNote(ctx, &notespb.GetNoteRequest{Id: in.NoteId})

	if err != nil {
		srv.logger.Error("failed to fetch note", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to retrieve note with id : %s", in.NoteId)
	}

	splittedEnumName := strings.Split(in.ExportFormat.String(), "_")
	formatName := strings.ToLower(splittedEnumName[len(splittedEnumName)-1])

	formatter, ok := m[formatName]

	if !ok {
		srv.logger.Error("format not recognized", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "format not recognized : %s", formatName)
	}

	fileBytes, err := formatter(note.Note)

	if err != nil {
		srv.logger.Error("failed to convert note", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to convert note to: %s", formatName)
	}

	return &notespb.ExportNoteResponse{File: fileBytes}, nil
}
