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

func (repo *repository) aggregate(ctx context.Context, pipeline interface{}, result interface{}, opts ...*options.AggregateOptions) error {
	repo.logger.Debug("aggregate", zap.Any("pipeline", pipeline))

	cur, err := repo.coll.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return repo.mongoAggregateErrorToModelsError(pipeline, err)
	}

	cur.All(ctx, result)
	if err != nil {
		return repo.mongoAggregateErrorToModelsError(pipeline, err)
	}

	return nil
}

func (repo *repository) updateOne(ctx context.Context, query interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	repo.logger.Debug("update one", zap.Any("query", query), zap.Any("update", update))
	res, err := repo.coll.UpdateOne(ctx, query, update, opts...)

	if err != nil {
		return repo.mongoUpdateOneErrorToModelsError(query, update, err)
	}
	if res.ModifiedCount == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (repo *repository) updateMany(ctx context.Context, query interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	repo.logger.Debug("update many", zap.Any("query", query), zap.Any("update", update))
	res, err := repo.coll.UpdateMany(ctx, query, update, opts...)

	if err != nil {
		return 0, repo.mongoUpdateManyErrorToModelsError(query, update, err)
	}
	return res.ModifiedCount, nil
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

func (repo *repository) deleteMany(ctx context.Context, query interface{}, opts ...*options.DeleteOptions) error {
	repo.logger.Debug("delete many", zap.Any("query", query))
	res, err := repo.coll.DeleteMany(ctx, query, opts...)
	if err != nil {
		return repo.mongoDeleteManyErrorToModelsError(query, err)
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

	opts = append(opts, options.Find().SetLimit(int64(lo.Limit)).SetSkip(int64(lo.Offset)))

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

func (repo *repository) findAll(ctx context.Context, query interface{}, results interface{}, opts ...*options.FindOptions) error {
	repo.logger.Debug("find", zap.Any("query", query))

	res, err := repo.coll.Find(ctx, query, opts...)
	if err != nil {
		return repo.mongoFindAllErrorToModelsError(query, err)
	}

	err = res.All(ctx, results)
	if err != nil {
		return repo.mongoFindAllErrorToModelsError(query, err)
	}

	return nil
}

func (repo *repository) mongoAggregateErrorToModelsError(query interface{}, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.ErrNotFound
	}
	repo.logger.Error("aggregate failed", zap.Any("query", query), zap.Error(err))
	return models.ErrUnknown
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

func (repo *repository) mongoDeleteManyErrorToModelsError(query interface{}, err error) error {
	repo.logger.Error("delete many failed", zap.Any("query", query), zap.Error(err))
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

func (repo *repository) mongoUpdateOneErrorToModelsError(query interface{}, update interface{}, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.ErrNotFound
	}
	repo.logger.Error("update one failed", zap.Any("query", query), zap.Any("update", update), zap.Error(err))
	return models.ErrUnknown
}

func (repo *repository) mongoUpdateManyErrorToModelsError(query interface{}, update interface{}, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.ErrNotFound
	}
	repo.logger.Error("update many failed", zap.Any("query", query), zap.Any("update", update), zap.Error(err))
	return models.ErrUnknown
}

func (repo *repository) mongoFindErrorToModelsError(query interface{}, lo *models.ListOptions, err error) error {
	repo.logger.Error("find failed", zap.Any("query", query), zap.Int32("limit", lo.Limit), zap.Int32("offset", lo.Offset), zap.Error(err))
	return models.ErrUnknown
}

func (repo *repository) mongoFindAllErrorToModelsError(query interface{}, err error) error {
	repo.logger.Error("findAll failed", zap.Any("query", query), zap.Error(err))
	return models.ErrUnknown
}
