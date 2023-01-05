package main

import (
	"context"
	notespb "notes-service/protorepo/noted/notes/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *notesService) GenerateWidgets(ctx context.Context, in *notespb.GenerateWidgetsRequest) (*notespb.GenerateWidgetsResponse, error) {
	_, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	/*err = validators.ValidateGetNoteRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}*/

	//Check if user is a note group member

	note, err := srv.repoNote.Get(ctx, in.NoteId)
	if err != nil {
		srv.logger.Error("failed to get note", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not get note.")
	}

	widgets := make([]*notespb.Widget, len(note.Keywords))

	for index, keyWord := range note.Keywords {
		//remplir les widget en fonction du type
	}

	return &notespb.GenerateWidgetsResponse{Widgets: widgets}, nil
}
