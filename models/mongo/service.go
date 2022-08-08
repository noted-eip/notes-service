package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

type Database struct {
	DB *mongo.Database

	client *mongo.Client
	logger *zap.Logger
}

func NewDatabase(ctx context.Context, uri string, name string, logger *zap.Logger) (*Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		logger.Error("failed to create mongo client", zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.Error("failed to connect to mongo server", zap.Error(err))
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Error("failed to ping mongo server", zap.Error(err))
		return nil, err
	}

	logger.Info("mongo database connection successful")

	return &Database{
		client: client,
		DB:     client.Database(name),
		logger: logger,
	}, nil
}

// Disconnect the TCP connection to the cluster.
func (s *Database) Disconnect(ctx context.Context) {
	if err := s.client.Disconnect(ctx); err != nil {
		s.logger.Error("failed to disconnect mongo client", zap.Error(err))
	}
}
