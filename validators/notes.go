package validators

import (
	notespb "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateCreateNoteRequest(in *notespb.CreateNoteRequest) error {
	err := validation.ValidateStruct(in, validation.Field(&in.Note, validation.NotNil))
	if err != nil {
		return err
	}
	err = validation.ValidateStruct(in.Note, validation.Field(&in.Note.AuthorId, validation.Required))
	if err != nil {
		return err
	}
	return nil
}

func ValidateGetNoteRequest(in *notespb.GetNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.Id, validation.Required, is.UUID),
	)
}

func ValidateUpdateNoteRequest(in *notespb.UpdateNoteRequest) error {
	err := validation.ValidateStruct(in,
		validation.Field(&in.Note, validation.NotNil),
		validation.Field(&in.Id, validation.Required, is.UUID),
	)
	if err != nil {
		return err
	}
	err = validation.ValidateStruct(in.Note, validation.Field(&in.Note.AuthorId, validation.Required))
	if err != nil {
		return err
	}
	return nil
}

func ValidateDeleteNoteRequest(in *notespb.DeleteNoteRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.Id, validation.Required, is.UUID),
	)
}

func ValidateListNoteRequest(in *notespb.ListNotesRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.AuthorId, validation.Required),
	)
}
