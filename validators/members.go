package validators

import (
	notesv1 "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateGetMemberRequest(req *notesv1.GetMemberRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.AccountId, validation.Required),
	)
}

func ValidateUpdateMemberRequest(req *notesv1.UpdateMemberRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.AccountId, validation.Required),
		validation.Field(&req.Member, validation.Required),
		validation.Field(&req.UpdateMask, validation.Required),
	)
}

func ValidateRemoveMemberRequest(req *notesv1.RemoveMemberRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.AccountId, validation.Required),
	)
}
