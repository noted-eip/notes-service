package main

import (
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"time"

	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestNotesSuite(t *testing.T) {
	tu := newTestUtilsOrDie(t)
	edouard := newTestAccount(t, tu)
	stranger := newTestAccount(t, tu)
	maxime := newTestAccount(t, tu)
	edouardGroup := newTestGroup(t, tu, edouard, maxime)

	t.Run("create-note", func(t *testing.T) {
		before := time.Now()
		res, err := tu.notes.CreateNote(edouard.Context, &notesv1.CreateNoteRequest{
			GroupId: edouardGroup.ID,
			Title:   "My New Note",
		})
		after := time.Now()
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "My New Note", res.Note.Title)
		require.Nil(t, res.Note.ModifiedAt)
		require.Nil(t, res.Note.AnalyzedAt)
		require.GreaterOrEqual(t, res.Note.CreatedAt.AsTime().Unix(), before.Unix())
		require.LessOrEqual(t, res.Note.CreatedAt.AsTime().Unix(), after.Unix())
	})

	t.Run("stranger-cannot-create-note", func(t *testing.T) {
		res, err := tu.notes.CreateNote(stranger.Context, &notesv1.CreateNoteRequest{
			GroupId: edouardGroup.ID,
			Title:   "My New Note",
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("member-can-create-note-with-blocks", func(t *testing.T) {
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

	t.Run("member-can-create-note-with-all-block-types", func(t *testing.T) {
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

	t.Run("member-can-create-note-with-invalid-blocks", func(t *testing.T) {
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

	t.Run("owner-can-update-note-title", func(t *testing.T) {
		before := time.Now()
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
		after := time.Now()

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "Brand New Title", res.Note.Title)
		require.GreaterOrEqual(t, res.Note.ModifiedAt.AsTime().Unix(), before.Unix())
		require.LessOrEqual(t, res.Note.ModifiedAt.AsTime().Unix(), after.Unix())
	})

	t.Run("member-cannot-update-note-title", func(t *testing.T) {
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

	t.Run("stranger-cannot-update-note-title", func(t *testing.T) {
		res, err := tu.notes.UpdateNote(stranger.Context, &notesv1.UpdateNoteRequest{
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

	t.Run("member-cannot-update-note-title", func(t *testing.T) {
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

	t.Run("owner-can-update-note-blocks", func(t *testing.T) {
		res, err := tu.notes.UpdateNote(edouard.Context, &notesv1.UpdateNoteRequest{
			NoteId:  edouardNote.ID,
			GroupId: edouardGroup.ID,
			Note: &notesv1.Note{
				Blocks: []*notesv1.Block{
					{
						Type: notesv1.Block_TYPE_HEADING_1,
						Data: &notesv1.Block_Heading{
							Heading: "Heading",
						},
					},
					{
						Type: notesv1.Block_TYPE_CODE,
						Data: &notesv1.Block_Code_{
							Code: &notesv1.Block_Code{
								Lang:    "go",
								Snippet: "package main",
							},
						},
					},
					{
						Type: notesv1.Block_TYPE_IMAGE,
						Data: &notesv1.Block_Image_{
							Image: &notesv1.Block_Image{
								Caption: "Image",
								Url:     "https://example.com/image.png",
							},
						},
					},
				},
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"blocks"},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "Heading", res.Note.Blocks[0].GetHeading())
		require.Equal(t, "go", res.Note.Blocks[1].GetCode().Lang)
		require.Equal(t, "package main", res.Note.Blocks[1].GetCode().Snippet)
		require.Equal(t, "Image", res.Note.Blocks[2].GetImage().Caption)
		require.Equal(t, "https://example.com/image.png", res.Note.Blocks[2].GetImage().Url)
		require.Equal(t, "Brand New Title", res.Note.Title)
		require.Len(t, res.Note.Blocks, 3)
	})

	t.Run("owner-can-update-note-blocks-with-empty-blocks-array", func(t *testing.T) {
		res, err := tu.notes.UpdateNote(edouard.Context, &notesv1.UpdateNoteRequest{
			NoteId:  edouardNote.ID,
			GroupId: edouardGroup.ID,
			Note: &notesv1.Note{
				Blocks: []*notesv1.Block{},
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"blocks"},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Note.Blocks, 0)
	})

	t.Run("owner-cannot-update-note-with-invalid-field-mask", func(t *testing.T) {
		res, err := tu.notes.UpdateNote(edouard.Context, &notesv1.UpdateNoteRequest{
			NoteId:  edouardNote.ID,
			GroupId: edouardGroup.ID,
			Note: &notesv1.Note{
				Title: "Brand New Title",
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"blocks"},
			},
		})
		requireErrorHasGRPCCode(t, codes.InvalidArgument, err)
		require.Nil(t, res)
	})

	t.Run("member-cannot-delete-note", func(t *testing.T) {
		res, err := tu.notes.DeleteNote(edouard.Context, &notesv1.DeleteNoteRequest{
			GroupId: maximeNote.Group.ID,
			NoteId:  maximeNote.ID,
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("stranger-cannot-delete-note", func(t *testing.T) {
		res, err := tu.notes.DeleteNote(stranger.Context, &notesv1.DeleteNoteRequest{
			GroupId: maximeNote.Group.ID,
			NoteId:  maximeNote.ID,
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	t.Run("owner-can-delete-note", func(t *testing.T) {
		res, err := tu.notes.DeleteNote(maxime.Context, &notesv1.DeleteNoteRequest{
			GroupId: maximeNote.Group.ID,
			NoteId:  maximeNote.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	newTestNote(t, tu, edouardGroup, maxime, nil)

	t.Run("member-can-list-notes-by-group", func(t *testing.T) {
		res, err := tu.notes.ListNotes(maxime.Context, &notesv1.ListNotesRequest{
			GroupId: edouardGroup.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Notes, 6)

		// Make sure only the note's metadata is returned.
		require.Nil(t, res.Notes[0].Blocks)
		require.Nil(t, res.Notes[1].Blocks)
		require.Nil(t, res.Notes[2].Blocks)
		require.Nil(t, res.Notes[3].Blocks)
		require.Nil(t, res.Notes[4].Blocks)
	})

	t.Run("stranger-cannot-list-notes", func(t *testing.T) {
		res, err := tu.notes.ListNotes(stranger.Context, &notesv1.ListNotesRequest{
			GroupId: edouardGroup.ID,
		})
		requireErrorHasGRPCCode(t, codes.NotFound, err)
		require.Nil(t, res)
	})

	maximeGroup := newTestGroup(t, tu, maxime, edouard)
	newTestNote(t, tu, maximeGroup, edouard, nil)

	t.Run("user-can-list-their-notes-across-groups", func(t *testing.T) {
		res, err := tu.notes.ListNotes(edouard.Context, &notesv1.ListNotesRequest{
			AuthorAccountId: edouard.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Notes, 6)
	})

	t.Run("user-can-list-their-notes-in-group", func(t *testing.T) {
		res, err := tu.notes.ListNotes(edouard.Context, &notesv1.ListNotesRequest{
			AuthorAccountId: edouard.ID,
			GroupId:         edouardGroup.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Notes, 5)
	})
}
