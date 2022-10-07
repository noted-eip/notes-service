package main

import (
	"context"

	"notes-service/auth"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"

	"github.com/google/uuid"
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
}

var _ notespb.NotesAPIServer = &notesService{}

func (srv *notesService) CreateNote(ctx context.Context, in *notespb.CreateNoteRequest) (*notespb.CreateNoteResponse, error) {
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
		srv.repoBlock.Create(ctx, &models.BlockWithIndex{NoteId: note.ID.String(), Type: uint32(in.Note.Blocks[index].Type), Index: uint32(index + 1), Content: blocks[index].Content})
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
	id, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get note")
	}

	note, err := srv.repoNote.Get(ctx, &models.NoteFilter{ID: id, AuthorId: ""})
	if err != nil {
		srv.logger.Errorw("failed to get note", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get note")
	}

	blocksTmp, err := srv.repoBlock.GetAllById(ctx, &models.BlockFilter{NoteId: note.ID.String()})
	if err != nil {
		srv.logger.Errorw("failed to get blocks", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid content provided for blocks form noteId : ", note.ID)
	}

	//convert []models.block to []notespb.Block
	blocks := make([]*notespb.Block, len(blocksTmp))
	for index, block := range blocksTmp {
		blocks[index] = &notespb.Block{}
		err := FillContentFromModelToApi(block, block.Type, blocks[index])
		if err != nil {
			srv.logger.Errorw("failed to the content of a block", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "fail to get content from block Id : ", block.ID)
		}
		blocks[index] = &notespb.Block{Id: block.ID, Type: notespb.Block_Type(block.Type), Data: blocks[index].Data}
	}
	noteToReturn := notespb.Note{Id: note.ID.String(), AuthorId: note.AuthorId, Title: note.Title, Blocks: blocks}

	return &notespb.GetNoteResponse{Note: &noteToReturn}, nil
}

func (srv *notesService) UpdateNote(ctx context.Context, in *notespb.UpdateNoteRequest) (*notespb.UpdateNoteResponse, error) {
	id, err := uuid.Parse(in.Note.Id)
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not update note")
	}

	//appeler deleteBlock avec le filtre note_id
	err = srv.repoBlock.Delete(ctx, &models.BlockFilter{NoteId: in.Id})
	if err != nil {
		srv.logger.Errorw("blocks weren't deleted : ", err.Error())
		return nil, status.Errorf(codes.Internal, "could not delete blocks")
	}

	//update juste les infos de la note et pas les blocks sinon
	err = srv.repoNote.Update(ctx, &models.NoteFilter{ID: id}, &models.NoteWithBlocks{AuthorId: in.Note.AuthorId, Title: in.Note.Title})
	if err != nil {
		srv.logger.Errorw("failed to update note", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not update note")
	}

	//appeller createBlock en boucle pour tout les autres blocks
	blocks := make([]models.Block, len(in.Note.Blocks))
	for index, block := range in.Note.Blocks {
		err := FillBlockContent(&blocks[index], block)
		if err != nil {
			srv.logger.Errorw("failed to update blocks", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "invalid content provided for block index : ", index)
		}
		srv.repoBlock.Create(ctx, &models.BlockWithIndex{NoteId: in.Id, Type: uint32(in.Note.Blocks[index].Type), Index: uint32(index + 1), Content: blocks[index].Content})
	}

	return nil, nil
}

func (srv *notesService) DeleteNote(ctx context.Context, in *notespb.DeleteNoteRequest) (*notespb.DeleteNoteResponse, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not delete note")
	}

	//appeler deleteBlock avec le filtre note_id
	err = srv.repoBlock.Delete(ctx, &models.BlockFilter{NoteId: in.Id})
	if err != nil {
		srv.logger.Errorw("blocks weren't deleted : ", err.Error())
		return nil, status.Errorf(codes.Internal, "could not delete blocks")
	}

	err = srv.repoNote.Delete(ctx, &models.NoteFilter{ID: id})
	if err != nil {
		srv.logger.Errorw("failed to delete note", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not delete note")
	}

	return nil, nil
}

func (srv *notesService) ListNotes(ctx context.Context, in *notespb.ListNotesRequest) (*notespb.ListNotesResponse, error) {
	if len(in.AuthorId) < 1 {
		srv.logger.Errorw("failed to lists notes, invalid parameters")
		return nil, status.Errorf(codes.Internal, "authorId is empty")
	}

	notes, err := srv.repoNote.List(ctx, &models.NoteFilter{AuthorId: in.AuthorId})
	if err != nil {
		srv.logger.Errorw("failed to get note", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get note")
	}

	notesResponse := make([]*notespb.Note, len(*notes))
	for index, note := range *notes {
		notesResponse[index] = &notespb.Note{Id: note.ID.String(), AuthorId: note.AuthorId, Title: note.Title}
	}

	return &notespb.ListNotesResponse{
		Notes: notesResponse,
	}, nil
}
