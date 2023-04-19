package validators

import (
	notesv1 "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateGenerateInviteLinkRequest(req *notesv1.GenerateInviteLinkRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
	)
}

func ValidateUseInviteLinkRequest(req *notesv1.UseInviteLinkRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.InviteLinkCode, validation.Required),
	)
}

func ValidateRevokeInviteLinkRequest(req *notesv1.RevokeInviteLinkRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.InviteLinkCode, validation.Required),
	)
}

func ValidateGetInviteLinkRequest(req *notesv1.GetInviteLinkRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.InviteLinkCode, validation.Required),
	)
}
