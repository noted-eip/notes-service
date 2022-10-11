package main

import (
	"context"

	"notes-service/auth"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"
	recommendationspb "notes-service/protorepo/noted/recommendations/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notesService struct {
	notespb.UnimplementedNotesAPIServer

	auth      auth.Service
	logger    *zap.SugaredLogger
	repoNote  models.NotesRepository
	repoBlock models.BlocksRepository

	recommendationClient recommendationspb.RecommendationsAPIClient
}

var _ notespb.NotesAPIServer = &notesService{}

func (srv *notesService) CreateNote(ctx context.Context, in *notespb.CreateNoteRequest) (*notespb.CreateNoteResponse, error) {
	/*
		recommendationRequest := &recommendationspb.ExtractKeywordsRequest{Content: "Alan Mathison Turing, né le 23 juin 1912 à Londres et mort le 7 juin 1954 à Wilmslow, est un mathématicien et cryptologue britannique, auteur de travaux qui fondent scientifiquement l'informatique.Pour résoudre le problème fondamental de la décidabilité en arithmétique, il présente en 1936 une expérience de pensée que l'on nommera ensuite machine de Turing et des concepts de programme et de programmation, qui prendront tout leur sens avec la diffusion des ordinateurs, dans la seconde moitié du XXe siècle. Son modèle a contribué à établir la thèse de Church, qui définit le concept mathématique intuitif de fonction calculable.Durant la Seconde Guerre mondiale, il joue un rôle majeur dans la cryptanalyse de la machine Enigma utilisée par les armées allemandes : l'invention de machines usant de procédés électroniques, les bombes1, fera passer le décryptage à plusieurs milliers de messages par jour. Mais tout ce travail doit forcément rester secret, et ne sera connu du public que dans les années 1970. Après la guerre, il travaille sur un des tout premiers ordinateurs, puis contribue au débat sur la possibilité de l'intelligence artificielle, en proposant le test de Turing."}
		if srv.recommendationClient == nil {
			fmt.Print("recommendation client is nil")
		}
		test, err := srv.recommendationClient.ExtractKeywords(ctx, recommendationRequest)
		if err != nil {
			fmt.Print("zebi : %v", err)
		}
		fmt.Print(test.Keywords)
	*/
	if len(in.Note.AuthorId) < 1 || len(in.Note.Title) < 1 {
		srv.logger.Errorw("failed to create note, invalid parameters")
		return nil, status.Errorf(codes.Internal, "authorId or title are empty")
	}

	note, err := srv.repoNote.Create(ctx, &models.NoteWithBlocks{AuthorId: in.Note.AuthorId, Title: in.Note.Title, Blocks: nil})

	if err != nil {
		srv.logger.Errorw("failed to create note", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not create note")
	}

	blocks := make([]models.Block, len(in.Note.Blocks))

	for index, block := range in.Note.Blocks {
		err := FillBlockContent(&blocks[index], block)
		if err != nil {
			srv.logger.Errorw("failed to create note", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "invalid content provided for block index : ", index)
		}
		//Get recommendation tags
		blockContent, err := GetDataContent(in.Note.Blocks[index])
		if err != nil {
			srv.logger.Errorw("failed to convert the content of the block", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "failed to convert the content of the block : ", index)
		}
		recommendationRequest := &recommendationspb.ExtractKeywordsRequest{Content: blockContent}
		clientResponse, err := srv.recommendationClient.ExtractKeywords(ctx, recommendationRequest)
		if err != nil {
			srv.logger.Errorw("failed to get the recommendation from client", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "failed to get the recommendation from client")
		}

		srv.repoBlock.Create(ctx, &models.BlockWithTags{NoteId: note.ID.String(), Type: uint32(in.Note.Blocks[index].Type), Index: uint32(index + 1), Content: blocks[index].Content, Tags: clientResponse.Keywords})
	}

	return &notespb.CreateNoteResponse{
		Note: &notespb.Note{
			Id:       note.ID.String(),
			AuthorId: note.AuthorId,
			Title:    note.Title,
			Blocks:   in.Note.Blocks,
		},
	}, nil
}

func (srv *notesService) GetNote(ctx context.Context, in *notespb.GetNoteRequest) (*notespb.GetNoteResponse, error) {
	return nil, nil
}

func (srv *notesService) UpdateNote(ctx context.Context, in *notespb.UpdateNoteRequest) (*notespb.UpdateNoteResponse, error) {
	return nil, nil
}

func (srv *notesService) DeleteNote(ctx context.Context, in *notespb.DeleteNoteRequest) (*notespb.DeleteNoteResponse, error) {
	return nil, nil
}

func (srv *notesService) ListNotes(ctx context.Context, in *notespb.ListNotesRequest) (*notespb.ListNotesResponse, error) {
	return nil, nil
}
