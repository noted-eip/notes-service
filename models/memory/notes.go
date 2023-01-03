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

func NewNotesDatabaseSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"note": {
				Name: "note",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"author_id": {
						Name:    "author_id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "AuthorId"},
					},
					"title": {
						Name:    "title",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Title"},
					},
					"blocks": {
						Name:    "blocks",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Blocks"},
					},
				},
			},
		},
	}
}

func (srv *notesRepository) Create(ctx context.Context, noteRequest *models.NotePayload) (*models.Note, error) {
	txn := srv.db.DB.Txn(true)
	id, err := uuid.NewRandom()
	defer txn.Abort()

	if noteRequest == nil {
		srv.logger.Error("NoteRequest is nil")
		return nil, status.Errorf(codes.Internal, "could not create account")
	}

	note := models.Note{ID: id.String(), AuthorId: noteRequest.AuthorId, Title: noteRequest.Title}

	err = txn.Insert("note", note)
	if err != nil {
		srv.logger.Error("mongo insert note failed", zap.Error(err), zap.String("note name", note.AuthorId))
		return nil, status.Errorf(codes.Internal, "could not create note")
	}
	return &note, nil
}

func (srv *notesRepository) Get(ctx context.Context, noteId string) (*models.Note, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) Delete(ctx context.Context, noteId string) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
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

func (srv *notesRepository) List(ctx context.Context, authorId string) (*[]models.Note, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
