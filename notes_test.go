package main

import (
	"context"
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"

	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestNotesSuite(t *testing.T) {
	tu := newTestUtilsOrDie(t)
	edouard := newTestAccount(t, tu)
	gabriel := newTestAccount(t, tu)
	maxime := newTestAccount(t, tu)
	edouardGroup := newTestGroup(t, tu, edouard, maxime)
	maximeGroup := newTestGroup(t, tu, maxime, edouard)

	t.Run("create-note", func(t *testing.T) {
		res, err := tu.notes.CreateNote(edouard.Context, &notesv1.CreateNoteRequest{
			GroupId: edouardGroup.ID,
			Title:   "My New Note",
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "My New Note", res.Note.Title)
	})

	t.Run("create-note-permission-denied", func(t *testing.T) {
		res, err := tu.notes.CreateNote(gabriel.Context, &notesv1.CreateNoteRequest{
			GroupId: edouardGroup.ID,
			Title:   "My New Note",
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("create-note-with-blocks", func(t *testing.T) {
		res, err := tu.notes.CreateNote(edouard.Context, &notesv1.CreateNoteRequest{
			GroupId: edouardGroup.ID,
			Title:   "My New Note",
			Blocks: []*notesv1.Block{
				{
					Type: notesv1.Block_TYPE_BULLET_POINT,
					Data: &notesv1.Block_BulletPoint{
						BulletPoint: "Sample Bullet Point",
					},
				},
				{
					Type: notesv1.Block_TYPE_CODE,
					Data: &notesv1.Block_Code_{
						Code: &notesv1.Block_Code{
							Snippet: "Sample Snippet",
							Lang:    "Sample Lang",
						},
					},
				},
				{
					Type: notesv1.Block_TYPE_HEADING_1,
					Data: &notesv1.Block_Heading{
						Heading: "Sample Heading",
					},
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Note.Blocks, 3)
		require.NotEmpty(t, res.Note.Blocks[0].Id)
		require.NotEmpty(t, res.Note.Blocks[1].Id)
		require.NotEmpty(t, res.Note.Blocks[2].Id)
		require.Equal(t, notesv1.Block_TYPE_BULLET_POINT, res.Note.Blocks[0].Type)
		require.Equal(t, notesv1.Block_TYPE_CODE, res.Note.Blocks[1].Type)
		require.Equal(t, notesv1.Block_TYPE_HEADING_1, res.Note.Blocks[2].Type)
		require.Equal(t, "Sample Bullet Point", res.Note.Blocks[0].GetBulletPoint())
		require.Equal(t, "Sample Snippet", res.Note.Blocks[1].GetCode().Snippet)
		require.Equal(t, "Sample Lang", res.Note.Blocks[1].GetCode().Lang)
		require.Equal(t, "Sample Heading", res.Note.Blocks[2].GetHeading())
	})

	t.Run("create-note-with-all-block-types", func(t *testing.T) {
		res, err := tu.notes.CreateNote(edouard.Context, &notesv1.CreateNoteRequest{
			GroupId: edouardGroup.ID,
			Title:   "My New Note",
			Blocks: []*notesv1.Block{
				{
					Type: notesv1.Block_TYPE_NUMBER_POINT,
					Data: &notesv1.Block_NumberPoint{
						NumberPoint: "Sample Number Point",
					},
				},
				{
					Type: notesv1.Block_TYPE_BULLET_POINT,
					Data: &notesv1.Block_BulletPoint{
						BulletPoint: "Sample Bullet Point",
					},
				},
				{
					Type: notesv1.Block_TYPE_MATH,
					Data: &notesv1.Block_Math{
						Math: "Sample Math",
					},
				},
				{
					Type: notesv1.Block_TYPE_CODE,
					Data: &notesv1.Block_Code_{
						Code: &notesv1.Block_Code{
							Snippet: "Sample Snippet",
							Lang:    "Sample Lang",
						},
					},
				},
				{
					Type: notesv1.Block_TYPE_IMAGE,
					Data: &notesv1.Block_Image_{
						Image: &notesv1.Block_Image{
							Url:     "Sample Url",
							Caption: "Sample Caption",
						},
					},
				},
				{
					Type: notesv1.Block_TYPE_HEADING_1,
					Data: &notesv1.Block_Heading{
						Heading: "Sample Heading 1",
					},
				},
				{
					Type: notesv1.Block_TYPE_HEADING_2,
					Data: &notesv1.Block_Heading{
						Heading: "Sample Heading 2",
					},
				},
				{
					Type: notesv1.Block_TYPE_HEADING_3,
					Data: &notesv1.Block_Heading{
						Heading: "Sample Heading 3",
					},
				},
				{
					Type: notesv1.Block_TYPE_PARAGRAPH,
					Data: &notesv1.Block_Paragraph{
						Paragraph: "Sample Paragraph",
					},
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Note.Blocks, 9)
		require.Equal(t, notesv1.Block_TYPE_NUMBER_POINT, res.Note.Blocks[0].Type)
		require.Equal(t, notesv1.Block_TYPE_BULLET_POINT, res.Note.Blocks[1].Type)
		require.Equal(t, notesv1.Block_TYPE_MATH, res.Note.Blocks[2].Type)
		require.Equal(t, notesv1.Block_TYPE_CODE, res.Note.Blocks[3].Type)
		require.Equal(t, notesv1.Block_TYPE_IMAGE, res.Note.Blocks[4].Type)
		require.Equal(t, notesv1.Block_TYPE_HEADING_1, res.Note.Blocks[5].Type)
		require.Equal(t, notesv1.Block_TYPE_HEADING_2, res.Note.Blocks[6].Type)
		require.Equal(t, notesv1.Block_TYPE_HEADING_3, res.Note.Blocks[7].Type)
		require.Equal(t, notesv1.Block_TYPE_PARAGRAPH, res.Note.Blocks[8].Type)
		require.Equal(t, "Sample Number Point", res.Note.Blocks[0].GetNumberPoint())
		require.Equal(t, "Sample Bullet Point", res.Note.Blocks[1].GetBulletPoint())
		require.Equal(t, "Sample Math", res.Note.Blocks[2].GetMath())
		require.Equal(t, "Sample Lang", res.Note.Blocks[3].GetCode().Lang)
		require.Equal(t, "Sample Snippet", res.Note.Blocks[3].GetCode().Snippet)
		require.Equal(t, "Sample Caption", res.Note.Blocks[4].GetImage().Caption)
		require.Equal(t, "Sample Url", res.Note.Blocks[4].GetImage().Url)
		require.Equal(t, "Sample Heading 1", res.Note.Blocks[5].GetHeading())
		require.Equal(t, "Sample Heading 2", res.Note.Blocks[6].GetHeading())
		require.Equal(t, "Sample Heading 3", res.Note.Blocks[7].GetHeading())
		require.Equal(t, "Sample Paragraph", res.Note.Blocks[8].GetParagraph())
	})

	t.Run("create-note-with-invalid-blocks", func(t *testing.T) {
		res, err := tu.notes.CreateNote(edouard.Context, &notesv1.CreateNoteRequest{
			GroupId: edouardGroup.ID,
			Title:   "Sample Title",
			Blocks: []*notesv1.Block{
				{
					Type: notesv1.Block_TYPE_HEADING_1,
					Data: &notesv1.Block_Code_{
						Code: &notesv1.Block_Code{
							Snippet: "Sample Snippet",
							Lang:    "Sample Lang",
						},
					},
				},
				{
					Type: notesv1.Block_TYPE_CODE,
					Data: nil,
				},
				{
					Type: notesv1.Block_TYPE_IMAGE,
					Data: nil,
				},
				{
					Type: notesv1.Block_TYPE_CODE,
					Data: &notesv1.Block_Heading{
						Heading: "Sample Heading",
					},
				},
				{},
				{
					Type: notesv1.Block_TYPE_INVALID,
				},
				{
					Data: &notesv1.Block_BulletPoint{
						BulletPoint: "Sample Bullet Point",
					},
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, notesv1.Block_TYPE_HEADING_1, res.Note.Blocks[0].Type)
		require.Equal(t, notesv1.Block_TYPE_CODE, res.Note.Blocks[1].Type)
		require.Equal(t, notesv1.Block_TYPE_IMAGE, res.Note.Blocks[2].Type)
		require.Equal(t, notesv1.Block_TYPE_CODE, res.Note.Blocks[3].Type)
		require.Equal(t, notesv1.Block_TYPE_INVALID, res.Note.Blocks[4].Type)
		require.Equal(t, notesv1.Block_TYPE_INVALID, res.Note.Blocks[5].Type)
		require.Equal(t, notesv1.Block_TYPE_INVALID, res.Note.Blocks[6].Type)
		require.Equal(t, "", res.Note.Blocks[0].GetHeading())
		require.Equal(t, "", res.Note.Blocks[1].GetCode().Lang)
		require.Equal(t, "", res.Note.Blocks[1].GetCode().Snippet)
		require.Equal(t, "", res.Note.Blocks[2].GetImage().Caption)
		require.Equal(t, "", res.Note.Blocks[2].GetImage().Url)
		require.Equal(t, "", res.Note.Blocks[3].GetCode().Lang)
		require.Equal(t, "", res.Note.Blocks[3].GetCode().Snippet)
		require.Equal(t, "", res.Note.Blocks[6].GetBulletPoint())
	})

	edouardNote := newTestNote(t, tu, edouardGroup, edouard, []*notesv1.Block{})

	t.Run("update-note-title", func(t *testing.T) {
		res, err := tu.notes.UpdateNote(edouard.Context, &notesv1.UpdateNoteRequest{
			GroupId: edouardGroup.ID,
			NoteId:  edouardNote.ID,
			Note: &notesv1.Note{
				Title: "Brand New Title",
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"title"},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "Brand New Title", res.Note.Title)
	})

	t.Run("update-note-title-group-member-cannot-edit", func(t *testing.T) {
		res, err := tu.notes.UpdateNote(maxime.Context, &notesv1.UpdateNoteRequest{
			GroupId: edouardGroup.ID,
			NoteId:  edouardNote.ID,
			Note: &notesv1.Note{
				Title: "Brand New Title",
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"title"},
			},
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("update-note-title-stranger-cannot-edit", func(t *testing.T) {
		res, err := tu.notes.UpdateNote(gabriel.Context, &notesv1.UpdateNoteRequest{
			GroupId: edouardGroup.ID,
			NoteId:  edouardNote.ID,
			Note: &notesv1.Note{
				Title: "Brand New Title",
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"title"},
			},
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("update-note-title-group-member-cannot-edit", func(t *testing.T) {
		res, err := tu.notes.UpdateNote(maxime.Context, &notesv1.UpdateNoteRequest{
			GroupId: edouardGroup.ID,
			NoteId:  edouardNote.ID,
			Note: &notesv1.Note{
				Title: "Brand New Title",
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"title"},
			},
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	maximeNote := newTestNote(t, tu, edouardGroup, maxime, nil)

	t.Run("delete-note-member-cannot-delete", func(t *testing.T) {
		res, err := tu.notes.DeleteNote(edouard.Context, &notesv1.DeleteNoteRequest{
			GroupId: maximeNote.Group.ID,
			NoteId:  maximeNote.ID,
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("delete-note-stranger-cannot-delete", func(t *testing.T) {
		res, err := tu.notes.DeleteNote(gabriel.Context, &notesv1.DeleteNoteRequest{
			GroupId: maximeNote.Group.ID,
			NoteId:  maximeNote.ID,
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("delete-note", func(t *testing.T) {
		res, err := tu.notes.DeleteNote(maxime.Context, &notesv1.DeleteNoteRequest{
			GroupId: maximeNote.Group.ID,
			NoteId:  maximeNote.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	// DeleteNote is a repository function, no auth
	_ = newTestNote(t, tu, edouardGroup, edouard, []*notesv1.Block{})
	_ = newTestNote(t, tu, edouardGroup, edouard, []*notesv1.Block{})
	_ = newTestNote(t, tu, maximeGroup, edouard, []*notesv1.Block{})

	t.Run("delete-notes-account", func(t *testing.T) {
		err := tu.notesRepository.DeleteNotes(context.TODO(), &models.ManyNotesFilter{
			AuthorAccountID: edouard.ID,
		})
		require.NoError(t, err)

		notes, err := tu.notesRepository.ListAllNotesInternal(context.TODO(), &models.ManyNotesFilter{
			AuthorAccountID: edouard.ID,
		})
		require.NoError(t, err)
		require.Zero(t, len(notes))
	})

	_ = newTestNote(t, tu, edouardGroup, edouard, []*notesv1.Block{})
	_ = newTestNote(t, tu, edouardGroup, edouard, []*notesv1.Block{})
	_ = newTestNote(t, tu, edouardGroup, maxime, []*notesv1.Block{})
	_ = newTestNote(t, tu, edouardGroup, maxime, []*notesv1.Block{})

	t.Run("delete-notes-group", func(t *testing.T) {
		err := tu.notesRepository.DeleteNotes(context.TODO(), &models.ManyNotesFilter{
			GroupID: edouardGroup.ID,
		})
		require.NoError(t, err)

		notes, err := tu.notesRepository.ListAllNotesInternal(context.TODO(), &models.ManyNotesFilter{
			GroupID: edouardGroup.ID,
		})
		require.NoError(t, err)
		require.Zero(t, len(notes))
	})

	_ = newTestNote(t, tu, edouardGroup, edouard, []*notesv1.Block{})
	_ = newTestNote(t, tu, edouardGroup, edouard, []*notesv1.Block{})
	_ = newTestNote(t, tu, edouardGroup, maxime, []*notesv1.Block{})
	_ = newTestNote(t, tu, edouardGroup, maxime, []*notesv1.Block{})

	t.Run("delete-notes-group-and-account", func(t *testing.T) {
		err := tu.notesRepository.DeleteNotes(context.TODO(), &models.ManyNotesFilter{
			GroupID:         edouardGroup.ID,
			AuthorAccountID: edouard.ID,
		})
		require.NoError(t, err)

		// Check that edouard doesn't have any notes left
		notes, err := tu.notesRepository.ListAllNotesInternal(context.TODO(), &models.ManyNotesFilter{
			AuthorAccountID: edouard.ID,
			GroupID:         edouardGroup.ID,
		})
		require.NoError(t, err)
		require.Zero(t, len(notes))

		// Check that the only remaining notes in the group are from maxime
		notes, err = tu.notesRepository.ListAllNotesInternal(context.TODO(), &models.ManyNotesFilter{
			GroupID: edouardGroup.ID,
		})
		require.NoError(t, err)
		require.Equal(t, len(notes), 2)

		for _, note := range notes {
			require.Equal(t, note.AuthorAccountID, maxime.ID)
		}
	})

}
