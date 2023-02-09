package validators

import (
	notesv1 "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateSendInviteRequest(req *notesv1.SendInviteRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.RecipientAccountId, validation.Required),
	)
}

func ValidateAcceptInviteRequest(req *notesv1.AcceptInviteRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.InviteId, validation.Required),
	)
}

func ValidateDenyInviteRequest(req *notesv1.DenyInviteRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.InviteId, validation.Required),
	)
}

func ValidateGetInviteRequest(req *notesv1.GetInviteRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.InviteId, validation.Required),
	)
}

func ValidateRevokeInviteRequest(req *notesv1.RevokeInviteRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.InviteId, validation.Required),
	)
}

func ValidateListInviteRequest(req *notesv1.ListInvitesRequest) error {
	return validation.ValidateStruct(req)
}
