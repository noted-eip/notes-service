package main

import (
	"context"
	"notes-service/models"
	"notes-service/models/memory"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBlocksService(t *testing.T) {
	suite.Run(t, &NotesAPISuite{})
}

func (s *NotesAPISuite) TestInsertBlockNoAuth() {
	res, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{})

	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestCreateBlockValidator() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	res, err := s.srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
	s.srv.auth = saveAuthPackage
}

func (s *NotesAPISuite) TestInsertBlockShouldReturnBlock() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	//get userId
	token, err := s.srv.auth.TokenFromContext(context.TODO())
	s.Require().NoError(err)
	userId := token.UserID.String()

	resCreateNote, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: userId,
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Require().NoError(err)

	blokcContent := "c-test-content"
	blockType := notespb.Block_TYPE_PARAGRAPH
	res, err := s.srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{
		Block: &notespb.Block{
			Type: blockType,
			Data: &notespb.Block_Paragraph{
				Paragraph: blokcContent,
			},
		},
		Index:  1,
		NoteId: resCreateNote.Note.Id,
	})
	s.Require().NoError(err)
	s.NotNil(res)

	var actualBlock models.Block
	convertApiBlockToModelBlock(&actualBlock, res.Block)
	s.Equal(blokcContent, actualBlock.Content)
	s.Equal(blockType, res.Block.Type)
	s.srv.auth = saveAuthPackage
}

func (s *NotesAPISuite) TestUpdateBlockNoAuth() {
	res, err := s.srv.UpdateBlock(context.TODO(), &notespb.UpdateBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.Unauthenticated)
	s.Nil(res)
}

func (s *NotesAPISuite) TestUpdateBlockValidator() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	res, err := s.srv.UpdateBlock(context.TODO(), &notespb.UpdateBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
	s.srv.auth = saveAuthPackage
}

/*
func (s *NotesAPISuite) TestUpdateBlockShouldReturnNoError() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	//get userId
	token, err := s.srv.auth.TokenFromContext(context.TODO())
	s.Require().NoError(err)
	userId := token.UserID.String()

	resCreateNote, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: userId,
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Require().NoError(err)

	resInsertBlock, err := s.srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{
		Block: &notespb.Block{
			Type: notespb.Block_TYPE_PARAGRAPH,
			Data: &notespb.Block_Paragraph{
				Paragraph: "c-test-content",
			},
		},
		Index:  1,
		NoteId: resCreateNote.Note.Id,
	})
	s.Require().NoError(err)

	res, err := s.srv.UpdateBlock(context.TODO(), &notespb.UpdateBlockRequest{
		Id: blockId.String(),
		Block: &notespb.Block{
			Type: notespb.Block_TYPE_PARAGRAPH,
			Data: &notespb.Block_Paragraph{
				Paragraph: "c-test-content-updated",
			},
		},
		Index: 1,
	})
	s.Require().NoError(err)

	var actualBlock models.Block
	convertApiBlockToModelBlock(&actualBlock, res.Block)
	s.Equal("c-test-content-updated", actualBlock.Content)
	s.Equal(notespb.Block_TYPE_PARAGRAPH, res.Block.Type)//switch le type aussi peu être
	s.srv.auth = saveAuthPackage
}*/

func (s *NotesAPISuite) TestDeleteBlockNoAuth() {
	res, err := s.srv.DeleteBlock(context.TODO(), &notespb.DeleteBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.Unauthenticated)
	s.Nil(res)
}

func (s *NotesAPISuite) TestDeleteBlockValidator() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	res, err := s.srv.DeleteBlock(context.TODO(), &notespb.DeleteBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
	s.srv.auth = saveAuthPackage
}

func (s *NotesAPISuite) TestDeleteBlockShouldReturnNoError() {
	saveAuthPackage := s.srv.auth
	s.srv.auth = NewMockService()
	//get userId
	token, err := s.srv.auth.TokenFromContext(context.TODO())
	s.Require().NoError(err)
	userId := token.UserID.String()

	resCreateNote, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: userId,
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Require().NoError(err)

	resInsertBlock, err := s.srv.InsertBlock(context.TODO(), &notespb.InsertBlockRequest{
		Block: &notespb.Block{
			Type: notespb.Block_TYPE_PARAGRAPH,
			Data: &notespb.Block_Paragraph{
				Paragraph: "c-test-content",
			},
		},
		Index:  1,
		NoteId: resCreateNote.Note.Id,
	})
	s.Require().NoError(err)

	res, err := s.srv.DeleteBlock(context.TODO(), &notespb.DeleteBlockRequest{
		Id: resInsertBlock.Block.Id,
	})
	s.Require().NoError(err)
	s.Nil(res)
	s.srv.auth = saveAuthPackage
}

func newBlocksDatabaseOrFail(t *testing.T, logger *zap.Logger) *memory.Database {
	db, err := memory.NewDatabase(context.Background(), memory.NewBlockDatabaseSchema(), logger)
	require.NoError(t, err, "could not instantiate in-memory database")
	return db
}
