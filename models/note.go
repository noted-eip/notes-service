// Package models defines the data types, payloads and repository interfaces
// of the accounts service.
package models

import (
	"context"
	"time"
)

type Note struct {
	ID               string    `json:"id" bson:"_id,omitempty"`
	AuthorId         string    `json:"authorId" bson:"authorId,omitempty"`
	Title            string    `json:"title" bson:"title,omitempty"`
	Blocks           []Block   `json:"blocks" bson:"blocks,omitempty"`
	CreationDate     time.Time `json:"creationDate" bson:"creationDate,omitempty"`
	ModificationDate time.Time `json:"modificationDate" bson:"modificationDate,omitempty"`
}

type NotePayload struct {
	ID       string  `json:"id" bson:"_id,omitempty"`
	AuthorId string  `json:"authorId" bson:"authorId,omitempty"`
	Title    string  `json:"title" bson:"title,omitempty"`
	Blocks   []Block `json:"blocks" bson:"blocks,omitempty"`
}

// NotesRepository is safe for use in multiple goroutines.
type NotesRepository interface {
	Create(ctx context.Context, noteRequest *NotePayload) (*Note, error)

	Get(ctx context.Context, noteId string) (*Note, error)

	Delete(ctx context.Context, noteId string) error

	Update(ctx context.Context, noteId string, noteRequest *NotePayload) error

	List(ctx context.Context, authorId string) ([]*Note, error)
}
