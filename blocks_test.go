package main

import (
	"context"
	"notes-service/models/memory"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/google/uuid"
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

func (s *BlocksAPISuite) InsertBlockShouldReturnNil() {
	res, err := s.srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

/*
func (s *BlocksAPISuite) InsertBlockShouldReturnBlock() {

	slice := []byte{0xFF, 0xFF, 0xFF, 0x7F, 0x7F, 0x7F, 0x7F, 0x7F, 0x7F}
	noteId := binary.LittleEndian.Uint32(slice)

	res, err := s.srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{
		Block: &notespb.Block{
			Type: notespb.Block_TYPE_BULLET_POINT,
			Data: &notespb.Block_BulletPoint{},
		},
		Index:  1,
		NoteId: noteId,
	})
	s.Nil(err)
	s.NotNil(res)
}*/

func (s *BlocksAPISuite) UpdateBlockShouldReturnError() {
	res, err := s.srv.UpdateBlock(context.TODO(), &notespb.UpdateBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.Internal)
	s.Nil(res)
}

func (s *BlocksAPISuite) UpdateBlockShouldReturnNoError() {
	blockId, err := uuid.NewRandom()

	res, err := s.srv.UpdateBlock(context.TODO(), &notespb.UpdateBlockRequest{
		Id: blockId.String(),
		Block: &notespb.Block{
			Type: notespb.Block_TYPE_BULLET_POINT,
			Data: &notespb.Block_BulletPoint{},
		},
		Index: 1,
	})
	s.Nil(err)
	s.Nil(res)
}

func (s *BlocksAPISuite) DeleteBlockShouldReturnError() {
	res, err := s.srv.DeleteBlock(context.TODO(), &notespb.DeleteBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

/*
func (s *BlocksAPISuite) DeleteBlockShouldReturnNoError() {
	id, err := uuid.NewRandom()

	res, err := s.srv.DeleteBlock(context.TODO(), &notespb.DeleteBlockRequest{
		Id: id.String(),
	})
	s.Nil(err)
	s.Nil(res)
}
*/
func newBlocksDatabaseOrFail(t *testing.T, logger *zap.Logger) *memory.Database {
	db, err := memory.NewDatabase(context.Background(), memory.NewBlockDatabaseSchema(), logger)
	require.NoError(t, err, "could not instantiate in-memory database")
	return db
}
