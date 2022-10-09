package mongo

import (
	"context"
	"notes-service/models"

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
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) Get(ctx context.Context, filter *models.NoteFilter) (*models.NoteWithBlocks, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) Delete(ctx context.Context, filter *models.NoteFilter) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) Update(ctx context.Context, filter *models.NoteFilter, noteRequest *models.NoteWithBlocks) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesRepository) List(ctx context.Context, filter *models.NoteFilter) (*[]models.NoteWithBlocks, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
