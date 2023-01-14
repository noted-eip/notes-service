package main

import (
	"context"
	"notes-service/auth"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBlocksAPI(t *testing.T) {
	if os.Getenv("NOTES_SERVICE_TEST_MONGODB_URI") == "" {
		t.Skipf("Skipping NotesAPI suite, missing NOTES_SERVICE_TEST_MONGODB_URI environment variable.")
		return
	}

	suite.Run(t, &NotesAPISuite{})
}

func (s *NotesAPISuite) TestInsertBlockNoAuth() {
	res, err := s.srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{})

	s.Require().Error(err)
	s.Equal(codes.Unauthenticated, status.Code(err))
	s.Nil(res)
}

func (s *NotesAPISuite) TestCreateBlockValidator() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)
	res, err := s.srv.InsertBlock(ctx, &notespb.InsertBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

func (s *NotesAPISuite) TestInsertBlockShouldReturnBlock() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)

	token, err := s.srv.auth.TokenFromContext(ctx)
	s.Require().NoError(err)
	userId := token.UserID.String()

	resCreateNote, err := s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: userId,
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Require().NoError(err)

	blokcContent := "c-test-content"
	blockType := notespb.Block_TYPE_PARAGRAPH
	res, err := s.srv.InsertBlock(ctx, &notespb.InsertBlockRequest{
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
}

func (s *NotesAPISuite) TestUpdateBlockNoAuth() {
	res, err := s.srv.UpdateBlock(context.TODO(), &notespb.UpdateBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.Unauthenticated)
	s.Nil(res)
}

func (s *NotesAPISuite) TestUpdateBlockValidator() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)
	res, err := s.srv.UpdateBlock(ctx, &notespb.UpdateBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
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
}*/

func (s *NotesAPISuite) TestDeleteBlockNoAuth() {
	res, err := s.srv.DeleteBlock(context.TODO(), &notespb.DeleteBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.Unauthenticated)
	s.Nil(res)
}

func (s *NotesAPISuite) TestDeleteBlockValidator() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)
	res, err := s.srv.DeleteBlock(ctx, &notespb.DeleteBlockRequest{})
	s.Require().Error(err)
	s.Equal(status.Code(err), codes.InvalidArgument)
	s.Nil(res)
}

func (s *NotesAPISuite) TestDeleteBlockShouldReturnNoError() {
	generatedUuid, err := uuid.NewRandom()
	s.Require().NoError(err)
	ctx, err := s.auth.ContextWithToken(context.TODO(), &auth.Token{UserID: generatedUuid})
	s.Require().NoError(err)
	//get userId

	resCreateNote, err := s.srv.CreateNote(ctx, &notespb.CreateNoteRequest{
		Note: &notespb.Note{
			AuthorId: generatedUuid.String(),
			Title:    "ci-test",
			Blocks:   nil,
		},
	})
	s.Require().NoError(err)

	resInsertBlock, err := s.srv.InsertBlock(ctx, &notespb.InsertBlockRequest{
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

	res, err := s.srv.DeleteBlock(ctx, &notespb.DeleteBlockRequest{
		Id: resInsertBlock.Block.Id,
	})
	s.Require().NoError(err)
	s.Nil(res)
}
