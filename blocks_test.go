package main

import (
	"context"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBlocksService(t *testing.T) {
	suite.Run(t, &NotesAPISuite{})
}

/*
func (s *NotesAPISuite) InsertBlockShouldReturnNil() {
	res, err := s.srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}*/

/*
func (s *NotesAPISuite) InsertBlockShouldReturnBlock() {

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

func (s *NotesAPISuite) TestUpdateBlockShouldReturnError() {
	res, err := s.srv.UpdateBlock(context.TODO(), &notespb.UpdateBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.Internal)
	s.Nil(res)
}

func (s *NotesAPISuite) TestUpdateBlockShouldReturnNoError() {
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

func (s *NotesAPISuite) TestDeleteBlockShouldReturnError() {
	res, err := s.srv.DeleteBlock(context.TODO(), &notespb.DeleteBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

/*
	func (s *NotesAPISuite) DeleteBlockShouldReturnNoError() {
		id, err := uuid.NewRandom()

		res, err := s.srv.DeleteBlock(context.TODO(), &notespb.DeleteBlockRequest{
			Id: id.String(),
		})
		s.Nil(err)
		s.Nil(res)
	}
*/
