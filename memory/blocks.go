// Package memory is an in-memory implementation of models.AccountsRepository
package memory

import (
	"context"
	"notes-service/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type blocksRepository struct {
	logger *zap.Logger
	db     *Database
}

func NewBlocksRepository(db *Database, logger *zap.Logger) models.BlocksRepository {
	return &blocksRepository{
		logger: logger.Named("memory").Named("blocks"),
		db:     db,
	}
}

func (srv *blocksRepository) GetBlock(ctx context.Context, blockId *string) (*models.BlockWithIndex, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) GetBlocks(ctx context.Context, noteId *string) ([]*models.BlockWithIndex, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) Create(ctx context.Context, blockRequest *models.BlockWithIndex) (*string, error) {
	txn := srv.db.DB.Txn(true)
	defer txn.Abort()

	id, err := uuid.NewRandom()
	if err != nil {
		srv.logger.Error("failed to generate new random uuid", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not create account")
	}
	blockId := id.String()
	block := models.BlockWithIndex{ID: id.String(), NoteId: blockRequest.NoteId, Type: blockRequest.Type, Index: blockRequest.Index, Content: blockRequest.Content}

	err = txn.Insert("note", block)
	if err != nil {
		srv.logger.Error("mongo insert block failed", zap.Error(err), zap.String("note id : ", blockRequest.NoteId))
		return nil, status.Errorf(codes.Internal, "could not insert block")
	}
	return &blockId, nil
}

func (srv *blocksRepository) Update(ctx context.Context, blockId *string, blockRequest *models.BlockWithIndex) (*models.BlockWithIndex, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) DeleteBlock(ctx context.Context, blockId *string) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) DeleteBlocks(ctx context.Context, noteId *string) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}
