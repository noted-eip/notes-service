package validators

import (
	notespb "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateInsertBlockRequest(req *notespb.InsertBlockRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.Block.Type, validation.Required),
	)
}

func ValidateUpdateBlockRequest(req *notespb.UpdateBlockRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.BlockId, validation.Required),
		validation.Field(&req.Block.Type, validation.Required),
	)
}

func ValidateDeleteBlockRequest(req *notespb.DeleteBlockRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.NoteId, validation.Required),
		validation.Field(&req.BlockId, validation.Required),
	)
}
