package memory

import (
	"context"

	"github.com/hashicorp/go-memdb"

	"go.uber.org/zap"
)

// Database manages a connection with a Mongo database.
type Database struct {
	DB *memdb.MemDB

	logger *zap.Logger
}

func NewDatabase(ctx context.Context, shema *memdb.DBSchema, logger *zap.Logger) (*Database, error) {
	db, err := memdb.NewMemDB(shema)
	if err != nil {
		logger.Error("failed to create in-memory database", zap.Error(err))
		return nil, err
	}

	logger.Info("in-memory database creation successful")

	return &Database{
		DB:     db,
		logger: logger.Named("memory"),
	}, nil
}
