// Package models defines the data types, payloads and repository interfaces
// of the accounts service.
package models

import (
	"context"
)

type Block struct {
	ID      string  `json:"id" bson:"_id,omitempty"`
	NoteId  string  `json:"noteId" bson:"noteId,omitempty"`
	Type    uint32  `json:"type" bson:"type,omitempty"`
	Content *string `json:"content" bson:"content,omitempty"`
}

type BlockWithIndex struct {
	ID      string  `json:"id" bson:"_id,omitempty"`
	NoteId  string  `json:"noteId" bson:"noteId,omitempty"`
	Type    uint32  `json:"type" bson:"type,omitempty"`
	Index   uint32  `json:"inxed" bson:"inxed,omitempty"`
	Content *string `json:"content" bson:"content,omitempty"`
}

type BlockFilter struct {
	BlockId string `json:"blockId" bson:"blockId,omitempty"`
	NoteId  string `json:"noteId" bson:"noteId,omitempty"`
}

// NotesRepository is safe for use in multiple goroutines.
type BlocksRepository interface {
	Create(ctx context.Context, blockRequest *BlockWithIndex) error

	Delete(ctx context.Context, filter *BlockFilter) (*NoteWithBlocks, error)

	Update(ctx context.Context, filter *BlockFilter, blockRequest *BlockWithIndex) (*NoteWithBlocks, error)
}
