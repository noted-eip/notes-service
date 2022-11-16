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

func (srv *blocksRepository) GetBlock(ctx context.Context, blockId *string) (*models.BlockWithIndex, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) GetBlocks(ctx context.Context, noteId *string) ([]*models.BlockWithIndex, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) Create(ctx context.Context, blockRequest *models.BlockWithIndex) (*string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		srv.logger.Error("failed to generate new random uuid", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not create account")
	}
	blockId := id.String()
	block := models.BlockWithIndex{ID: id.String(), NoteId: blockRequest.NoteId, Type: blockRequest.Type, Index: blockRequest.Index, Content: blockRequest.Content}

	_, err = srv.db.Collection("blocks").InsertOne(ctx, block)
	if err != nil {
		srv.logger.Error("mongo insert block failed", zap.Error(err), zap.String("note id : ", blockRequest.NoteId))
		return nil, status.Error(codes.Internal, "could not insert block")
	}
	return &blockId, nil
}

func (srv *blocksRepository) Update(ctx context.Context, blockId *string, blockRequest *models.BlockWithIndex) (*models.BlockWithIndex, error) {
	update, err := srv.db.Collection("blocks").UpdateOne(ctx, buildBlockQuery(blockId), bson.D{{Key: "$set", Value: &blockRequest}})

	if err != nil {
		srv.logger.Error("failed to convert object id from hex", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if update.MatchedCount == 0 {
		srv.logger.Error("mongo update note query matched none", zap.String("block_id : ", *blockId))
		return nil, status.Error(codes.Internal, "could not update block")
	}

	return blockRequest, nil
}

//delete one block with BlockId
func (srv *blocksRepository) DeleteBlock(ctx context.Context, blockId *string) error {
	delete, err := srv.db.Collection("blocks").DeleteOne(ctx, buildBlockQuery(blockId))

	if err != nil {
		srv.logger.Error("delete note db query failed", zap.Error(err))
		return status.Error(codes.Internal, "could not delete note")
	}
	if delete.DeletedCount == 0 {
		srv.logger.Info("mongo delete block matched none", zap.String("block_id", *blockId))
		return status.Error(codes.Internal, "could not delete block")
	}
	return nil
}

//delete multiple blocks with NoteId
func (srv *blocksRepository) DeleteBlocks(ctx context.Context, noteId *string) error {
	delete, err := srv.db.Collection("blocks").DeleteOne(ctx, buildBlocksQuery(noteId))

	if err != nil {
		srv.logger.Error("delete note db query failed", zap.Error(err))
		return status.Error(codes.Internal, "could not delete note")
	}
	if delete.DeletedCount == 0 {
		srv.logger.Info("mongo delete block matched none", zap.String("note_id", *noteId))
		return status.Error(codes.Internal, "could not delete block")
	}
	return nil
}

func buildBlockQuery(blockId *string) bson.M {
	query := bson.M{}
	if *blockId != "" {
		query["_id"] = blockId
	}
	return query
}

func buildBlocksQuery(noteId *string) bson.M {
	query := bson.M{}
	if *noteId != "" {
		query["noteId"] = noteId
	}
	return query
}
