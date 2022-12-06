// Package models defines the data types, payloads and repository interfaces
// of the accounts service.
package models

import (
	"context"
)

type Widget struct {
	Title       string `json:"title" bson:"title,omitempty"`
	Description string `json:"description" bson:"description,omitempty"`
}

//Inheritance from Widget
type WebsiteWidget struct {
	Widget
	Url  string `json:"url" bson:"url,omitempty"`
	Type string `json:"Type" bson:"Type,omitempty"`
}

type Note struct {
	ID       string  `json:"id" bson:"_id,omitempty"`
	AuthorId string  `json:"authorId" bson:"authorId,omitempty"`
	Title    string  `json:"title" bson:"title,omitempty"`
	Blocks   []Block `json:"blocks" bson:"blocks,omitempty"`
}

// NotesRepository is safe for use in multiple goroutines.
type NotesRepository interface {
	Create(ctx context.Context, noteRequest *Note) (*Note, error)

	Get(ctx context.Context, noteId *string) (*Note, error)

	Delete(ctx context.Context, noteId *string) error

	Update(ctx context.Context, noteId *string, noteRequest *Note) error

	List(ctx context.Context, authorId *string) (*[]Note, error)
}
