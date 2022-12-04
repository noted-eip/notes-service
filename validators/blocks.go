package validators

import (
	notespb "notes-service/protorepo/noted/notes/v1"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateInsertBlockRequest(in *notespb.InsertBlockRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.NoteId, validation.When(strconv.Itoa(int(in.NoteId)) == "", validation.Required)),
		validation.Field(&in.Block.Data, validation.When(in.Block.Data == nil, validation.Required)),
		validation.Field(&in.Index, validation.When(in.Index < 1, validation.Required)),
		validation.Field(&in.Block.Type, validation.When(in.Block.Type < 1, validation.Required)),
	)
}

func ValidateUpdateBlockRequest(in *notespb.UpdateBlockRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.Id, validation.When(in.Id == "", validation.Required)),
	)
}

func ValidateDeleteBlockRequest(in *notespb.DeleteBlockRequest) error {
	return validation.ValidateStruct(in,
		validation.Field(&in.Id, validation.When(in.Id == "", validation.Required)),
	)
}
