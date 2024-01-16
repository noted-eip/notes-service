package main

import (
	"context"
	"notes-service/models"
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
	maximeGroup := newTestGroup(t, tu, maxime, edouard)

	testUser := newTestAccount(t, tu)
	testGroup := newTestGroup(t, tu, testUser)
	note := newTestNote(t, tu, testGroup, testUser, []*notesv1.Block{
		{
			Type: notesv1.Block_TYPE_HEADING_1,
			Data: &notesv1.Block_Heading{
				Heading: "Ada Lovelace",
			},
		},
		{ // TODO: Put placeholder texts in separate file
			Type: notesv1.Block_TYPE_PARAGRAPH,
			Data: &notesv1.Block_Paragraph{
				Paragraph: "Ada Lovelace, de son nom complet Augusta Ada King, comtesse de Lovelace, née Ada Byron le 10 décembre 1815 à Londres et morte le 27 novembre 1852 à Marylebone dans la même ville, est une pionnière de la science informatique. Elle est principalement connue pour avoir réalisé le premier véritable programme informatique, lors de son travail sur un ancêtre de l'ordinateur : la machine analytique de Charles Babbage. Dans ses notes, on trouve en effet le premier programme publié, destiné à être exécuté par une machine, ce qui fait d'Ada Lovelace la première personne à avoir programmé au monde. Elle a également entrevu et décrit certaines possibilités offertes par les calculateurs universels, allant bien au-delà du calcul numérique et de ce qu'imaginaient Babbage et ses contemporains. ",
			},
		},
		{
			Type: notesv1.Block_TYPE_PARAGRAPH,
			Data: &notesv1.Block_Paragraph{
				Paragraph: "Ada était la seule fille légitime du poète George Gordon Byron et de son épouse Annabella Milbanke, une femme intelligente et cultivée, cousine de Caroline Lamb, dont la liaison avec Byron fut à l'origine d'un scandale. Le premier prénom d'Ada, Augusta, aurait été choisi en hommage à Augusta Leigh, la demi-sœur de Byron, avec qui ce dernier aurait eu des relations incestueusesSwade 1. Le prénom Ada aurait été choisi par Byron lui-mêmeStein 1, car il était « court, ancien et vocalique »Wolfram 1. C'est Augusta qui encouragea Byron à se marier pour éviter un scandale, et il épousa Annabella à contrecœur[réf. souhaitée], en janvier 1815. Ada naît en décembre de cette même année. À la suite de quatre tentatives de viol en état d'ivresse de la part de ByronSwade 1, Annabella quitte Byron le 16 janvier 1816, gardant Ada avec elle. Le 21 avril, Byron signe l'acte de séparation, puis quitte le Royaume-Uni pour toujours. Il ne les revit jamais.",
			},
		},
	})

	//
	//
	// Notes/Blocks CRUD tests
	//
	//

	t.Run("create-note", func(t *testing.T) {
		before := time.Now()
		res, err := tu.notes.CreateNote(edouard.Context, &notesv1.CreateNoteRequest{
			GroupId: edouardGroup.ID,
			Title:   "My New Note",
			Lang:    "en",
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
			Lang:    "en",
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
			Lang: "en",
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
			Lang: "en",
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
			Lang: "en",
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

	t.Run("member-no-edit-rights-cannot-update-note-title", func(t *testing.T) {
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
		requireErrorHasGRPCCode(t, codes.PermissionDenied, err)
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

	t.Run("owner-can-update-note-styles", func(t *testing.T) {
		note, err := tu.notes.GetNote(edouard.Context, &notesv1.GetNoteRequest{
			NoteId:  edouardNote.ID,
			GroupId: edouardGroup.ID,
		})
		require.NoError(t, err)

		_, err = tu.notes.UpdateBlock(edouard.Context, &notesv1.UpdateBlockRequest{
			GroupId: edouardGroup.ID,
			NoteId:  edouardNote.ID,
			BlockId: note.Note.Blocks[0].Id,
			Block: &notesv1.Block{
				Styles: []*notesv1.Block_TextStyle{
					{
						Style: notesv1.Block_TextStyle_STYLE_BOLD,
						Pos: &notesv1.Block_TextStyle_Position{
							Start:  12,
							Length: 10,
						},
						Color: &notesv1.Block_TextStyle_Color{
							R: 12,
							G: 120,
							B: 12,
						},
					},
				},
			},
		},
		)

		require.NoError(t, err)

		note, err = tu.notes.GetNote(edouard.Context, &notesv1.GetNoteRequest{
			NoteId:  edouardNote.ID,
			GroupId: edouardGroup.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, note)
		style := note.Note.Blocks[0].Styles[0]
		require.Equal(t, style.Style, notesv1.Block_TextStyle_STYLE_BOLD)
		require.Equal(t, style.Color.R, int32(12))
		require.Equal(t, style.Color.G, int32(120))
		require.Equal(t, style.Color.B, int32(12))
		require.Equal(t, style.Pos.Start, int64(12))
		require.Equal(t, style.Pos.Length, int64(10))
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

	t.Run("owner-cannot-update-note-with-invalid-field-mask-path", func(t *testing.T) {
		res, err := tu.notes.UpdateNote(edouard.Context, &notesv1.UpdateNoteRequest{
			NoteId:  edouardNote.ID,
			GroupId: edouardGroup.ID,
			Note: &notesv1.Note{
				Title: "Brand New Title",
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"invalid-field"},
			},
		})
		requireErrorHasGRPCCode(t, codes.InvalidArgument, err)
		require.Nil(t, res)
	})

	t.Run("owner-cannot-update-note-with-several-invalid-field-mask-path", func(t *testing.T) {
		res, err := tu.notes.UpdateNote(edouard.Context, &notesv1.UpdateNoteRequest{
			NoteId:  edouardNote.ID,
			GroupId: edouardGroup.ID,
			Note: &notesv1.Note{
				Title: "Brand New Title",
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"invalid-field",
					"very-invalid-field",
				},
			},
		})
		requireErrorHasGRPCCode(t, codes.InvalidArgument, err)
		require.Nil(t, res)
	})

	t.Run("owner-can-update-note-with-one-valid-field-mask-path", func(t *testing.T) {
		res, err := tu.notes.UpdateNote(edouard.Context, &notesv1.UpdateNoteRequest{
			NoteId:  edouardNote.ID,
			GroupId: edouardGroup.ID,
			Note: &notesv1.Note{
				Title: "Brand New Title",
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"invalid-field",
					"title",
					"very-invalid-field",
				},
			},
		})
		require.Nil(t, err)
		require.NotNil(t, res)
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

	_ = newTestNote(t, tu, edouardGroup, edouard, []*notesv1.Block{})
	_ = newTestNote(t, tu, edouardGroup, edouard, []*notesv1.Block{})
	_ = newTestNote(t, tu, maximeGroup, edouard, []*notesv1.Block{})

	t.Run("delete-notes-account", func(t *testing.T) {
		// DeleteNotes is a repository function, no auth
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

	// Clean-up maximeGroup made notes (Because of background service bug)
	err := tu.notesRepository.DeleteNotes(context.TODO(), &models.ManyNotesFilter{
		AuthorAccountID: maxime.ID,
	})
	require.NoError(t, err)

	//
	//
	// Quiz tests
	//
	//

	var quizIDContainer string
	// NOTE: This test takes at least 5 seconds
	t.Run("generate-quiz-success", func(t *testing.T) {
		res, err := tu.notes.GenerateQuiz(testUser.Context, &notesv1.GenerateQuizRequest{
			GroupId: note.Group.ID,
			NoteId:  note.ID,
		})
		require.NoError(t, err)

		require.NotZero(t, len(res.Quiz.Questions))
		for _, question := range res.Quiz.Questions {
			require.NotZero(t, len(question.Question))
			require.NotZero(t, len(question.Answers))
			require.NotZero(t, len(question.Solutions))
		}

		quizIDContainer = res.Quiz.Id
	})

	t.Run("quiz-stored-after-generated-quiz-and-author-can-list-quizs", func(t *testing.T) {
		res, err := tu.notes.ListQuizs(note.Author.Context, &notesv1.ListQuizsRequest{
			GroupId: note.Group.ID,
			NoteId:  note.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		require.Equal(t, 1, len(res.Quizs))
	})

	t.Run("stranger-cannot-list-quiz", func(t *testing.T) {
		res, err := tu.notes.ListQuizs(stranger.Context, &notesv1.ListQuizsRequest{
			GroupId: note.Group.ID,
			NoteId:  note.ID,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("stranger-cannot-delete-quiz", func(t *testing.T) {
		err := tu.notesRepository.DeleteQuiz(stranger.Context, &models.OneNoteFilter{
			GroupID: note.Group.ID,
			NoteID:  note.ID,
		}, quizIDContainer, stranger.ID)
		require.Error(t, err)
	})

	t.Run("author-can-delete-quiz", func(t *testing.T) {
		err := tu.notesRepository.DeleteQuiz(note.Author.Context, &models.OneNoteFilter{
			GroupID: note.Group.ID,
			NoteID:  note.ID,
		}, quizIDContainer, note.Author.ID)
		require.NoError(t, err)

		res, err := tu.notes.ListQuizs(note.Author.Context, &notesv1.ListQuizsRequest{
			GroupId: note.Group.ID,
			NoteId:  note.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		require.Zero(t, len(res.Quizs))
	})

	//
	//
	// Edit permissions tests
	//
	//

	// Maxime joins test group
	invite := testUser.SendInvite(t, tu, maxime, testGroup)
	maxime.AcceptInvite(t, tu, invite)

	// Edouard joins test group
	invite = testUser.SendInvite(t, tu, edouard, testGroup)
	edouard.AcceptInvite(t, tu, invite)

	t.Run("member-can-list-quizs", func(t *testing.T) {
		res, err := tu.notes.ListQuizs(edouard.Context, &notesv1.ListQuizsRequest{
			GroupId: note.Group.ID,
			NoteId:  note.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("non-author-cannot-grant-permission", func(t *testing.T) {
		res, err := tu.notes.ChangeNoteEditPermission(maxime.Context, &notesv1.ChangeNoteEditPermissionRequest{
			GroupId:            note.Group.ID,
			NoteId:             note.ID,
			RecipientAccountId: edouard.ID,
			Type:               notesv1.ChangeNoteEditPermissionRequest_ACTION_GRANT,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("author-can-grant-edit-permissions", func(t *testing.T) {
		res, err := tu.notes.ChangeNoteEditPermission(note.Author.Context, &notesv1.ChangeNoteEditPermissionRequest{
			GroupId:            note.Group.ID,
			NoteId:             note.ID,
			RecipientAccountId: maxime.ID,
			Type:               notesv1.ChangeNoteEditPermissionRequest_ACTION_GRANT,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("non-author-cannot-grant-permission-even-with-edit-rights", func(t *testing.T) {
		res, err := tu.notes.ChangeNoteEditPermission(maxime.Context, &notesv1.ChangeNoteEditPermissionRequest{
			GroupId:            note.Group.ID,
			NoteId:             note.ID,
			RecipientAccountId: edouard.ID,
			Type:               notesv1.ChangeNoteEditPermissionRequest_ACTION_GRANT,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("author-cannot-grant-edit-permissions-to-stranger", func(t *testing.T) {
		res, err := tu.notes.ChangeNoteEditPermission(note.Author.Context, &notesv1.ChangeNoteEditPermissionRequest{
			GroupId:            note.Group.ID,
			NoteId:             note.ID,
			RecipientAccountId: stranger.ID,
			Type:               notesv1.ChangeNoteEditPermissionRequest_ACTION_GRANT,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("user-can-update-note-with-edit-permission", func(t *testing.T) {
		newTitle := "Hacked by Maximator"

		res, err := tu.notes.UpdateNote(
			maxime.Context,
			&notesv1.UpdateNoteRequest{
				GroupId: note.Group.ID,
				NoteId:  note.ID,
				Note: &notesv1.Note{
					Title: newTitle,
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"title"}},
			})
		require.NoError(t, err)
		require.Equal(t, res.Note.Title, newTitle)
	})

	t.Run("remove-edit-permissions-when-leaving-a-group", func(t *testing.T) {
		res, err := tu.notesRepository.GetNote(testUser.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)
		require.NoError(t, err)
		OldNumberOfEditors := len(res.AccountsWithEditPermissions)

		_, err = tu.groups.RemoveMember(testUser.Context,
			&notesv1.RemoveMemberRequest{
				GroupId:   testGroup.ID,
				AccountId: maxime.ID,
			})
		require.NoError(t, err)

		res, err = tu.notesRepository.GetNote(testUser.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)

		require.NoError(t, err)
		require.Equal(t, len(res.AccountsWithEditPermissions), OldNumberOfEditors-1)
	})

	// Granting permissions to Edouard
	tu.notes.ChangeNoteEditPermission(note.Author.Context, &notesv1.ChangeNoteEditPermissionRequest{
		GroupId:            note.Group.ID,
		NoteId:             note.ID,
		RecipientAccountId: edouard.ID,
		Type:               notesv1.ChangeNoteEditPermissionRequest_ACTION_GRANT,
	})

	t.Run("non-author-cannot-remove-others-permission", func(t *testing.T) {
		res, err := tu.notes.ChangeNoteEditPermission(edouard.Context, &notesv1.ChangeNoteEditPermissionRequest{
			GroupId:            note.Group.ID,
			NoteId:             note.ID,
			RecipientAccountId: testUser.ID,
			Type:               notesv1.ChangeNoteEditPermissionRequest_ACTION_REMOVE,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("author-can-remove-edit-permissions", func(t *testing.T) {
		res, err := tu.notesRepository.GetNote(testUser.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)
		require.NoError(t, err)
		OldNumberOfEditors := len(res.AccountsWithEditPermissions)

		r, err := tu.notes.ChangeNoteEditPermission(note.Author.Context, &notesv1.ChangeNoteEditPermissionRequest{
			GroupId:            note.Group.ID,
			NoteId:             note.ID,
			RecipientAccountId: edouard.ID,
			Type:               notesv1.ChangeNoteEditPermissionRequest_ACTION_REMOVE,
		})
		require.NoError(t, err)
		require.NotNil(t, r)

		res, err = tu.notesRepository.GetNote(testUser.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)

		require.NoError(t, err)
		require.Equal(t, len(res.AccountsWithEditPermissions), OldNumberOfEditors-1)

	})

	t.Run("user-cannot-update-note-after-removing-permissions", func(t *testing.T) {
		newTitle := "Hacked by Edouardino"

		res, err := tu.notes.UpdateNote(
			edouard.Context,
			&notesv1.UpdateNoteRequest{
				GroupId: note.Group.ID,
				NoteId:  note.ID,
				Note: &notesv1.Note{
					Title: newTitle,
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"title"}},
			})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("author-cannot-remove-edit-permissions-to-stranger", func(t *testing.T) {
		res, err := tu.notes.ChangeNoteEditPermission(note.Author.Context, &notesv1.ChangeNoteEditPermissionRequest{
			GroupId:            note.Group.ID,
			NoteId:             note.ID,
			RecipientAccountId: stranger.ID,
			Type:               notesv1.ChangeNoteEditPermissionRequest_ACTION_REMOVE,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	// Granting permissions to Edouard
	tu.notes.ChangeNoteEditPermission(note.Author.Context, &notesv1.ChangeNoteEditPermissionRequest{
		GroupId:            note.Group.ID,
		NoteId:             note.ID,
		RecipientAccountId: edouard.ID,
		Type:               notesv1.ChangeNoteEditPermissionRequest_ACTION_GRANT,
	})

	t.Run("user-can-remove-his-own-edit-permissions", func(t *testing.T) {
		res, err := tu.notesRepository.GetNote(testUser.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)
		require.NoError(t, err)
		OldNumberOfEditors := len(res.AccountsWithEditPermissions)

		r, err := tu.notes.ChangeNoteEditPermission(edouard.Context, &notesv1.ChangeNoteEditPermissionRequest{
			GroupId:            note.Group.ID,
			NoteId:             note.ID,
			RecipientAccountId: edouard.ID,
			Type:               notesv1.ChangeNoteEditPermissionRequest_ACTION_REMOVE,
		})

		require.NoError(t, err)
		require.NotNil(t, r)

		res, err = tu.notesRepository.GetNote(testUser.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)

		require.NoError(t, err)
		require.Equal(t, len(res.AccountsWithEditPermissions), OldNumberOfEditors-1)

	})

	t.Run("author-can-comment-on-his-note", func(t *testing.T) {
		res, err := tu.notesRepository.GetNote(note.Author.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)

		require.NoError(t, err)
		require.NotZero(t, len(*res.Blocks))

		blockID := (*res.Blocks)[0].ID
		commentContent := "very nice comment"

		r, err := tu.notes.CreateBlockComment(note.Author.Context, &notesv1.CreateBlockCommentRequest{
			GroupId: testGroup.ID,
			NoteId:  note.ID,
			BlockId: blockID,
			Comment: &notesv1.Block_Comment{
				AuthorId: note.Author.ID,
				Content:  commentContent,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, r)
		require.Equal(t, r.Comment.Content, commentContent)
	})

	t.Run("stranger-cannot-comment-on-author-note", func(t *testing.T) {
		res, err := tu.notesRepository.GetNote(note.Author.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)

		require.NoError(t, err)
		require.NotZero(t, len(*res.Blocks))

		blockID := (*res.Blocks)[0].ID
		commentContent := "very nice comment"

		r, err := tu.notes.CreateBlockComment(stranger.Context, &notesv1.CreateBlockCommentRequest{
			GroupId: testGroup.ID,
			NoteId:  note.ID,
			BlockId: blockID,
			Comment: &notesv1.Block_Comment{
				AuthorId: note.Author.ID,
				Content:  commentContent,
			},
		})

		require.Error(t, err)
		require.Nil(t, r)
	})

	// t.Run("user-cannot-comment-on-note-on-which-he-has-not-access-to", func(t *testing.T) {
	// 	res, err := tu.notesRepository.GetNote(note.Author.Context, &models.OneNoteFilter{
	// 		GroupID: testGroup.ID,
	// 		NoteID:  note.ID,
	// 	}, testUser.ID)

	// 	require.NoError(t, err)
	// 	require.NotZero(t, len(*res.Blocks))

	// 	blockID := (*res.Blocks)[0].ID
	// 	commentContent := "very nice comment"

	// 	r, err := tu.notes.CreateBlockComment(edouard.Context, &notesv1.CreateBlockCommentRequest{
	// 		GroupId: testGroup.ID,
	// 		NoteId:  note.ID,
	// 		BlockId: blockID,
	// 		Comment: &notesv1.Block_Comment{
	// 			AuthorId: note.Author.ID,
	// 			Content:  commentContent,
	// 		},
	// 	})

	// 	require.Error(t, err)
	// 	require.Nil(t, r)
	// })

	// t.Run("user-cannot-list-comments-on-note-on-which-he-has-not-access-to", func(t *testing.T) {
	// 	res, err := tu.notesRepository.GetNote(note.Author.Context, &models.OneNoteFilter{
	// 		GroupID: testGroup.ID,
	// 		NoteID:  note.ID,
	// 	}, testUser.ID)

	// 	require.NoError(t, err)
	// 	require.NotZero(t, len(*res.Blocks))

	// 	blockID := (*res.Blocks)[0].ID

	// 	r, err := tu.notes.ListBlockComments(edouard.Context, &notesv1.ListBlockCommentsRequest{
	// 		GroupId: testGroup.ID,
	// 		NoteId:  note.ID,
	// 		BlockId: blockID,
	// 	})

	// 	require.Error(t, err)
	// 	require.Nil(t, r)
	// })

	// Change back permissions for edouard on testUser's note
	r, err := tu.notes.ChangeNoteEditPermission(testUser.Context, &notesv1.ChangeNoteEditPermissionRequest{
		GroupId:            note.Group.ID,
		NoteId:             note.ID,
		RecipientAccountId: edouard.ID,
		Type:               notesv1.ChangeNoteEditPermissionRequest_ACTION_GRANT,
	})
	require.NoError(t, err)
	require.NotNil(t, r)

	t.Run("user-can-comment-on-note-on-which-he-has-access-to", func(t *testing.T) {
		res, err := tu.notesRepository.GetNote(note.Author.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)

		require.NoError(t, err)
		require.NotZero(t, len(*res.Blocks))

		blockID := (*res.Blocks)[0].ID
		commentContent := "very nice comment"

		r, err := tu.notes.CreateBlockComment(edouard.Context, &notesv1.CreateBlockCommentRequest{
			GroupId: testGroup.ID,
			NoteId:  note.ID,
			BlockId: blockID,
			Comment: &notesv1.Block_Comment{
				AuthorId: note.Author.ID,
				Content:  commentContent,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, r)
		require.Equal(t, commentContent, r.Comment.Content)
	})

	t.Run("user-can-list-comments-on-note-on-which-he-has-access-to", func(t *testing.T) {
		res, err := tu.notesRepository.GetNote(note.Author.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)

		require.NoError(t, err)
		require.NotZero(t, len(*res.Blocks))

		blockID := (*res.Blocks)[0].ID

		r, err := tu.notes.ListBlockComments(edouard.Context, &notesv1.ListBlockCommentsRequest{
			GroupId: testGroup.ID,
			NoteId:  note.ID,
			BlockId: blockID,
		})

		require.NoError(t, err)
		require.NotNil(t, r)
		require.Equal(t, len(r.Comments), 2)
	})

	t.Run("non-author-cannot-delete-an-other-user-comment", func(t *testing.T) {
		res, err := tu.notesRepository.GetNote(note.Author.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)
		require.NoError(t, err)
		require.NotZero(t, len(*res.Blocks))

		blockID := (*res.Blocks)[0].ID
		commentID := ""
		for _, cmt := range *(*res.Blocks)[0].Thread {
			if cmt.AuthorAccountID == note.Author.ID {
				commentID = cmt.ID
				break
			}
		}
		require.NotEmpty(t, commentID)

		r, err := tu.notes.DeleteBlockComment(maxime.Context, &notesv1.DeleteBlockCommentRequest{
			GroupId:   testGroup.ID,
			NoteId:    note.ID,
			BlockId:   blockID,
			CommentId: commentID,
		})

		require.Error(t, err)
		require.Nil(t, r)
	})

	t.Run("is-possible-to-delete-own-comment", func(t *testing.T) {
		res, err := tu.notesRepository.GetNote(note.Author.Context, &models.OneNoteFilter{
			GroupID: testGroup.ID,
			NoteID:  note.ID,
		}, testUser.ID)

		require.NoError(t, err)
		require.NotZero(t, len(*res.Blocks))

		blockID := (*res.Blocks)[0].ID
		commentID := ""
		for _, cmt := range *(*res.Blocks)[0].Thread {
			if cmt.AuthorAccountID == note.Author.ID {
				commentID = cmt.ID
				break
			}
		}
		require.NotEmpty(t, commentID)

		r, err := tu.notes.DeleteBlockComment(note.Author.Context, &notesv1.DeleteBlockCommentRequest{
			GroupId:   testGroup.ID,
			NoteId:    note.ID,
			BlockId:   blockID,
			CommentId: commentID,
		})

		require.NoError(t, err)
		require.NotNil(t, r)
	})
}
