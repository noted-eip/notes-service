// Package memory is an in-memory implementation of models.AccountsRepository
package memory

import (
	"context"
	"notes-service/models"

	"github.com/google/uuid"
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

func (srv *notesRepository) Create(ctx context.Context, noteRequest *models.Note) (*models.Note, error) {
	txn := srv.db.DB.Txn(true)
	defer txn.Abort()

	id, err := uuid.NewRandom()

	if err != nil {
		srv.logger.Error("failed to generate new random uuid", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not create account")
	}
	noteRequest.ID = id
	note := models.Note{ID: noteRequest.ID, AuthorId: noteRequest.AuthorId, Title: noteRequest.Title, Blocks: noteRequest.Blocks}

	err = txn.Insert("note", note)
	if err != nil {
		srv.logger.Error("mongo insert note failed", zap.Error(err), zap.String("note name", note.AuthorId))
		return nil, status.Errorf(codes.Internal, "could not create note")
	}
	return noteRequest, nil
}

func (srv *notesRepository) Get(ctx context.Context, noteId *string) (*models.Note, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) Delete(ctx context.Context, noteId *string) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) Update(ctx context.Context, noteId *string, noteRequest *models.Note) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) List(ctx context.Context, authorId *string) (*[]models.Note, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
