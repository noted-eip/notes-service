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

func (srv *blocksRepository) GetBlock(ctx context.Context, blockId *string) (*models.BlockWithIndex, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) GetBlocks(ctx context.Context, noteId *string) ([]*models.BlockWithIndex, error) {
	blockCursor, err := srv.db.Collection("blocks").Find(ctx, buildBlockQuery(noteId))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "blocks not found")
		}
		srv.logger.Error("unable to query blocks", zap.Error(err))
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	//convert blocks from mongo to []*BlockWithIndex
	var blocks []bson.M
	if err := blockCursor.All(context.TODO(), &blocks); err != nil {
		srv.logger.Error("unable to parse blocks", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	blocksResponse := make([]*models.BlockWithIndex, blockCursor.RemainingBatchLength())
	for index, block := range blocks {
		id, err := uuid.Parse(block["_id"].(string))
		if err != nil {
			srv.logger.Error("unable to retrieve id of the block", zap.Error(err))
			return nil, status.Errorf(codes.Aborted, err.Error())
		}
		blocksResponse[index] = &models.BlockWithIndex{ID: id.String(), NoteId: *noteId, Type: uint32(block["type"].(int64)), Index: uint32(block["index"].(int64)), Content: block["content"].(string)}
	}

	return blocksResponse, nil
}

func (srv *blocksRepository) Create(ctx context.Context, blockRequest *models.BlockWithIndex) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) Update(ctx context.Context, blockId *string, blockRequest *models.BlockWithIndex) (*models.BlockWithIndex, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) Delete(ctx context.Context, blockId *string) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}

func buildBlockQuery(noteId *string) bson.M {
	query := bson.M{}
	if *noteId != "" {
		query["noteId"] = noteId
	}
	return query
}
