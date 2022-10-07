package mongo

import (
	"context"
	"errors"
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
	var note note

	err := srv.db.Collection("notes").FindOne(ctx, buildNoteQuery(filter)).Decode(&note)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "note not found")
		}
		srv.logger.Error("unable to query note", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	uuid, err := uuid.Parse(note.ID)
	if err != nil {
		srv.logger.Error("failed to convert uuid from string", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not get note")
	}

	return &models.NoteWithBlocks{ID: uuid, AuthorId: note.AuthorId, Title: *note.Title, Blocks: note.Blocks}, nil
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
