package validators

import (
	notesv1 "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

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
