package mongo

import (
	"context"
	"errors"
	"notes-service/models"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

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

func (srv *notesRepository) Create(ctx context.Context, noteRequest *models.NotePayload) (*models.Note, error) {
	id, err := uuid.NewRandom()

	if err != nil {
		srv.logger.Error("failed to generate new random uuid", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not create account")
	}

	note := models.Note{ID: id.String(), AuthorId: noteRequest.AuthorId, Title: noteRequest.Title, Blocks: noteRequest.Blocks, CreationDate: time.Now().UTC(), ModificationDate: time.Now().UTC()}

	_, err = srv.db.Collection("notes").InsertOne(ctx, note)
	if err != nil {
		srv.logger.Error("mongo insert note failed", zap.Error(err), zap.String("note name", note.AuthorId))
		return nil, status.Error(codes.Internal, "could not create note")
	}
	return &note, nil
}

func (srv *notesRepository) Get(ctx context.Context, noteId *string) (*models.Note, error) {
	var note models.Note

	err := srv.db.Collection("notes").FindOne(ctx, buildIdQuery(noteId)).Decode(&note)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "note not found")
		}
		srv.logger.Error("unable to query note", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	_, err = uuid.Parse(note.ID)
	if err != nil {
		srv.logger.Error("failed to convert uuid from string", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not get note")
	}
	return &note, nil
}

func (srv *notesRepository) Delete(ctx context.Context, noteId *string) error {
	delete, err := srv.db.Collection("notes").DeleteOne(ctx, buildIdQuery(noteId))

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

func (srv *notesRepository) Update(ctx context.Context, noteId *string, noteRequest *models.NotePayload) error {
	var note models.Note
	err := copier.Copy(&note, &noteRequest)
	if err != nil {
		srv.logger.Error("failed to copy struct", zap.Error(err))
		return status.Error(codes.Internal, err.Error())
	}
	note.ModificationDate = time.Now().UTC()

	update, err := srv.db.Collection("notes").UpdateOne(ctx, buildIdQuery(noteId), bson.D{{Key: "$set", Value: &note}})
	if err != nil {
		srv.logger.Error("failed to convert object id from hex", zap.Error(err))
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if update.MatchedCount == 0 {
		srv.logger.Error("mongo update note query matched none", zap.String("user_id", *noteId))
		return status.Error(codes.Internal, "could not update note")
	}
	return nil
}

func (srv *notesRepository) List(ctx context.Context, authorId *string) (*[]models.Note, error) {
	cursor, err := srv.db.Collection("notes").Find(ctx, buildAuthodIdQuery(authorId))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "note not found")
		}
		srv.logger.Error("unable to query note", zap.Error(err))
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	var notesResponse []models.Note
	for cursor.Next(context.TODO()) {
		var note models.Note

		if err := cursor.Decode(&note); err != nil {
			srv.logger.Error("unable to parse notes", zap.Error(err))
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		if _, err := uuid.Parse(note.ID); err != nil {
			srv.logger.Error("wrong id from note", zap.Error(err))
			return nil, status.Errorf(codes.Aborted, err.Error())
		}
		notesResponse = append(notesResponse, note)
	}
	return &notesResponse, nil
}

func buildIdQuery(noteId *string) bson.M {
	query := bson.M{}
	if *noteId != "" {
		query["_id"] = noteId
	}
	return query
}

func buildAuthodIdQuery(authorId *string) bson.M {
	query := bson.M{}
	if *authorId != "" {
		query["authorId"] = authorId
	}
	return query
}
