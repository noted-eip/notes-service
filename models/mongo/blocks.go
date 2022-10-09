package mongo

import (
	"context"
	"notes-service/models"

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
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) GetAllById(ctx context.Context, filter *models.BlockFilter) ([]*models.BlockWithIndex, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) Create(ctx context.Context, blockRequest *models.BlockWithIndex) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) Update(ctx context.Context, filter *models.BlockFilter, blockRequest *models.BlockWithIndex) (*models.BlockWithIndex, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) Delete(ctx context.Context, filter *models.BlockFilter) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}
