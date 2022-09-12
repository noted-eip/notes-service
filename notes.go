package main

import (
	"context"
	"fmt"
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

	//sol 1 parser et fill model.Block
	//coppier -> automapper
	err := srv.repo.Create(ctx, &models.NoteWithBlocks{AuthorId: in.Note.AuthorId, Title: &in.Note.Title, Blocks: nil /*in.Blocks*/})

	if err != nil {
		srv.logger.Errorw("failed to create note", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not create note")
	}

	return nil, nil
}

func (srv *notesService) GetNote(ctx context.Context, in *notespb.GetNoteRequest) (*notespb.GetNoteResponse, error) {
	fmt.Print("on est la\n")
	id, err := uuid.Parse(in.Id)
	if err != nil {
		//ca crash la prcq il faut que le CreatNote créer son propre UUID après on verra si on peut le parser
		fmt.Print("EROR HERE\n")
		fmt.Print(err)
		fmt.Print("\n")
		if srv.logger == nil {
			fmt.Print("logger == nil\n")
		}
		srv.logger.Errorw("failed to convert uuid from string", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get note")
	}
	fmt.Print("on est la 2\n")

	note, err := srv.repo.Get(ctx, &models.NoteFilter{ID: id, AuthorId: ""})
	if err != nil {
		srv.logger.Errorw("failed to get account", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "could not get note")
	}
	noteToReturn := notespb.Note{Id: note.ID.String(), AuthorId: note.AuthorId, Title: *note.Title, Blocks: nil /*note.Blocks*/}
	fmt.Print("on est la 3\n")
	return &notespb.GetNoteResponse{Note: &noteToReturn}, nil
}

func (srv *notesService) UpdateNote(ctx context.Context, in *notespb.UpdateNoteRequest) (*notespb.UpdateNoteResponse, error) {
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
	/*
		filter := bson.M{"authorId": in.AuthorId}
		cur, err := NotesCollection.Find(context.TODO(), filter)

		if err != nil {
			return nil, err
		}
		fmt.Println(cur)

		var results *notespb.Notes = nil
		for cur.Next(context.TODO()) {
			var elem *notespb.Note
			err := cur.Decode(&elem)
			if err != nil {
				return nil, err
			}
			results.Notes = append(results.Notes, elem)
		}
		err = cur.Err()

		if err != nil {
			return nil, err
		}

		fmt.Println(results)
	*/
	return nil, status.Errorf(codes.OK, "Note found")
}
