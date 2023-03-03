package mongo

import (
	"context"
	"notes-service/models"
	"time"

	"github.com/jaevor/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type activitiesRepository struct {
	repository
}

func NewActivitiesRepository(db *mongo.Database, logger *zap.Logger) models.ActivitiesRepository {
	newUUID, err := nanoid.Standard(21)
	if err != nil {
		panic(err)
	}

	return &activitiesRepository{
		repository: repository{
			logger:  logger.Named("mongo").Named("activities"),
			coll:    db.Collection("activities"),
			newUUID: newUUID,
		},
	}
}

func (repo *activitiesRepository) ListActivitiesInternal(ctx context.Context, filter *models.ManyActivitiesFilter, lo *models.ListOptions) ([]*models.Activity, error) {
	activities := make([]*models.Activity, 0)

	query := bson.D{
		{Key: "groupId", Value: filter.GroupID},
	}

	err := repo.find(ctx, query, &activities, lo)
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func (repo *activitiesRepository) GetActivityInternal(ctx context.Context, filter *models.OneActivityFilter) (*models.Activity, error) {
	activity := &models.Activity{}

	query := bson.D{
		{Key: "_id", Value: filter.ActivityId},
		{Key: "groupId", Value: filter.GroupID},
	}

	err := repo.findOne(ctx, query, activity)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

func (repo *activitiesRepository) CreateActivityInternal(ctx context.Context, payload *models.ActivityPayload) (*models.Activity, error) {
	activity := &models.Activity{
		ID:        repo.newUUID(),
		GroupID:   payload.GroupID,
		Type:      string(payload.Type),
		Event:     payload.Event,
		CreatedAt: time.Now(),
	}

	err := repo.insertOne(ctx, activity)
	if err != nil {
		return nil, err
	}

	return activity, nil
}
