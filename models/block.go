// Package models defines the data types, payloads and repository interfaces
// of the accounts service.
package models

import (
	"context"
)

type Code struct {
	snippet *string
	lang    *string
}

type Image struct {
	url     *string
	caption *string
}

type Block struct {
	ID      string `json:"id" bson:"_id,omitempty"`
	NoteId  string `json:"noteId" bson:"noteId,omitempty"`
	Type    uint32 `json:"type" bson:"type,omitempty"`
	Content string `json:"content" bson:"content,omitempty"`
}

type BlockWithIndex struct {
	ID      string `json:"id" bson:"_id,omitempty"`
	NoteId  string `json:"noteId" bson:"noteId,omitempty"`
	Type    uint32 `json:"type" bson:"type,omitempty"`
	Index   uint32 `json:"index" bson:"index,omitempty"`
	Content string `json:"content" bson:"content,omitempty"`
}

type BlockFilter struct {
	BlockId string `json:"blockId" bson:"blockId,omitempty"`
	NoteId  string `json:"noteId" bson:"noteId,omitempty"`
}

// NotesRepository is safe for use in multiple goroutines.
type BlocksRepository interface {
	GetBlock(ctx context.Context, blockId *string) (*BlockWithIndex, error)

	GetBlocks(ctx context.Context, noteId *string) ([]*BlockWithIndex, error)

	Create(ctx context.Context, blockRequest *BlockWithIndex) error

	Update(ctx context.Context, blockId *string, blockRequest *BlockWithIndex) (*BlockWithIndex, error)

	Delete(ctx context.Context, blockId *string) error
}
