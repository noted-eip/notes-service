package main

import (
	"context"
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"notes-service/validators"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (srv *groupsAPI) ListActivities(ctx context.Context, req *notesv1.ListActivitiesRequest) (*notesv1.ListActivitiesResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateListActivitiesRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check user is part of the group.
	_, err = srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	activities, err := srv.activities.ListActivitiesInternal(ctx,
		&models.ManyActivitiesFilter{GroupID: req.GroupId},
		&models.ListOptions{Limit: int32(req.Limit), Offset: int32(req.Offset)},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.ListActivitiesResponse{Activities: modelsGroupActivitiesToProtobufGroupActivities(activities)}, nil
}

func (srv *groupsAPI) GetActivity(ctx context.Context, req *notesv1.GetActivityRequest) (*notesv1.GetActivityResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateGetActivityRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check user is part of the group.
	_, err = srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	activity, err := srv.activities.GetActivityInternal(ctx, &models.OneActivityFilter{GroupID: req.GroupId, ActivityId: req.ActivityId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.GetActivityResponse{Activity: modelsGroupActivityToProtobufGroupActivity(activity)}, nil
}

func modelsGroupActivitiesToProtobufGroupActivities(activities []*models.Activity) []*notesv1.GroupActivity {
	protoActivities := make([]*notesv1.GroupActivity, len(activities))

	for i := range activities {
		protoActivities[i] = modelsGroupActivityToProtobufGroupActivity(activities[i])
	}

	return protoActivities
}

func modelsGroupActivityToProtobufGroupActivity(activity *models.Activity) *notesv1.GroupActivity {
	return &notesv1.GroupActivity{
		Id:        activity.ID,
		GroupId:   activity.GroupID,
		Type:      string(activity.Type),
		Event:     activity.Event,
		CreatedAt: timestamppb.New(activity.CreatedAt),
	}
}
