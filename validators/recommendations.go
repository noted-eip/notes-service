package validators

import (
	notesv1 "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateGenerateWidgetsRequest(req *notesv1.GenerateWidgetsRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.GroupId, validation.Required, validation.Length(1, 64)),
	)
}
