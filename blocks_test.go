package main

import (
	"context"
	"notes-service/memory"
	notespb "notes-service/protorepo/noted/notes/v1"
	"strconv"
	"testing"

	"github.com/google/uuid"
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

func (s *BlocksAPISuite) TestBlocksServiceInsertBlockShouldReturnNil() {
	res, err := s.srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

func (s *BlocksAPISuite) TestBlocksServiceInsertBlockShouldReturnBlock() {
	noteId, err := uuid.NewRandom()
	noteIdInt, err := strconv.Atoi(noteId.String())
	res, err := s.srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{
		Block: &notespb.Block{
			Type: notespb.Block_TYPE_BULLET_POINT,
			Data: &notespb.Block_BulletPoint{},
		},
		Index:  1,
		NoteId: uint32(noteIdInt),
	})
	s.Nil(err)
	s.NotNil(res)
}

func (s *BlocksAPISuite) TestBlocksServiceUpdateBlock() {
	res, err := s.srv.UpdateBlock(context.TODO(), &notespb.UpdateBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.Unimplemented)
	s.Nil(res)
}

func (s *BlocksAPISuite) TestBlocksServiceDeleteBlock() {
	res, err := s.srv.DeleteBlock(context.TODO(), &notespb.DeleteBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.Unimplemented)
	s.Nil(res)
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
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
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
