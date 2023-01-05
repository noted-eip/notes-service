// Package models defines the data types, payloads and repository interfaces
// of the accounts service.
package models

import (
	"context"
	"time"
)

type KeywordType string

// Note: Enum just to be sure that there won't be any other strings used in the `Type` field, dk if it's useful
const (
	Unknown      KeywordType = "Unknown"
	Person       KeywordType = "Person"
	Location     KeywordType = "Location"
	Organization KeywordType = "Organization"
	Event        KeywordType = "Event"
	WorkOfArt    KeywordType = "Work of art"
	ConsumerGood KeywordType = "Consumer good"
	Other        KeywordType = "Other"
	PhoneNumber  KeywordType = "Phone number"
	Address      KeywordType = "Address"
	Date         KeywordType = "Date"
	Number       KeywordType = "Number"
	Price        KeywordType = "Price"
)

type Keyword struct {
	Keyword string      `json:"keyword" bson:"keyword,omitempty"`
	Type    KeywordType `json:"type" bson:"type,omitempty"`
	URL     string      `json:"url" bson:"url,omitempty"`
}

type Keywords []Keyword

type Note struct {
	ID               string    `json:"id" bson:"_id,omitempty"`
	AuthorId         string    `json:"authorId" bson:"authorId,omitempty"`
	Title            string    `json:"title" bson:"title,omitempty"`
	Blocks           []Block   `json:"blocks" bson:"blocks,omitempty"`
	Keywords         Keywords  `json:"keywords" bson:"keywords,omitempty"`
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

	List(ctx context.Context, authorId string) (*[]Note, error)
}
