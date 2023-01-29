package mongo

import (
	"context"
	"errors"
	"notes-service/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type repository struct {
	logger  *zap.Logger
	coll    *mongo.Collection
	newUUID func() string
}

func (repo *repository) findOneAndUpdate(ctx context.Context, query interface{}, update interface{}, result interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	repo.logger.Debug("find one and update", zap.Any("query", query), zap.Any("update", update))
	opts = append(opts, options.FindOneAndUpdate().SetReturnDocument(options.After))
	err := repo.coll.FindOneAndUpdate(ctx, query, update, opts...).Decode(result)
	if err != nil {
		return repo.mongoFindOneAndUpdateErrorToModelsError(query, update, err)
	}
	return nil
}

func (repo *repository) deleteOne(ctx context.Context, query interface{}, opts ...*options.DeleteOptions) error {
	repo.logger.Debug("delete one", zap.Any("query", query))
	res, err := repo.coll.DeleteOne(ctx, query, opts...)
	if err != nil {
		return repo.mongoDeleteOneErrorToModelsError(query, err)
	}
	if res.DeletedCount == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (repo *repository) findOne(ctx context.Context, query interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	repo.logger.Debug("find one", zap.Any("query", query))
	err := repo.coll.FindOne(ctx, query, opts...).Decode(result)
	if err != nil {
		return repo.mongoFindOneErrorToModelsError(query, err)
	}
	return nil
}

func (repo *repository) insertOne(ctx context.Context, payload interface{}, opts ...*options.InsertOneOptions) error {
	repo.logger.Debug("insert one", zap.Any("payload", payload))
	_, err := repo.coll.InsertOne(ctx, payload, opts...)
	if err != nil {
		return repo.mongoInsertOneErrorToModelsError(payload, err)
	}
	return nil
}

func (repo *repository) find(ctx context.Context, query interface{}, results interface{}, lo *models.ListOptions, opts ...*options.FindOptions) error {
	repo.logger.Debug("find", zap.Any("query", query))
	if lo == nil {
		lo = &models.ListOptions{Limit: 20, Offset: 0}
	}

	opts = append(opts, options.Find().SetLimit(lo.Limit).SetSkip(lo.Offset))

	res, err := repo.coll.Find(ctx, query, opts...)
	if err != nil {
		return repo.mongoFindErrorToModelsError(query, lo, err)
	}

	err = res.All(ctx, results)
	if err != nil {
		return repo.mongoFindErrorToModelsError(query, lo, err)
	}

	return nil
}

func (repo *repository) mongoFindOneErrorToModelsError(query interface{}, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.ErrNotFound
	}
	repo.logger.Error("find one failed", zap.Any("query", query), zap.Error(err))
	return models.ErrUnknown
}

func (repo *repository) mongoDeleteOneErrorToModelsError(query interface{}, err error) error {
	repo.logger.Error("delete one failed", zap.Any("query", query), zap.Error(err))
	return models.ErrUnknown
}

func (repo *repository) mongoInsertOneErrorToModelsError(query interface{}, err error) error {
	if mongo.IsDuplicateKeyError(err) {
		return models.ErrAlreadyExists
	}
	repo.logger.Error("find one failed", zap.Any("query", query), zap.Error(err))
	return models.ErrUnknown
}

func (repo *repository) mongoFindOneAndUpdateErrorToModelsError(query interface{}, update interface{}, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.ErrNotFound
	}
	repo.logger.Error("find one and update failed", zap.Any("query", query), zap.Any("update", update), zap.Error(err))
	return models.ErrUnknown
}

func (repo *repository) mongoFindErrorToModelsError(query interface{}, lo *models.ListOptions, err error) error {
	repo.logger.Error("find failed", zap.Any("query", query), zap.Int64("limit", lo.Limit), zap.Int64("offset", lo.Offset), zap.Error(err))
	return models.ErrUnknown
}
