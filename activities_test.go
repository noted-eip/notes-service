package main

import (
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"time"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestActivitiesSuite(t *testing.T) {
	tu := newTestUtilsOrDie(t)
	gabi := newTestAccount(t, tu)
	diego := newTestAccount(t, tu)
	stranger := newTestAccount(t, tu)
	gabiGroup := newTestGroup(t, tu, gabi, diego)

	t.Run("get-activity", func(t *testing.T) {
		before := time.Now()
		activity, err := tu.activitiesRepository.CreateActivityInternal(gabi.Context, &models.ActivityPayload{
			GroupID: gabiGroup.ID,
			Type:    models.NoteAdded,
			Event:   "New event content",
		})
		after := time.Now()
		res, err := tu.groups.GetActivity(gabi.Context, &notesv1.GetActivityRequest{
			GroupId:    gabiGroup.ID,
			ActivityId: activity.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, string(models.NoteAdded), res.Activity.Type)
		require.Equal(t, "New event content", res.Activity.Event)
		require.GreaterOrEqual(t, res.Activity.CreatedAt.AsTime().Unix(), before.Unix())
		require.LessOrEqual(t, res.Activity.CreatedAt.AsTime().Unix(), after.Unix())
	})

	t.Run("stranger-cannot-get-activity", func(t *testing.T) {
		activity, err := tu.activitiesRepository.CreateActivityInternal(gabi.Context, &models.ActivityPayload{
			GroupID: gabiGroup.ID,
			Type:    models.NoteAdded,
			Event:   "New event content",
		})
		res, err := tu.groups.GetActivity(stranger.Context, &notesv1.GetActivityRequest{
			GroupId:    gabiGroup.ID,
			ActivityId: activity.ID,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("cannot-get-activity-because-of-validators", func(t *testing.T) {
		activity, err := tu.activitiesRepository.CreateActivityInternal(gabi.Context, &models.ActivityPayload{
			GroupID: gabiGroup.ID,
			Type:    models.NoteAdded,
			Event:   "New event content",
		})

		res, err := tu.groups.GetActivity(stranger.Context, &notesv1.GetActivityRequest{})
		require.Error(t, err)
		require.Nil(t, res)

		res, err = tu.groups.GetActivity(stranger.Context, &notesv1.GetActivityRequest{
			GroupId: gabiGroup.ID,
		})
		require.Error(t, err)
		require.Nil(t, res)

		res, err = tu.groups.GetActivity(stranger.Context, &notesv1.GetActivityRequest{
			ActivityId: activity.ID,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("list-activities", func(t *testing.T) {
		res, err := tu.groups.ListActivities(gabi.Context, &notesv1.ListActivitiesRequest{
			GroupId: gabiGroup.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, 4, len(res.Activities))
	})

	t.Run("stanger-cannot-list-activities", func(t *testing.T) {
		res, err := tu.groups.ListActivities(stranger.Context, &notesv1.ListActivitiesRequest{
			GroupId: gabiGroup.ID,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("cannot-list-activities-because-of-validators", func(t *testing.T) {
		res, err := tu.groups.ListActivities(gabi.Context, &notesv1.ListActivitiesRequest{})
		require.Error(t, err)
		require.Nil(t, res)

		res, err = tu.groups.ListActivities(gabi.Context, &notesv1.ListActivitiesRequest{
			GroupId: "",
		})
		require.Error(t, err)
		require.Nil(t, res)
	})
}
