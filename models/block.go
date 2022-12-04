// Package models defines the data types, payloads and repository interfaces
// of the accounts service.
package models

import (
	"context"
)

type Code struct {
	Snippet *string
	Lang    *string
}

type Image struct {
	Url     *string
	Caption *string
}

type Block struct {
	ID      string `json:"id" bson:"_id,omitempty"`
	NoteId  string `json:"noteId" bson:"noteId,omitempty"`
	Type    uint32 `json:"type" bson:"type,omitempty"`
	Index   uint32 `json:"index" bson:"index,omitempty"`
	Content string `json:"content" bson:"content,omitempty"`
}

// NotesRepository is safe for use in multiple goroutines.
type BlocksRepository interface {
	GetBlock(ctx context.Context, blockId string) (*Block, error)

	GetBlocks(ctx context.Context, noteId string) ([]*Block, error)

	Create(ctx context.Context, blockRequest *Block) (*string, error)

	Update(ctx context.Context, blockId string, blockRequest *Block) (*Block, error)

	DeleteBlock(ctx context.Context, blockId string) error

	DeleteBlocks(ctx context.Context, noteId string) error
}
