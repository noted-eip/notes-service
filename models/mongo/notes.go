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
	return nil, nil
}

func (srv *notesRepository) Delete(ctx context.Context, filter *models.NoteFilter) error {
	return nil
}

func (srv *notesRepository) Update(ctx context.Context, filter *models.NoteFilter, noteRequest *models.NoteWithBlocks) error {
	return nil
}

func (srv *notesRepository) List(ctx context.Context, filter *models.NoteFilter) (*[]models.NoteWithBlocks, error) {
	notesCursor, err := srv.db.Collection("notes").Find(ctx, buildNoteQuery(filter))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "note not found")
		}
		srv.logger.Error("unable to query note", zap.Error(err))
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	notesResponse := make([]models.NoteWithBlocks, notesCursor.RemainingBatchLength())

	//convert notes from mongo to []models.NoteWithBlocks
	var notes []bson.M
	if err := notesCursor.All(context.TODO(), &notes); err != nil {
		srv.logger.Error("unable to parse notes", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	for index, note := range notes {
		id, err := uuid.Parse(note["_id"].(string))
		if err != nil {
			srv.logger.Error("unable to retrieve id of the note", zap.Error(err))
			return nil, status.Errorf(codes.Aborted, err.Error())
		}
		notesResponse[index] = models.NoteWithBlocks{ID: id, AuthorId: note["authorId"].(string), Title: note["title"].(string)}
	}

	return &notesResponse, nil
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
