// Package models defines the data types, payloads and repository interfaces
// of the accounts service.
package models

import (
	"context"

	"github.com/google/uuid"
)

type NoteWithBlocks struct {
	ID       uuid.UUID `json:"id" bson:"_id,omitempty"`
	AuthorId string    `json:"authorId" bson:"authorId,omitempty"`
	Title    string    `json:"title" bson:"title,omitempty"`
	Blocks   []Block   `json:"blocks" bson:"blocks,omitempty"`
}

/*type NoteProprieties struct {
	ID       uuid.UUID `json:"id" bson:"_id,omitempty"`
	AuthorId string    `json:"authorId" bson:"authorId,omitempty"`
	Title    *string   `json:"title" bson:"title,omitempty"`
}*/

type NoteFilter struct {
	ID       uuid.UUID `json:"id" bson:"_id,omitempty"`
	AuthorId string    `json:"authorId" bson:"authorId,omitempty"`
}

// NotesRepository is safe for use in multiple goroutines.
type NotesRepository interface {
	Create(ctx context.Context, noteRequest *NoteWithBlocks) (*NoteWithBlocks, error)

	Get(ctx context.Context, filter *NoteFilter) (*NoteWithBlocks, error)

	Delete(ctx context.Context, filter *NoteFilter) error

	Update(ctx context.Context, filter *NoteFilter, noteRequest *NoteWithBlocks) error

	List(ctx context.Context, filter *NoteFilter) (*[]NoteWithBlocks, error)
}
