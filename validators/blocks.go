package validators

import (
	notespb "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateInsertBlockRequest(req *notespb.InsertBlockRequest) error {
	err := validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
	)
	if err != nil {
		return err
	}
	return validation.Validate(req.Block, validation.Required)
}

func ValidateUpdateBlockRequest(req *notespb.UpdateBlockRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.BlockId, validation.Required),
	)
}

func ValidateUpdateBlockIndexRequest(req *notespb.UpdateBlockIndexRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.BlockId, validation.Required),
		validation.Field(&req.Index, validation.Required),
	)
}

func ValidateDeleteBlockRequest(req *notespb.DeleteBlockRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.BlockId, validation.Required),
	)
}
