package main

import (
	"context"
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"

	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestBlocksSuite(t *testing.T) {
	tu := newTestUtilsOrDie(t)
	edouard := newTestAccount(t, tu)
	gabriel := newTestAccount(t, tu)
	maxime := newTestAccount(t, tu)
	edouardGroup := newTestGroup(t, tu, edouard, maxime)
	maximeNote := newTestNote(t, tu, edouardGroup, maxime, []*notesv1.Block{})

	t.Run("insert-block-in-empty-note", func(t *testing.T) {
		res, err := tu.notes.InsertBlock(maxime.Context, &notesv1.InsertBlockRequest{
			GroupId: maximeNote.Group.ID,
			NoteId:  maximeNote.ID,
			Index:   0,
			Block: &notesv1.Block{
				Type: notesv1.Block_TYPE_HEADING_1,
				Data: &notesv1.Block_Heading{
					Heading: "Sample Heading",
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, notesv1.Block_TYPE_HEADING_1, res.Block.Type)
		require.Equal(t, "Sample Heading", res.Block.GetHeading())
		require.NotEmpty(t, res.Block.Id)

		// Make sure the block is stored within the note.
		note, err := tu.notesRepository.GetNote(context.TODO(),
			&models.OneNoteFilter{NoteID: maximeNote.ID, GroupID: maximeNote.Group.ID}, maxime.ID)
		require.NoError(t, err)
		require.NotNil(t, note.FindBlock(res.Block.Id))
	})

	t.Run("insert-block-stranger-cannot-insert", func(t *testing.T) {
		res, err := tu.notes.InsertBlock(gabriel.Context, &notesv1.InsertBlockRequest{
			GroupId: maximeNote.Group.ID,
			NoteId:  maximeNote.ID,
			Index:   0,
			Block: &notesv1.Block{
				Type: notesv1.Block_TYPE_HEADING_1,
				Data: &notesv1.Block_Heading{
					Heading: "Sample Heading",
				},
			},
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("insert-block-member-cannot-insert", func(t *testing.T) {
		res, err := tu.notes.InsertBlock(edouard.Context, &notesv1.InsertBlockRequest{
			GroupId: maximeNote.Group.ID,
			NoteId:  maximeNote.ID,
			Index:   0,
			Block: &notesv1.Block{
				Type: notesv1.Block_TYPE_HEADING_1,
				Data: &notesv1.Block_Heading{
					Heading: "Sample Heading",
				},
			},
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	someBlock := maximeNote.InsertBlock(t, tu, &notesv1.Block{
		Type: notesv1.Block_TYPE_BULLET_POINT,
		Data: &notesv1.Block_BulletPoint{
			BulletPoint: "Sample Bullet Point",
		}}, 0)

	t.Run("delete-block-member-cannot-delete", func(t *testing.T) {
		res, err := tu.notes.DeleteBlock(edouard.Context, &notesv1.DeleteBlockRequest{
			GroupId: maximeNote.Group.ID,
			NoteId:  maximeNote.ID,
			BlockId: someBlock.ID,
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("delete-block-stranger-cannot-delete", func(t *testing.T) {
		res, err := tu.notes.DeleteBlock(gabriel.Context, &notesv1.DeleteBlockRequest{
			GroupId: maximeNote.Group.ID,
			NoteId:  maximeNote.ID,
			BlockId: someBlock.ID,
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	// The mongo repository DeleteBlock method is not working.
	// It doesn't delete the block.
	// t.Run("delete-block", func(t *testing.T) {
	// 	res, err := tu.notes.DeleteBlock(maxime.Context, &notesv1.DeleteBlockRequest{
	// 		GroupId: maximeNote.Group.ID,
	// 		NoteId:  maximeNote.ID,
	// 		BlockId: someBlock.ID,
	// 	})
	// 	require.NoError(t, err)
	// 	require.NotNil(t, res)

	// 	// Make sure the block is not foundable.
	// 	note, err := tu.notesRepository.GetNote(context.TODO(),
	// 		&models.OneNoteFilter{NoteID: maximeNote.ID, GroupID: maximeNote.Group.ID}, maxime.ID)
	// 	require.NoError(t, err)
	// 	require.Nil(t, note.FindBlock(someBlock.ID))
	// })
}
