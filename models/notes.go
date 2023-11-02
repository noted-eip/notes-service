package models

import (
	"context"
	"time"
)

type NoteAction uint

type NoteIdentifier struct {
	NoteId     string
	ActionType NoteAction
	Metadata   interface{}
}

const (
	NoteUpdateKeyword NoteAction = iota
	NoteDeleteQuiz
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
	Styles      []TextStyle     `json:"styles,omitempty" bson:"styles,omitempty"`
	Thread      *[]BlockComment `json:"thread,omitempty" bson:"thread,omitempty"`
}

type BlockComment struct {
	ID              string `json:"id" bson:"id"`
	AuthorAccountID string `json:"authorAccountId" bson:"authorAccountId"`
	Content         string `json:"content,omitempty" bson:"content,omitempty"`
}

type TextStyle struct {
	Style    string   `json:"style,omitempty" bson:"style,omitempty"`
	Position Position `json:"pos,omitempty" bson:"pos,omitempty"`
	Color    *Color   `json:"color,omitempty" bson:"color,omitempty"`
}

type Position struct {
	Start  int64 `json:"styles,omitempty" bson:"styles,omitempty"`
	Length int64 `json:"length,omitempty" bson:"length,omitempty"`
}

type TextStyleType string

type Color struct {
	R int32 `json:"r,omitempty" bson:"r,omitempty"`
	G int32 `json:"g,omitempty" bson:"g,omitempty"`
	B int32 `json:"b,omitempty" bson:"b,omitempty"`
}

const (
	Unknown      string = "Inconnu"
	Person       string = "Personne"
	Location     string = "Lieu"
	Organization string = "Organisation"
	Event        string = "Evenement"
	WorkOfArt    string = "Chef d'oeuvre"
	ConsumerGood string = "Bien de consommation"
	Other        string = "Autre"
	PhoneNumber  string = "Numéro de téléphone"
	Address      string = "Addresse"
	Date         string = "Date"
	Number       string = "Nombre"
	Price        string = "Prix"
)

type Keyword struct {
	Keyword  string `json:"keyword,omitempty" bson:"keyword,omitempty"`
	Type     string `json:"type,omitempty" bson:"type,omitempty"`
	URL      string `json:"url,omitempty" bson:"url,omitempty"`
	Summary  string `json:"summary,omitempty" bson:"summary,omitempty"`
	ImageURL string `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
}

type Note struct {
	ID                          string       `json:"id" bson:"_id"`
	Title                       string       `json:"title" bson:"title"`
	AuthorAccountID             string       `json:"authorAccountId" bson:"authorAccountId"`
	GroupID                     string       `json:"groupId" bson:"groupId"`
	CreatedAt                   time.Time    `json:"createdAt" bson:"createdAt"`
	ModifiedAt                  *time.Time   `json:"modifiedAt" bson:"modifiedAt"`
	AnalyzedAt                  *time.Time   `json:"analyzedAt" bson:"analyzedAt"`
	Keywords                    []*Keyword   `json:"keywords" bson:"keywords"`
	Blocks                      *[]NoteBlock `json:"blocks" bson:"blocks"`
	Quizs                       *[]Quiz      `json:"quizs" bson:"quizs"`
	Lang                        string       `json:"lang" bson:"lang"`
	AccountsWithEditPermissions []string     `json:"accountsWithEditPermissions" bson:"accountsWithEditPermissions"`
}

type Quiz struct {
	ID            string         `json:"id" bson:"id"`
	CreatedAt     time.Time      `json:"createdAt" bson:"createdAt"`
	QuizQuestions []QuizQuestion `json:"questions,omitempty" bson:"questions,omitempty"`
}

type Summary struct {
	Content string `json:"content,omitempty" bson:"content,omitempty"`
}

type QuizQuestion struct {
	Question  string   `json:"question,omitempty" bson:"question,omitempty"`
	Answers   []string `json:"answers,omitempty" bson:"answers,omitempty"`
	Solutions []string `json:"solutions,omitempty" bson:"solutions,omitempty"`
}

func (note *Note) FindBlock(blockID string) *NoteBlock {
	for i := 0; i < len(*note.Blocks); i++ {
		if (*note.Blocks)[i].ID == blockID {
			return &(*note.Blocks)[i]
		}
	}
	return nil
}

func (block *NoteBlock) FindComment(commentID string) *BlockComment {
	for i := 0; i < len(*block.Thread); i++ {
		if (*block.Thread)[i].ID == commentID {
			return &(*block.Thread)[i]
		}
	}
	return nil
}

type CreateNotePayload struct {
	Title           string
	AuthorAccountID string
	GroupID         string
	FolderID        string
	Lang            string
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
	Keywords []*Keyword `json:"keywords" bson:"keywords"`
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
	StoreNewQuiz(ctx context.Context, filter *OneNoteFilter, payload *Quiz, accountID string) (*Quiz, error)
	ListQuizs(ctx context.Context, filter *OneNoteFilter, accountID string) (*[]Quiz, error)
	DeleteQuiz(ctx context.Context, filter *OneNoteFilter, quizID string, accountID string) error

	// Permisions
	GrantNoteEditPermission(ctx context.Context, filter *OneNoteFilter, AccountID string, RecipientAccountId string) error

	// Blocks
	InsertBlock(ctx context.Context, filter *OneNoteFilter, payload *InsertNoteBlockPayload, accountID string) (*NoteBlock, error)
	UpdateBlock(ctx context.Context, filter *OneBlockFilter, payload *UpdateBlockPayload, accountID string) (*NoteBlock, error)
	GetBlock(ctx context.Context, filter *OneBlockFilter, accountID string) (*NoteBlock, error)
	DeleteBlock(ctx context.Context, filter *OneBlockFilter, accountID string) error
	CreateBlockComment(ctx context.Context, filter *OneBlockFilter, payload *BlockComment, accountID string) (*BlockComment, error)
	DeleteBlockComment(ctx context.Context, filter *OneBlockFilter, payload *BlockComment, accountID string) (*BlockComment, error)
	ListBlockComments(ctx context.Context, filter *OneBlockFilter, lo *ListOptions, accountID string) (*[]BlockComment, error)

	// Utils
	RemoveEditPermissions(ctx context.Context, filter *OneNoteFilter, accountID string) error
}
