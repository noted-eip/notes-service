// Package memory is an in-memory implementation of models.AccountsRepository
package memory

import (
	"context"
	"notes-service/models"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notesRepository struct {
	logger *zap.Logger
	db     *Database
}

func NewNotesRepository(db *Database, logger *zap.Logger) models.NotesRepository {
	return &notesRepository{
		logger: logger.Named("memory").Named("notes"),
		db:     db,
	}
}

func (srv *notesRepository) Create(ctx context.Context, noteRequest *models.NotePayload) (*models.Note, error) {
	txn := srv.db.DB.Txn(true)
	defer txn.Abort()
	id, err := uuid.NewRandom()

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if noteRequest == nil {
		srv.logger.Error("NoteRequest is nil")
		return nil, status.Errorf(codes.Internal, "could not create account")
	}

	note := models.Note{ID: id.String(), AuthorId: noteRequest.AuthorId, Title: noteRequest.Title}

	err = txn.Insert("note", &note)
	if err != nil {
		srv.logger.Error("mongo insert note failed", zap.Error(err), zap.String("note name", note.AuthorId))
		return nil, status.Errorf(codes.Internal, "could not create note")
	}

	txn.Commit()
	return &note, nil
}

func (srv *notesRepository) Get(ctx context.Context, noteId string) (*models.Note, error) {
	txn := srv.db.DB.Txn(false)

	raw, err := txn.First("note", "id", noteId)

	if err != nil {
		srv.logger.Error("unable to query note", zap.Error(err))
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound, "note not found")
	}
	txn.Commit()
	return raw.(*models.Note), nil
}

func (srv *notesRepository) Delete(ctx context.Context, noteId string) error {
	txn := srv.db.DB.Txn(true)
	defer txn.Abort()

	err := txn.Delete("note", models.Note{ID: noteId})

	if err == memdb.ErrNotFound {
		srv.logger.Error("unable to find note", zap.Error(err))
	}
	if err != nil {
		srv.logger.Error("unable to delete note", zap.Error(err))
		return err
	}
	txn.Commit()
	return nil
}

func (srv *notesRepository) Update(ctx context.Context, noteId string, noteRequest *models.NotePayload) error {
	/*txn := srv.db.DB.Txn(true)
	defer txn.Abort()

	//update, err := srv.db.Collection("notes").UpdateOne(ctx, buildNoteQuery(noteId), bson.D{{Key: "$set", Value: &noteRequest}})
	err := txn.Update("note", buildNoteQuery(noteId), bson.D{{Key: "$set", Value: &noteRequest}})

	if err != nil {
		srv.logger.Error("failed to convert object id from hex", zap.Error(err))
		return status.Error(codes.InvalidArgument, err.Error())
	}*/
	return nil
}

func (srv *notesRepository) List(ctx context.Context, authorId string) ([]*models.Note, error) {
	var notes []*models.Note

	txn := srv.db.DB.Txn(false)

	it, err := txn.Get("note", "author_id", authorId)
	if err != nil {
		srv.logger.Error("unable to list notes", zap.Error(err))
		return nil, err
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		notes = append(notes, obj.(*models.Note))
	}

	return notes, nil
}
