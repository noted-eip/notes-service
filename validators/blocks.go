package validators

import (
	notespb "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateInsertBlockRequest(in *notespb.InsertBlockRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.GroupId, validation.Required),
		validation.Field(&in.Block.Type, validation.Required),
	)
}

func ValidateUpdateBlockRequest(in *notespb.UpdateBlockRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.GroupId, validation.Required),
		validation.Field(&in.NoteId, validation.Required),
		validation.Field(&in.BlockId, validation.Required),
		validation.Field(&in.Block.Type, validation.Required),
	)
}

func ValidateDeleteBlockRequest(in *notespb.DeleteBlockRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.GroupId, validation.Required),
		validation.Field(&in.NoteId, validation.Required),
		validation.Field(&in.BlockId, validation.Required),
	)
}
