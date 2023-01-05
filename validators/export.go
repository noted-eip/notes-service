package validators

import (
	notespb "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateExportNoteRequest(in *notespb.ExportNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.NoteId, validation.Required, is.UUID),
		validation.Field(&in.ExportFormat, validation.Required))
}
