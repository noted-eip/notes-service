package validators

import (
	notespb "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateCreateNoteRequest(req *notespb.CreateNoteRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.Title, validation.Required, validation.Length(1, 64)),
	)
}

func ValidateGetNoteRequest(req *notespb.GetNoteRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
	)
}

func ValidateUpdateNoteRequest(req *notespb.UpdateNoteRequest) error {
	err := validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.UpdateMask, validation.Required),
	)
	if err != nil {
		return err
	}

	// Check update mask is not empty.
	err = validation.Validate(&req.UpdateMask.Paths, validation.Required)
	if err != nil {
		return err
	}

	cptValideFieldMask := 0

	for i := 0; i < len(req.UpdateMask.Paths); i++ {
		path := req.UpdateMask.Paths[i]
		if path == "title" {
			cptValideFieldMask++
			err = validation.Validate(&req.Note.Title, validation.Required, validation.Length(1, 64))
		} else if path == "blocks" {
			cptValideFieldMask++
			err = validation.Validate(&req.Note.Blocks, validation.NotNil)
		} else {
			// if update mask is not allowed, we remove it from the list
			req.UpdateMask.Paths = append(req.UpdateMask.Paths[:i], req.UpdateMask.Paths[i+1:]...)
			i--
			continue
		}
	}

	if cptValideFieldMask < 1 {
		return validation.NewError("update_mask", "At least one valide update mask should be valid")
	}

	return err
}

func ValidateDeleteNoteRequest(req *notespb.DeleteNoteRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
	)
}

func ValidateListNoteRequest(req *notespb.ListNotesRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.AuthorAccountId, validation.When(req.GroupId == "", validation.Required)),
		validation.Field(&req.GroupId, validation.When(req.AuthorAccountId == "", validation.Required)),
	)
}

func ValidateExportNoteRequest(req *notespb.ExportNoteRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.ExportFormat, validation.Required),
	)
}

func ValidateGenerateQuizzRequest(req *notespb.GenerateQuizRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
	)
}

func ValidateGenerateSummaryRequest(req *notespb.GenerateSummaryRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
	)
}

func ValidateCreateBlockCommentRequest(req *notespb.CreateBlockCommentRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.BlockId, validation.Required),
		validation.Field(&req.Comment, validation.Required),
	)
}

func ValidateListBlockCommentRequest(req *notespb.ListBlockCommentsRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.BlockId, validation.Required),
	)
}

func ValidateDeleteBlockCommentRequest(req *notespb.DeleteBlockCommentRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.BlockId, validation.Required),
		validation.Field(&req.CommentId, validation.Required),
	)
}

func ValidateChangeEditPermissionsRequest(req *notespb.ChangeNoteEditPermissionRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.Type, validation.Required),
		validation.Field(&req.RecipientAccountId, validation.Required),
	)
}
