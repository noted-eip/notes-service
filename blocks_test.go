package main

import (
	"context"
	"notes-service/memory"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BlocksAPISuite struct {
	suite.Suite
	srv *blocksService
}

func TestBlocksService(t *testing.T) {
	suite.Run(t, &BlocksAPISuite{})
}

func (s *BlocksAPISuite) SetupSuite() {
	logger := newLoggerOrFail(s.T())
	db := newBlocksDatabaseOrFail(s.T(), logger)

	s.srv = &blocksService{
		//auth:      auth.NewService(genKeyOrFail(s.T())),
		logger: logger,
		repo:   memory.NewBlocksRepository(db, logger),
	}
}

func TestBlocksServiceInsertBlockShouldReturnNil(t *testing.T) {
	srv := blocksService{}

	res, err := srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.InvalidArgument)
	require.Nil(t, res)
}

func TestBlocksServiceInsertBlockShouldReturnBlock(t *testing.T) {
	srv := blocksService{}

	res, err := srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{
		Block: &notespb.Block{
			Type: notespb.Block_TYPE_BULLET_POINT,
			Data: &notespb.Block_BulletPoint{},
		},
		Index:  1,
		NoteId: 1,
	})
	require.Nil(t, err)
	require.NotNil(t, res)
}

func TestBlocksServiceUpdateBlock(t *testing.T) {
	srv := blocksService{}

	res, err := srv.UpdateBlock(context.TODO(), &notespb.UpdateBlockRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.Unimplemented)
	require.Nil(t, res)
}

func TestBlocksServiceDeleteBlock(t *testing.T) {
	srv := blocksService{}

	res, err := srv.DeleteBlock(context.TODO(), &notespb.DeleteBlockRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.Unimplemented)
	require.Nil(t, res)
}

func newBlocksDatabaseOrFail(t *testing.T, logger *zap.Logger) *memory.Database {
	db, err := memory.NewDatabase(context.Background(), newBlockDatabaseSchema(), logger)
	require.NoError(t, err, "could not instantiate in-memory database")
	return db
}

func newBlockDatabaseSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"block": {
				Name: "block",
				Indexes: map[string]*memdb.IndexSchema{
					"note_id": {
						Name:    "note_id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "NoteId"},
					},
					"type": {
						Name:    "type",
						Unique:  true,
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
