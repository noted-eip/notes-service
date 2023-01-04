// Package memory is an in-memory implementation of models.AccountsRepository
package memory

import (
	"context"
	"notes-service/models"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
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

func NewBlockDatabaseSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"block": {
				Name: "block",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"note_id": {
						Name:    "note_id",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "NoteId"},
					},
					"type": {
						Name:    "type",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Type"},
					},
					"index": {
						Name:    "index",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Index"},
					},
					"content": {
						Name:    "content",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Content"},
					},
				},
			},
		},
	}
}

func (srv *blocksRepository) GetBlock(ctx context.Context, blockId string) (*models.Block, error) {
	txn := srv.db.DB.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("block", "id", blockId)

	if err != nil {
		srv.logger.Error("unable to query block", zap.Error(err))
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound, "block not found")
	}
	return raw.(*models.Block), nil
}

func (srv *blocksRepository) GetBlocks(ctx context.Context, noteId string) ([]*models.Block, error) {
	var blocks []*models.Block

	txn := srv.db.DB.Txn(false)

	it, err := txn.Get("block", "note_id", noteId)
	if err != nil {
		srv.logger.Error("unable to list blocks", zap.Error(err))
		return nil, err
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		blocks = append(blocks, obj.(*models.Block))
	}

	return blocks, nil
}

func (srv *blocksRepository) Create(ctx context.Context, blockRequest *models.Block) (*string, error) {
	txn := srv.db.DB.Txn(true)
	defer txn.Abort()

	id, err := uuid.NewRandom()
	if err != nil {
		srv.logger.Error("failed to generate new random uuid", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not create account")
	}
	blockId := id.String()
	block := models.Block{ID: id.String(), NoteId: blockRequest.NoteId, Type: blockRequest.Type, Index: blockRequest.Index, Content: blockRequest.Content}

	err = txn.Insert("block", &block)
	if err != nil {
		srv.logger.Error("mongo insert block failed", zap.Error(err), zap.String("note id : ", blockRequest.NoteId))
		return nil, status.Errorf(codes.Internal, "could not insert block")
	}
	txn.Commit()
	return &blockId, nil
}

func (srv *blocksRepository) Update(ctx context.Context, blockId string, blockRequest *models.Block) (*models.Block, error) {
	/*update, err := srv.db.Collection("blocks").UpdateOne(ctx, buildIdQuery(blockId), bson.D{{Key: "$set", Value: &blockRequest}})
	if err != nil {
		srv.logger.Error("failed to convert object id from hex", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if update.MatchedCount == 0 {
		srv.logger.Error("mongo update block query matched none", zap.String("block_id : ", blockId))
		return nil, status.Error(codes.Internal, "could not update block")
	}

	return blockRequest, nil
	*/
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *blocksRepository) DeleteBlock(ctx context.Context, blockId string) error {
	txn := srv.db.DB.Txn(true)
	defer txn.Abort()

	err := txn.Delete("block", models.Block{ID: blockId})

	if err == memdb.ErrNotFound {
		srv.logger.Error("unable to find block", zap.Error(err))
	}
	if err != nil {
		srv.logger.Error("delete block db query failed", zap.Error(err))
		return status.Error(codes.Internal, "could not delete block")
	}
	return nil
}

func (srv *blocksRepository) DeleteBlocks(ctx context.Context, noteId string) error {
	txn := srv.db.DB.Txn(true)
	defer txn.Abort()

	//err := txn.Delete("block", models.Block{NoteId: noteId})
	_, err := txn.DeleteAll("block", "note_id", noteId)

	if err != nil {
		srv.logger.Error("delete blocks db query failed", zap.Error(err))
		return status.Error(codes.Internal, "could not delete blocks")
	}
	return nil
}
