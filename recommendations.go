package main

import (
	"context"
	"notes-service/auth"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type recommendationsService struct {
	notespb.UnimplementedRecommendationsAPIServer

	auth     auth.Service
	logger   *zap.Logger
	repoNote models.NotesRepository
}

var _ notespb.RecommendationsAPIServer = &recommendationsService{}

func (srv *recommendationsService) GenerateWidgets(ctx context.Context, in *notespb.GenerateWidgetsRequest) (*notespb.GenerateWidgetsResponse, error) {
	_, err := Authenticate1(srv, ctx)
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

	var widgets []*notespb.Widget

	for _, keyWord := range note.Keywords {

		if keyWord.URL != "" {
			widgets = append(widgets, &notespb.Widget{
				Type: &notespb.Widget_WebsiteWidget{
					WebsiteWidget: &notespb.WebsiteWidget{
						Title:       keyWord.Keyword,
						Url:         keyWord.URL,
						Description: string(keyWord.Type),
					},
				},
			})
		}
	}

	return &notespb.GenerateWidgetsResponse{Widgets: widgets}, nil
}

func Authenticate1(srv *recommendationsService, ctx context.Context) (*auth.Token, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return token, nil
}

func (srv *recommendationsService) authenticate(ctx context.Context) (*auth.Token, error) {
	token, err := srv.auth.TokenFromContext(ctx)
	if err != nil {
		srv.logger.Debug("failed to authenticate request", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return token, nil
}
