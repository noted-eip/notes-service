package models

import (
	"context"
	"time"
)

type Action uint

type NoteIdentifier struct {
	NoteId     string
	ActionType Action
}

const (
	NoteUpdateKeyword Action = 1
	// Put in enum the other type of actions
	//...
)

type NoteBlockImage struct {
	Url     string `json:"url" bson:"url"`
	Caption string `json:"caption" bson:"caption"`
}

type NoteBlockCode struct {
	Snippet string `json:"snippet" bson:"snippet"`
	Lang    string `json:"lang" bson:"lang"`
}

type NoteBlockType = string

type NoteBlock struct {
	ID          string          `json:"id" bson:"id"`
	Type        NoteBlockType   `json:"type" bson:"type"`
	Heading     *string         `json:"heading,omitempty" bson:"heading,omitempty"`
	Paragraph   *string         `json:"paragraph,omitempty" bson:"paragraph,omitempty"`
	NumberPoint *string         `json:"numberPoint,omitempty" bson:"numberPoint,omitempty"`
	BulletPoint *string         `json:"bulletPoint,omitempty" bson:"bulletPoint,omitempty"`
	Math        *string         `json:"math,omitempty" bson:"math,omitempty"`
	Image       *NoteBlockImage `json:"image,omitempty" bson:"image,omitempty"`
	Code        *NoteBlockCode  `json:"code,omitempty" bson:"code,omitempty"`
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
	ModifiedAt      *time.Time  `json:"modifiedAt" bson:"modifiedAt"`
	AnalyzedAt      *time.Time  `json:"analyzedAt" bson:"analyzedAt"`
	Keywords        []Keyword   `json:"keywords" bson:"keywords"`
	Blocks          []NoteBlock `json:"blocks" bson:"blocks"`
}

func (note *Note) FindBlock(blockID string) *NoteBlock {
	for i := 0; i < len(note.Blocks); i++ {
		if note.Blocks[i].ID == blockID {
			return &note.Blocks[i]
		}
	}
	return nil
}

type CreateNotePayload struct {
	Title           string
	AuthorAccountID string
	GroupID         string
	FolderID        string
	Blocks          []NoteBlock
}

type InsertNoteBlockPayload struct {
	Block NoteBlock
	Index uint
}

type ManyNotesFilter struct {
	// (Optional) List notes belonging to group.
	GroupID string
	// (Optional) List notes belonging to account.
	AuthorAccountID string
}

type UpdateBlockPayload struct {
	Block NoteBlock
}

type UpdateNotePayload struct {
	Title  string       `json:"title,omitempty" bson:"title,omitempty"`
	Blocks *[]NoteBlock `json:"blocks,omitempty" bson:"blocks,omitempty"`

	// TODO: Remove
	Keywords []Keyword `json:"keywords" bson:"keywords"`
}

type UpdateNoteGroupPayload struct {
	GroupID string `json:"groupId" bson:"groupId"`
}

type OneNoteFilter struct {
	GroupID string
	NoteID  string
}

type OneBlockFilter struct {
	GroupID string
	NoteID  string
	BlockID string
}

type NotesRepository interface {
	// Notes
	CreateNote(ctx context.Context, payload *CreateNotePayload, accountID string) (*Note, error)
	GetNote(ctx context.Context, filter *OneNoteFilter, accountID string) (*Note, error)
	UpdateNote(ctx context.Context, filter *OneNoteFilter, payload *UpdateNotePayload, accountID string) (*Note, error)
	UpdateNotesInternal(ctx context.Context, filter *ManyNotesFilter, payload interface{}) (*Note, error)
	DeleteNote(ctx context.Context, filter *OneNoteFilter, accountID string) error
	DeleteNotes(ctx context.Context, filter *ManyNotesFilter) error
	ListNotesInternal(ctx context.Context, filter *ManyNotesFilter, opts *ListOptions) ([]*Note, error)
	ListAllNotesInternal(ctx context.Context, filter *ManyNotesFilter) ([]*Note, error)

	// Blocks
	InsertBlock(ctx context.Context, filter *OneNoteFilter, payload *InsertNoteBlockPayload, accountID string) (*NoteBlock, error)
	UpdateBlock(ctx context.Context, filter *OneBlockFilter, payload *UpdateBlockPayload, accountID string) (*NoteBlock, error)
	DeleteBlock(ctx context.Context, filter *OneBlockFilter, accountID string) error
}
