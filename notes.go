package main

import (
	"context"
	"fmt"

	automapper "github.com/stroiman/go-automapper"

	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notesService struct {
	notespb.UnimplementedNotesAPIServer

	logger *zap.SugaredLogger
	repo   models.NotesRepository
}

var _ notespb.NotesAPIServer = &notesService{}

func (srv *notesService) CreateNote(ctx context.Context, in *notespb.CreateNoteRequest) (*notespb.CreateNoteResponse, error) {
	//fmt.Print("on passe 1\n")
	//blocks := []*models.BlockTest{}

	//fmt.Print("on passe 3\n")
	//automapper.Map(in.Note.Blocks, &blocks)

	fmt.Print("on passe 1\n")
	fmt.Print("authorid grpc : ", &in.Note.AuthorId, "\n")
	fmt.Print("len blocks grpc : ", len(in.Note.Blocks), "\n")

	blocks := make([]models.BlockTest, len(in.Note.Blocks))

	fmt.Print("on passe 2\n")
	for i := range in.Note.Blocks {
		automapper.Map(in.Note.Blocks[i], &blocks[i])
	}

	fmt.Print("on passe 3\n")
	err := srv.repo.Create(ctx, &models.NoteWithBlocks{AuthorId: in.Note.AuthorId, Title: &in.Note.Title, Blocks: nil})

	if err != nil {
		srv.logger.Errorw("failed to create note", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not create note")
	}

	return nil, nil
}

func (srv *notesService) GetNote(ctx context.Context, in *notespb.GetNoteRequest) (*notespb.GetNoteResponse, error) {

	//Appeler GetBlock avec le filtre noteId
	//les classer selon leur index

	id, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get note")
	}

	note, err := srv.repo.Get(ctx, &models.NoteFilter{ID: id, AuthorId: ""})
	if err != nil {
		srv.logger.Errorw("failed to get account", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get note")
	}
	noteToReturn := notespb.Note{Id: note.ID.String(), AuthorId: note.AuthorId, Title: *note.Title, Blocks: nil /*note.Blocks*/}
	return &notespb.GetNoteResponse{Note: &noteToReturn}, nil
}

func (srv *notesService) UpdateNote(ctx context.Context, in *notespb.UpdateNoteRequest) (*notespb.UpdateNoteResponse, error) {

	//appeler deleteBlock avec le filtre note_id
	//appeller createBlock pour tout les autres

	id, err := uuid.Parse(in.Note.Id)
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not update note")
	}

	err = srv.repo.Update(ctx, &models.NoteFilter{ID: id, AuthorId: ""}, &models.NoteWithBlocks{AuthorId: in.Note.AuthorId, Title: &in.Note.Title, Blocks: nil /*in.Note.Blocks*/})
	if err != nil {
		srv.logger.Errorw("failed to create note", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not create note")
	}

	return nil, nil
}

func (srv *notesService) DeleteNote(ctx context.Context, in *notespb.DeleteNoteRequest) (*notespb.DeleteNoteResponse, error) {

	//deleteBlock avec le filtre noteId
	//delete la note

	id, err := uuid.Parse(in.Id)
	if err != nil {
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not delete note")
	}

	err = srv.repo.Delete(ctx, &models.NoteFilter{ID: id})
	if err != nil {
		srv.logger.Errorw("failed to delete note", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not delete note")
	}

	return nil, nil
}

func (srv *notesService) ListNotes(ctx context.Context, in *notespb.ListNotesRequest) (*notespb.ListNotesResponse, error) {

	//appeler GetNote avec le filtre authorName

	return nil, status.Errorf(codes.OK, "Note found")
}
