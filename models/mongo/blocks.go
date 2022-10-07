package mongo

import (
	"context"
	"errors"
	"notes-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type block struct {
	ID      string `json:"id" bson:"_id,omitempty"`
	NoteId  string `json:"noteId" bson:"noteId,omitempty"`
	Type    uint32 `json:"type" bson:"type,omitempty"`
	Content string `json:"content" bson:"content,omitempty"`
}

type blockWithIndex struct {
	ID      string `json:"id" bson:"_id,omitempty"`
	NoteId  string `json:"noteId" bson:"noteId,omitempty"`
	Type    uint32 `json:"type" bson:"type,omitempty"`
	Index   uint32 `json:"index" bson:"index,omitempty"`
	Content string `json:"content" bson:"content,omitempty"`
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

func (srv *blocksRepository) GetByFilter(ctx context.Context, filter *models.BlockFilter) (*models.BlockWithIndex, error) {
	var block blockWithIndex

	err := srv.db.Collection("notes").FindOne(ctx, buildBlockQuery(filter)).Decode(&block)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "block not found")
		}
		srv.logger.Error("unable to query block", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &models.BlockWithIndex{ID: block.ID, NoteId: filter.NoteId, Type: block.Type, Index: block.Index, Content: block.Content}, nil
}

func (srv *blocksRepository) GetAllById(ctx context.Context, filter *models.BlockFilter) ([]*models.BlockWithIndex, error) {
	return nil, nil
}

func (srv *blocksRepository) Create(ctx context.Context, blockRequest *models.BlockWithIndex) error {
	return nil
}

func (srv *blocksRepository) Update(ctx context.Context, filter *models.BlockFilter, blockRequest *models.BlockWithIndex) (*models.BlockWithIndex, error) {
	return nil, nil
}

func (srv *blocksRepository) Delete(ctx context.Context, filter *models.BlockFilter) error {
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
