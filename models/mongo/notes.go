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
	update, err := srv.db.Collection("notes").UpdateOne(ctx, buildNoteQuery(filter), bson.D{{Key: "$set", Value: &noteRequest}})
	if err != nil {
		srv.logger.Error("failed to convert object id from hex", zap.Error(err))
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	if update.MatchedCount == 0 {
		srv.logger.Error("mongo update note query matched none", zap.String("user_id", filter.ID.String()))
		return status.Errorf(codes.Internal, "could not update note")
	}
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
