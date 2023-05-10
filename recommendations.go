package main

import (
	"context"
	"notes-service/auth"
	"notes-service/language"
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"

	"notes-service/validators"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type recommendationsAPI struct {
	notesv1.UnimplementedRecommendationsAPIServer

	logger *zap.Logger

	auth     auth.Service
	language language.Service
	notes    models.NotesRepository
}

var _ notesv1.RecommendationsAPIServer = &recommendationsAPI{}

func (srv *recommendationsAPI) GenerateWidgets(ctx context.Context, req *notesv1.GenerateWidgetsRequest) (*notesv1.GenerateWidgetsResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = validators.ValidateGenerateWidgetsRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	note, err := srv.notes.GetNote(ctx, &models.OneNoteFilter{NoteID: req.NoteId, GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		srv.logger.Error("failed to get note", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not get note.")
	}

	var widgets []*notesv1.Widget

	for _, keyword := range note.Keywords {

		widgets = append(widgets, &notesv1.Widget{
			Type: &notesv1.Widget_WebsiteWidget{
				WebsiteWidget: &notesv1.WebsiteWidget{
					Keyword:  keyword.Keyword,
					Type:     keyword.Type,
					Url:      keyword.URL,
					Summary:  keyword.Summary,
					ImageUrl: keyword.ImageURL,
				},
			},
		})
	}

	return &notesv1.GenerateWidgetsResponse{Widgets: widgets}, nil
}

func (srv *recommendationsAPI) authenticate(ctx context.Context) (*auth.Token, error) {
	token, err := srv.auth.TokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return token, nil
}
