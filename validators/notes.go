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

	// Check update mask paths are set.
	for _, path := range req.UpdateMask.Paths {
		switch path {
		case "title":
			err = validation.Validate(&req.Note.Title, validation.Required, validation.Length(1, 64))
		case "blocks":
			err = validation.Validate(&req.Note.Blocks, validation.NotNil)
		default:
			return validation.NewError("update_mask", "update to "+path+" is forbidden")
		}
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
