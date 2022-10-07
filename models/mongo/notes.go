package mongo

import (
	"context"
	"notes-service/models"

	"github.com/google/uuid"
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
	id, err := uuid.NewRandom()

	if err != nil {
		srv.logger.Error("failed to generate new random uuid", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not create account")
	}
	noteRequest.ID = id

	note := note{ID: noteRequest.ID.String(), AuthorId: noteRequest.AuthorId, Title: &noteRequest.Title, Blocks: noteRequest.Blocks}

	_, err = srv.db.Collection("notes").InsertOne(ctx, note)
	if err != nil {
		srv.logger.Error("mongo insert note failed", zap.Error(err), zap.String("note name", note.AuthorId))
		return nil, status.Errorf(codes.Internal, "could not create note")
	}
	return noteRequest, nil
}

func (srv *notesRepository) Get(ctx context.Context, filter *models.NoteFilter) (*models.NoteWithBlocks, error) {
	return nil, nil
}

func (srv *notesRepository) Delete(ctx context.Context, filter *models.NoteFilter) error {
	return nil
}

func (srv *notesRepository) Update(ctx context.Context, filter *models.NoteFilter, noteRequest *models.NoteWithBlocks) error {
	return nil
}

func (srv *notesRepository) List(ctx context.Context, filter *models.NoteFilter) (*[]models.NoteWithBlocks, error) {
	return nil, nil
}
