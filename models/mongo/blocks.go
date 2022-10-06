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

type block struct {
	ID      string  `json:"id" bson:"_id,omitempty"`
	NoteId  string  `json:"noteId" bson:"noteId,omitempty"`
	Type    uint32  `json:"type" bson:"type,omitempty"`
	Content *string `json:"content" bson:"content,omitempty"`
}

type blockWithIndex struct {
	ID      string  `json:"id" bson:"_id,omitempty"`
	NoteId  string  `json:"noteId" bson:"noteId,omitempty"`
	Type    uint32  `json:"type" bson:"type,omitempty"`
	Index   uint32  `json:"inxed" bson:"inxed,omitempty"`
	Content *string `json:"content" bson:"content,omitempty"`
}

type blocksRepository struct {
	logger          *zap.Logger
	db              *mongo.Database
	notesRepository models.NotesRepository
}

func NewBlocksRepository(db *mongo.Database, logger *zap.Logger, notesRepository models.NotesRepository) models.BlocksRepository {
	return &blocksRepository{
		logger:          logger,
		db:              db,
		notesRepository: notesRepository,
	}
}

func (srv *blocksRepository) Create(ctx context.Context, blockRequest *models.BlockWithIndex) error {
	id, err := uuid.NewRandom()
	if err != nil {
		srv.logger.Error("failed to generate new random uuid", zap.Error(err))
		return status.Errorf(codes.Internal, "could not create account")
	}

	block := blockWithIndex{ID: id.String(), NoteId: blockRequest.NoteId, Type: blockRequest.Type, Index: blockRequest.Index, Content: blockRequest.Content}

	_, err = srv.db.Collection("blocks").InsertOne(ctx, block)
	if err != nil {
		srv.logger.Error("mongo insert block failed", zap.Error(err), zap.String("note id : ", blockRequest.NoteId))
		return status.Errorf(codes.Internal, "could not insert block")
	}
	return nil
}

func (srv *blocksRepository) Update(ctx context.Context, filter *models.BlockFilter, blockRequest *models.BlockWithIndex) (*models.BlockWithIndex, error) {
	update, err := srv.db.Collection("blocks").UpdateOne(ctx, buildBlockQuery(filter), bson.D{{Key: "$set", Value: &blockRequest}})

	if err != nil {
		srv.logger.Error("failed to convert object id from hex", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if update.MatchedCount == 0 {
		srv.logger.Error("mongo update note query matched none", zap.String("block_id : ", filter.BlockId))
		return nil, status.Errorf(codes.Internal, "could not update block")
	}

	return blockRequest, nil
}

func (srv *blocksRepository) Delete(ctx context.Context, filter *models.BlockFilter) error {
	delete, err := srv.db.Collection("blocks").DeleteOne(ctx, buildBlockQuery(filter))

	if err != nil {
		srv.logger.Error("delete note db query failed", zap.Error(err))
		return status.Errorf(codes.Internal, "could not delete note")
	}
	if delete.DeletedCount == 0 {
		srv.logger.Info("mongo delete block matched none", zap.String("block_id", filter.BlockId))
		return status.Errorf(codes.Internal, "could not delete block")
	}
	return nil
}

func buildBlockQuery(filter *models.BlockFilter) bson.M {
	query := bson.M{}
	if filter.BlockId != "" {
		query["_id"] = filter.BlockId
	}
	if filter.NoteId != "" {
		query["noteId"] = filter.NoteId
	}
	return query
}
