package validators

import (
	notespb "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateCreateNoteRequest(in *notespb.CreateNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.Note, validation.When(in.Note == nil, validation.Required)),
		validation.Field(&in.Note.AuthorId, validation.When(in.Note.AuthorId == "", validation.Required)),
	)
}

func ValidateGetNoteRequest(in *notespb.GetNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.Id, validation.When(in.Id == "", validation.Required)),
	)
}

func ValidateUpdateNoteRequest(in *notespb.UpdateNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.Id, validation.When(in.Id == "", validation.Required)),
		validation.Field(&in.Note, validation.When(in.Note == nil, validation.Required)),
		validation.Field(&in.Note.AuthorId, validation.When(in.Note.AuthorId == "", validation.Required)),
	)
}

func ValidateDeleteNoteRequest(in *notespb.DeleteNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.Id, validation.When(in.Id == "", validation.Required)),
	)
}

func ValidateListNoteRequest(in *notespb.ListNotesRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.AuthorId, validation.When(in.AuthorId == "", validation.Required)),
	)
}
