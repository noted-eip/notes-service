package validators

import (
	"errors"
	notesv1 "notes-service/protorepo/noted/notes/v1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateListActivitiesRequest(req *notesv1.ListActivitiesRequest) error {
	errorGroupId := validation.ValidateStruct(req, validation.Field(&req.GroupId, validation.Required, validation.Required))
	errorAccountId := validation.ValidateStruct(req, validation.Field(&req.AccountId, validation.Required, validation.Required))

	if errorGroupId != nil && errorAccountId != nil {
		return errors.New("A GroupId or AccountId should be provided")
	}
	return nil
}

func ValidateGetActivityRequest(req *notesv1.GetActivityRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.GroupId, validation.Required),
		validation.Field(&req.ActivityId, validation.Required),
	)
}
