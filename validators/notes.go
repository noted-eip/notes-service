package validators

import (
	notespb "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateCreateNoteRequest(in *notespb.CreateNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.GroupId, validation.Required),
		validation.Field(&in.Title, validation.Required, validation.Length(1, 64)),
	)
}

func ValidateGetNoteRequest(in *notespb.GetNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.GroupId, validation.Required),
		validation.Field(&in.NoteId, validation.Required),
	)
}

func ValidateUpdateNoteRequest(in *notespb.UpdateNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.GroupId, validation.Required),
		validation.Field(&in.NoteId, validation.Required),
		validation.Field(&in.Note.Title, validation.Length(1, 64)),
	)
}

func ValidateDeleteNoteRequest(in *notespb.DeleteNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.GroupId, validation.Required),
		validation.Field(&in.NoteId, validation.Required),
	)
}

func ValidateListNoteRequest(in *notespb.ListNotesRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.AuthorId, validation.Required),
	)
}

func ValidateExportNoteRequest(in *notespb.ExportNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.GroupId, validation.Required),
		validation.Field(&in.NoteId, validation.Required),
		validation.Field(&in.ExportFormat, validation.Required),
	)
}
