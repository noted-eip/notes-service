package validators

import (
	notespb "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateInsertBlockRequest(in *notespb.InsertBlockRequest) error {
	err := validation.ValidateStruct(in,
		validation.Field(&in.NoteId, validation.Required),
		validation.Field(&in.Index, validation.Required),
		validation.Field(&in.Block, validation.NotNil),
	)
	if err != nil {
		return err
	}
	err = validation.ValidateStruct(in.Block,
		validation.Field(&in.Block.Data, validation.NotNil),
		validation.Field(&in.Block.Type, validation.Required),
	)
	if err != nil {
		return err
	}
	return nil
}

func ValidateUpdateBlockRequest(in *notespb.UpdateBlockRequest) error {
	return validation.ValidateStruct(in, validation.Field(&in.Id, validation.Required))
}

func ValidateDeleteBlockRequest(in *notespb.DeleteBlockRequest) error {
	return validation.ValidateStruct(in, validation.Field(&in.Id, validation.Required))
}
