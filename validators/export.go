package validators

import (
	notespb "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateExportNoteRequest(in *notespb.ExportNoteRequest) error {
	err := validation.ValidateStruct(in,
		validation.Field(&in.NoteId, validation.Required),
		validation.Field(&in.ExportFormat, validation.Required),
	)
	return err
}
