package mongo

import (
	"context"
	"notes-service/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type note struct {
	ID       string         `json:"id" bson:"_id,omitempty"`
	AuthorId string         `json:"authorId" bson:"authorId,omitempty"`
	Title    *string        `json:"title" bson:"title,omitempty"`
	Blocks   []models.Block `json:"blocks" bson:"blocks,omitempty"`
}

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

func (srv *notesRepository) Create(ctx context.Context, noteRequest *models.NoteWithBlocks) (*models.NoteWithBlocks, error) {
	return nil, nil
}

func (srv *notesRepository) Get(ctx context.Context, filter *models.NoteFilter) (*models.NoteWithBlocks, error) {
	return nil, nil
}

func (srv *notesRepository) Delete(ctx context.Context, filter *models.NoteFilter) error {
	delete, err := srv.db.Collection("notes").DeleteOne(ctx, buildNoteQuery(filter))

	if err != nil {
		srv.logger.Error("delete note db query failed", zap.Error(err))
		return status.Errorf(codes.Internal, "could not delete note")
	}
	if delete.DeletedCount == 0 {
		srv.logger.Info("mongo delete note matched none", zap.String("note_id", filter.ID.String()))
		return status.Errorf(codes.Internal, "could not delete note")
	}
	return nil
}

func (srv *notesRepository) Update(ctx context.Context, filter *models.NoteFilter, noteRequest *models.NoteWithBlocks) error {
	return nil
}

func (srv *notesRepository) List(ctx context.Context, filter *models.NoteFilter) (*[]models.NoteWithBlocks, error) {
	return nil, nil
}

func buildNoteQuery(filter *models.NoteFilter) bson.M {
	query := bson.M{}
	if filter.ID != uuid.Nil {
		query["_id"] = filter.ID.String()
	}
	if filter.AuthorId != "" {
		query["authorId"] = filter.AuthorId
	}
	return query
}
