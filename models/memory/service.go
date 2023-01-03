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

func NewDatabase(ctx context.Context, logger *zap.Logger) (*Database, error) {
	db, err := memdb.NewMemDB(&memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"note": {
				Name: "note",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"author_id": {
						Name:    "author_id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "AuthorId"},
					},
					"title": {
						Name:    "title",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Title"},
					},
					"blocks": {
						Name:    "blocks",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Blocks"},
					},
				},
			},
			"block": {
				Name: "block",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"note_id": {
						Name:    "note_id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "NoteId"},
					},
					"type": {
						Name:    "type",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Type"},
					},
					"index": {
						Name:    "index",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Index"},
					},
					"content": {
						Name:    "content",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Content"},
					},
				},
			},
		},
	})
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
