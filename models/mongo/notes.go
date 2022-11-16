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
	return status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) Update(ctx context.Context, noteId *string, noteRequest *models.Note) error {
	update, err := srv.db.Collection("notes").UpdateOne(ctx, buildNoteQuery(noteId), bson.D{{Key: "$set", Value: &noteRequest}})
	if err != nil {
		srv.logger.Error("failed to convert object id from hex", zap.Error(err))
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	if update.MatchedCount == 0 {
		srv.logger.Error("mongo update note query matched none", zap.String("user_id", *noteId))
		return status.Errorf(codes.Internal, "could not update note")
	}
	return nil
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
