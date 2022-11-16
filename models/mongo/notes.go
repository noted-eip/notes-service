package mongo

import (
	"context"
	"notes-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notesRepository struct {
	logger *zap.Logger
	db     *mongo.Database
}

func NewNotesRepository(db *mongo.Database, logger *zap.Logger) models.NotesRepository {
	return &notesRepository{
		logger: logger,
		db:     db,
	}
}

func (srv *notesRepository) Create(ctx context.Context, noteRequest *models.Note) (*models.Note, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) Get(ctx context.Context, noteId *string) (*models.Note, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) Delete(ctx context.Context, noteId *string) error {
	delete, err := srv.db.Collection("notes").DeleteOne(ctx, buildNoteQuery(noteId))

	if err != nil {
		srv.logger.Error("delete note db query failed", zap.Error(err))
		return status.Errorf(codes.Internal, "could not delete note")
	}
	if delete.DeletedCount == 0 {
		srv.logger.Info("mongo delete note matched none", zap.String("note_id", *noteId))
		return status.Errorf(codes.Internal, "could not delete note")
	}
	return nil
}

func (srv *notesRepository) Update(ctx context.Context, noteId *string, noteRequest *models.Note) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) List(ctx context.Context, authorId *string) (*[]models.Note, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func buildNoteQuery(noteId *string) bson.M {
	query := bson.M{}
	if *noteId != "" {
		query["_id"] = noteId
	}
	return query
}
