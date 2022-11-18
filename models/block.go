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

type BlockWithTags struct {
	ID      string   `json:"id" bson:"_id,omitempty"`
	NoteId  string   `json:"noteId" bson:"noteId,omitempty"`
	Type    uint32   `json:"type" bson:"type,omitempty"`
	Index   uint32   `json:"index" bson:"index,omitempty"`
	Content string   `json:"content" bson:"content,omitempty"`
	Tags    []string `json:"tags" bson:"tags,omitempty"`
}

// NotesRepository is safe for use in multiple goroutines.
type BlocksRepository interface {
	GetBlock(ctx context.Context, blockId *string) (*BlockWithIndex, error)

	GetBlocks(ctx context.Context, noteId *string) ([]*BlockWithIndex, error)

	GetTagsByFilter(ctx context.Context, noteId *string) (*BlockWithTags, error)

	Create(ctx context.Context, blockRequest *BlockWithTags) (*string, error)

	Update(ctx context.Context, blockId *string, blockRequest *BlockWithIndex) (*BlockWithIndex, error)

	DeleteBlock(ctx context.Context, blockId *string) error

	DeleteBlocks(ctx context.Context, noteId *string) error
}
