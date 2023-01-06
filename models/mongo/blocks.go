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
	logger *zap.Logger
	db     *mongo.Database
}

func NewBlocksRepository(db *mongo.Database, logger *zap.Logger) models.BlocksRepository {
	return &blocksRepository{
		logger: logger,
		db:     db,
	}
}

func (srv *blocksRepository) GetBlock(ctx context.Context, blockId string) (*models.Block, error) {
	var block models.Block

	err := srv.db.Collection("blocks").FindOne(ctx, buildIdQuery(blockId)).Decode(&block)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "block not found")
		}
		srv.logger.Error("unable to query block", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &models.Block{ID: block.ID, NoteId: block.NoteId, Type: block.Type, Index: block.Index, Content: block.Content}, nil
}

func (srv *blocksRepository) GetBlocks(ctx context.Context, noteId string) ([]*models.Block, error) {
	blockCursor, err := srv.db.Collection("blocks").Find(ctx, buildNoteIdQuery(noteId))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "blocks not found")
		}
		srv.logger.Error("unable to query blocks", zap.Error(err))
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	//convert blocks from mongo to []*Block
	var blocks []bson.M
	if err := blockCursor.All(context.TODO(), &blocks); err != nil {
		srv.logger.Error("unable to parse blocks", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	blocksResponse := make([]*models.Block, blockCursor.RemainingBatchLength())
	for index, block := range blocks {
		id, err := uuid.Parse(block["_id"].(string))
		if err != nil {
			srv.logger.Error("unable to retrieve id of the block", zap.Error(err))
			return nil, status.Errorf(codes.Aborted, err.Error())
		}
		blocksResponse[index] = &models.Block{ID: id.String(), NoteId: noteId, Type: uint32(block["type"].(int64)), Index: uint32(block["index"].(int64)), Content: block["content"].(string)}
	}

	return blocksResponse, nil
}

func (srv *blocksRepository) Create(ctx context.Context, blockRequest *models.Block) (*string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		srv.logger.Error("failed to generate new random uuid", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not create account")
	}
	blockId := id.String()
	block := models.Block{ID: id.String(), NoteId: blockRequest.NoteId, Type: blockRequest.Type, Index: blockRequest.Index, Content: blockRequest.Content}

	_, err = srv.db.Collection("blocks").InsertOne(ctx, block)
	if err != nil {
		srv.logger.Error("mongo insert block failed", zap.Error(err), zap.String("note id : ", blockRequest.NoteId))
		return nil, status.Error(codes.Internal, "could not insert block")
	}
	return &blockId, nil
}

func (srv *blocksRepository) Update(ctx context.Context, blockId string, blockRequest *models.Block) (*models.Block, error) {
	//index update ?
	update, err := srv.db.Collection("blocks").UpdateOne(ctx, buildIdQuery(blockId), bson.D{{Key: "$set", Value: &blockRequest}})

	if err != nil {
		srv.logger.Error("failed to convert object id from hex", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if update.MatchedCount == 0 {
		srv.logger.Error("mongo update block query matched none", zap.String("block_id : ", blockId))
		return nil, status.Error(codes.Internal, "could not update block")
	}

	return blockRequest, nil
}

func (srv *blocksRepository) DeleteBlock(ctx context.Context, blockId string) error {
	delete, err := srv.db.Collection("blocks").DeleteOne(ctx, buildIdQuery(blockId))

	if err != nil {
		srv.logger.Error("delete block db query failed", zap.Error(err))
		return status.Error(codes.Internal, "could not delete block")
	}
	if delete.DeletedCount == 0 {
		srv.logger.Info("mongo delete block matched none", zap.String("block_id", blockId))
		return status.Error(codes.Internal, "could not delete block")
	}
	return nil
}

func (srv *blocksRepository) DeleteBlocks(ctx context.Context, noteId string) error {
	delete, err := srv.db.Collection("blocks").DeleteMany(ctx, buildNoteIdQuery(noteId))

	if err != nil {
		srv.logger.Error("delete blocks db query failed", zap.Error(err))
		return status.Error(codes.Internal, "could not delete blocks")
	}
	if delete.DeletedCount == 0 {
		srv.logger.Info("mongo delete block matched none", zap.String("note_id", noteId))
		return status.Error(codes.Internal, "could not delete block")
	}
	return nil
}

func buildNoteIdQuery(noteId string) bson.M {
	query := bson.M{}
	if noteId != "" {
		query["noteId"] = noteId
	}
	return query
}
