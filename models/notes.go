package models

import (
	"context"
	"time"
)

type ImageBlock struct {
	Url     string `json:"url" bson:"url"`
	Caption string `json:"caption" bson:"caption"`
}

type CodeBlock struct {
	Snippet string `json:"snippet" bson:"snippet"`
	Caption string `json:"caption" bson:"caption"`
}

type NoteBlockType string

type NoteBlock struct {
	ID          string        `json:"id" bson:"_id"`
	Type        NoteBlockType `json:"type" bson:"type"`
	Heading     *string       `json:"heading,omitempty" bson:"heading,omitempty"`
	Paragraph   *string       `json:"paragraph,omitempty" bson:"paragraph,omitempty"`
	NumberPoint *string       `json:"numberPoint,omitempty" bson:"numberPoint,omitempty"`
	BulletPoint *string       `json:"bulletPoint,omitempty" bson:"bulletPoint,omitempty"`
	Math        *string       `json:"math,omitempty" bson:"math,omitempty"`
	Image       *ImageBlock   `json:"image,omitempty" bson:"image,omitempty"`
	Code        *CodeBlock    `json:"code,omitempty" bson:"code,omitempty"`
}

type KeywordType string

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
	Keyword string      `json:"keyword,omitempty" bson:"keyword,omitempty"`
	Type    KeywordType `json:"type,omitempty" bson:"type,omitempty"`
	URL     string      `json:"url,omitempty" bson:"url,omitempty"`
}

type Note struct {
	ID              string      `json:"id" bson:"_id"`
	Title           string      `json:"title" bson:"title"`
	AuthorAccountID string      `json:"authorAccountId" bson:"authorAccountId"`
	GroupID         string      `json:"groupId" bson:"groupId"`
	CreatedAt       time.Time   `json:"createdAt" bson:"createdAt"`
	ModifiedAt      time.Time   `json:"modifiedAt" bson:"modifiedAt"`
	AnalyzedAt      time.Time   `json:"analyzedAt" bson:"analyzedAt"`
	Keywords        []Keyword   `json:"keywords" bson:"keywords"`
	Blocks          []NoteBlock `json:"blocks" bson:"blocks"`
}

type CreateNotePayload struct {
	Title           string
	AuthorAccountID string
	GroupID         string
	FolderID        string
	Blocks          []CreateNoteBlockPayload
}

type CreateNoteBlockPayload struct {
	Type        NoteBlockType
	Heading     *string
	Paragraph   *string
	NumberPoint *string
	BulletPoint *string
	Math        *string
	Image       *ImageBlock
	Code        *CodeBlock
}

type ManyNotesFilter struct {
	// (Optional) List notes belonging to group.
	GroupID *string
	// (Optional) List notes belonging to account.
	AuthorAccountID *string
}

type UpdateBlockPayload struct {
	Type        NoteBlockType
	Heading     *string
	Paragraph   *string
	NumberPoint *string
	BulletPoint *string
	Math        *string
	Image       *ImageBlock
	Code        *CodeBlock
}

type NotesRepository interface {
	CreateNote(ctx context.Context, note *CreateNotePayload) (*Note, error)

	GetNote(ctx context.Context, noteID string) (*Note, error)

	UpdateNote(ctx context.Context, noteID string) (*Note, error)

	DeleteNote(ctx context.Context, noteID string) (*Note, error)

	ListNotes(ctx context.Context, filter *ManyNotesFilter, opts *ListOptions) ([]*Note, error)

	InsertBlock(ctx context.Context, noteID string, block *CreateNoteBlockPayload) (*Block, error)

	UpdateBlock(ctx context.Context, noteID string, blockID string, block *UpdateBlockPayload) (*Block, error)

	DeleteBlock(ctx context.Context, noteID string, blockID string) error
}
