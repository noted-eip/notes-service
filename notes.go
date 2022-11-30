package main

import (
	"context"

	"notes-service/auth"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"
	recommendationspb "notes-service/protorepo/noted/recommendations/v1"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notesService struct {
	notespb.UnimplementedNotesAPIServer

	auth                 auth.Service
	logger               *zap.Logger
	repoNote             models.NotesRepository
	repoBlock            models.BlocksRepository
	recommendationClient recommendationspb.RecommendationsAPIClient
}

var _ notespb.NotesAPIServer = &notesService{}

func (srv *notesService) CreateNote(ctx context.Context, in *notespb.CreateNoteRequest) (*notespb.CreateNoteResponse, error) {
	if in == nil {
		srv.logger.Error("failed to create note, Request is empty")
		return nil, status.Error(codes.InvalidArgument, "CreateNoteRequest is empty")
	}
	if in.Note == nil {
		srv.logger.Error("failed to create note, Note Request is empty")
		return nil, status.Error(codes.InvalidArgument, "Note is empty")
	}
	if len(in.Note.AuthorId) < 1 || len(in.Note.Title) < 1 {
		srv.logger.Error("failed to create note, invalid parameters")
		return nil, status.Error(codes.InvalidArgument, "authorId or title are empty")
	}

	//stopper la goroutine sur blockId

	note, err := srv.repoNote.Create(ctx, &models.Note{AuthorId: in.Note.AuthorId, Title: in.Note.Title, Blocks: nil})

	if err != nil {
		srv.logger.Error("failed to create note", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not create note")
	}
	//create block with tags in function
	err = CreateBlockWithTags(srv, ctx, &in.Note.Id, in.Note.Blocks)
	if err != nil {
		srv.logger.Error("failed to create blocks", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//lancer la goroutine sur blockId

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
	_, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Error("failed to convert uuid from string", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "could not get note")
	}

	note, err := srv.repoNote.Get(ctx, &in.Id)
	if err != nil {
		srv.logger.Error("failed to get note", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "could not get note")
	}

	noteId := note.ID.String()
	blocksTmp, err := srv.repoBlock.GetBlocks(ctx, &noteId)
	if err != nil {
		srv.logger.Error("failed to get blocks", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "invalid content provided for blocks form noteId : %d", note.ID)
	}

	//convert []models.block to []notespb.Block
	blocks := make([]*notespb.Block, len(blocksTmp))
	for index, block := range blocksTmp {
		blocks[index] = &notespb.Block{}
		err := FillContentFromModelToApi(block, block.Type, blocks[index])
		if err != nil {
			srv.logger.Error("failed to the content of a block", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "fail to get content from block Id : %s", block.ID)
		}
		blocks[index] = &notespb.Block{Id: block.ID, Type: notespb.Block_Type(block.Type), Data: blocks[index].Data}
	}
	noteToReturn := notespb.Note{Id: note.ID.String(), AuthorId: note.AuthorId, Title: note.Title, Blocks: blocks}

	return &notespb.GetNoteResponse{Note: &noteToReturn}, nil
}

func (srv *notesService) UpdateNote(ctx context.Context, in *notespb.UpdateNoteRequest) (*notespb.UpdateNoteResponse, error) {
	if in == nil {
		srv.logger.Error("failed to update note, Request is empty")
		return nil, status.Error(codes.InvalidArgument, "UpdateNoteRequest is empty")
	}
	if in.Note == nil {
		srv.logger.Error("failed to update note, Note Request is empty")
		return nil, status.Error(codes.InvalidArgument, "Note is empty")
	}

	id, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Error("failed to convert uuid from string", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not update note")
	}

	//appeler deleteBlock avec le filtre note_id
	err = srv.repoBlock.DeleteBlocks(ctx, &in.Id)
	if err != nil {
		srv.logger.Error("blocks weren't deleted : ", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not delete blocks")
	}

	//update juste les infos de la note
	noteId := id.String()
	err = srv.repoNote.Update(ctx, &noteId, &models.Note{AuthorId: in.Note.AuthorId, Title: in.Note.Title})
	if err != nil {
		srv.logger.Error("failed to update note", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not update note")
	}

	//appeller createBlock en boucle pour tout les autres blocks
	err = CreateBlockWithTags(srv, ctx, &in.Id, in.Note.Blocks)
	if err != nil {
		srv.logger.Error("failed to create blocks", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func (srv *notesService) DeleteNote(ctx context.Context, in *notespb.DeleteNoteRequest) (*notespb.DeleteNoteResponse, error) {
	if in == nil {
		srv.logger.Error("failed to delete note, Request is empty")
		return nil, status.Error(codes.InvalidArgument, "DeleteNoteRequest is empty")
	}
	_, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Error("failed to convert uuid from string", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "could not delete note")
	}

	//appeler deleteBlock avec le filtre note_id
	err = srv.repoBlock.DeleteBlocks(ctx, &in.Id)
	if err != nil {
		srv.logger.Error("blocks weren't deleted : ", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not delete blocks")
	}

	err = srv.repoNote.Delete(ctx, &in.Id)
	if err != nil {
		srv.logger.Error("failed to delete note", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not delete note")
	}

	return nil, nil
}

func (srv *notesService) ListNotes(ctx context.Context, in *notespb.ListNotesRequest) (*notespb.ListNotesResponse, error) {
	if in == nil {
		srv.logger.Error("failed to delete note, Request is empty")
		return nil, status.Error(codes.InvalidArgument, "DeleteNoteRequest is empty")
	}
	if len(in.AuthorId) < 1 {
		srv.logger.Error("failed to lists notes, invalid parameters")
		return nil, status.Errorf(codes.InvalidArgument, "authorId is empty")
	}

	notes, err := srv.repoNote.List(ctx, &in.AuthorId)
	if err != nil {
		srv.logger.Error("failed to get note", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "could not get note")
	}

	notesResponse := make([]*notespb.Note, len(*notes))
	for index, note := range *notes {
		notesResponse[index] = &notespb.Note{Id: note.ID.String(), AuthorId: note.AuthorId, Title: note.Title}
	}

	return &notespb.ListNotesResponse{
		Notes: notesResponse,
	}, nil
}

//Mettre ca en processBackGround
func CreateBlockWithTags(srv *notesService, ctx context.Context, noteId *string, blocksGrpc []*notespb.Block) error {
	blocks := make([]models.Block, len(blocksGrpc))

	for index, block := range blocksGrpc {
		err := FillBlockContent(&blocks[index], block)
		if err != nil {
			srv.logger.Error("failed get block content", zap.Error(err))
			return status.Errorf(codes.Internal, "invalid content provided for block index : ", index)
		}

		content, err := GetDataContent(blocksGrpc[index])
		if err != nil {
			srv.logger.Error("failed to convert the content of the block", zap.Error(err))
			return status.Errorf(codes.Internal, "failed to convert the content of the block : ", index)
		}
		//find tags from block -> enelever ca après juste créer les block et faire ca en BackGround
		keywords, err := FindTags(srv, ctx, &content)
		if err != nil {
			srv.logger.Error("failed to get the recommendation from client", zap.Error(err))
			return status.Errorf(codes.Internal, "failed to get the recommendation from client")
		}

		srv.repoBlock.Create(ctx, &models.BlockWithTags{NoteId: *noteId, Type: uint32(blocksGrpc[index].Type), Index: uint32(index + 1), Content: blocks[index].Content, Tags: keywords})
	}
	return nil
}

func FindTags(srv *notesService, ctx context.Context, blockContent *string) ([]string, error) {
	recommendationRequest := &recommendationspb.ExtractKeywordsRequest{Content: *blockContent}
	clientResponse, err := srv.recommendationClient.ExtractKeywords(ctx, recommendationRequest)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get the recommendation from client")
	}
	return clientResponse.Keywords, nil
}
