package validators

import (
	notesv1 "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateCreateGroupRequest(req *notesv1.CreateGroupRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Name, validation.Required, validation.Length(1, 32)),
		validation.Field(&req.Description, validation.Required, validation.Length(1, 256)),
	)
}

func ValidateGetGroupRequest(req *notesv1.GetGroupRequest) error {
	return validation.ValidateStruct(req)
}

func ValidateUpdateGroupRequest(req *notesv1.UpdateGroupRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.Name, validation.Length(1, 32)),
		validation.Field(&req.Description, validation.Length(1, 256)),
	)
}

func ValidateDeleteGroupRequest(req *notesv1.DeleteGroupRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
	)
}

func ValidateListGroupsRequest(req *notesv1.ListGroupsRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.AccountId, validation.Required),
	)
}

func ValidateListActivitiesRequest(req *notesv1.ListActivitiesRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required, validation.Required),
	)
}

func ValidateGetActivityRequest(req *notesv1.GetActivityRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.ActivityId, validation.Required),
	)
}
