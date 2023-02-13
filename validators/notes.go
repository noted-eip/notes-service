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
	)
	if err != nil {
		return err
	}
	return validation.Validate(req.Note.Title, validation.Length(0, 64))
}

func ValidateDeleteNoteRequest(req *notespb.DeleteNoteRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
	)
}

func ValidateListNoteRequest(req *notespb.ListNotesRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.AuthorAccountId, validation.Required),
	)
}

func ValidateExportNoteRequest(req *notespb.ExportNoteRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.ExportFormat, validation.Required),
	)
}
