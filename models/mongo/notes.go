package mongo

import (
	"context"
	"notes-service/models"
	"time"

	notesv1 "notes-service/protorepo/noted/notes/v1"

	"github.com/jaevor/go-nanoid"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type notesRepository struct {
	repository
}

func NewNotesRepository(db *mongo.Database, logger *zap.Logger) models.NotesRepository {
	newUUID, err := nanoid.Standard(21)
	if err != nil {
		panic(err)
	}

	return &notesRepository{
		repository: repository{
			logger:  logger.Named("mongo").Named("notes"),
			coll:    db.Collection("notes"),
			newUUID: newUUID,
		},
	}
}

func (repo *notesRepository) CreateNote(ctx context.Context, payload *models.CreateNotePayload, accountID string) (*models.Note, error) {
	note := &models.Note{
		ID:              repo.newUUID(),
		Title:           payload.Title,
		AuthorAccountID: accountID,
		GroupID:         payload.GroupID,
		CreatedAt:       time.Now(),
		ModifiedAt:      time.Now(),
		AnalyzedAt:      time.Now(),
		Keywords:        []models.Keyword{},
		Blocks:          []models.NoteBlock{},
	}

	err := repo.insertOne(ctx, note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (repo *notesRepository) GetNote(ctx context.Context, filter *models.OneNoteFilter, accountID string) (*models.Note, error) {
	note := &models.Note{}
	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "group_id", Value: filter.GroupID},
	}

	err := repo.findOne(ctx, query, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (repo *notesRepository) UpdateNote(ctx context.Context, filter *models.OneNoteFilter, payload *models.UpdateNotePayload, accountID string) (*models.Note, error) {
	note := &models.Note{}
	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
	}
	update := bson.D{
		{Key: "$set", Value: payload},
		{Key: "$set", Value: bson.D{
			{Key: "modifiedAt", Value: time.Now()},
		}}}

	err := repo.findOneAndUpdate(ctx, query, update, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (repo *notesRepository) DeleteNote(ctx context.Context, filter *models.OneNoteFilter, accountID string) error {
	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
	}

	return repo.deleteOne(ctx, query)
}

func (repo *notesRepository) ListNotesInternal(ctx context.Context, filter *models.ManyNotesFilter, lo *models.ListOptions) ([]*models.Note, error) {
	notes := make([]*models.Note, 0)

	query := bson.D{}
	if filter != nil {
		if filter.AuthorAccountID != nil {
			query = append(query, bson.E{Key: "authorAccountId", Value: filter.AuthorAccountID})
		}
		if filter.GroupID != nil {
			query = append(query, bson.E{Key: "groupId", Value: filter.GroupID})
		}
	}

	err := repo.find(ctx, query, &notes, lo)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (repo *notesRepository) InsertBlock(ctx context.Context, filter *models.OneNoteFilter, payload *models.CreateNoteBlockPayload, accountID string) (*models.NoteBlock, error) {
	note := &models.Note{}
	block := &models.NoteBlock{}
	copier.Copy(block, payload)
	block.ID = repo.newUUID()

	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		// Only the author can modify the note.
		{Key: "authorAccountId", Value: accountID},
	}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "blocks", Value: bson.D{
				{Key: "$each", Value: bson.A{block}},
				{Key: "$position", Value: payload.Index},
			}},
		}},
	}

	err := repo.findOneAndUpdate(ctx, query, update, note)
	if err != nil {
		return nil, err
	}

	return note.FindBlock(block.ID), nil
}

func (repo *notesRepository) UpdateBlock(ctx context.Context, filter *models.OneBlockFilter, payload *models.UpdateBlockPayload, accountID string) (*models.NoteBlock, error) {
	note := &models.Note{}

	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		{Key: "blocks.id", Value: filter.BlockID},
		// Only the author can modify the note.
		{Key: "authorAccountId", Value: accountID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "blocks.$.type", Value: payload.Type},
			updateBlockPayloadToDocument(payload),
		}},
	}

	err := repo.findOneAndUpdate(ctx, query, update, note)
	if err != nil {
		return nil, err
	}

	return note.FindBlock(filter.BlockID), nil
}

func (repo *notesRepository) DeleteBlock(ctx context.Context, filter *models.OneBlockFilter, accountID string) error {
	note := &models.Note{}

	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		{Key: "blocks.id", Value: filter.BlockID},
	}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "blocks", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "id", Value: filter.BlockID},
				}},
			}},
		}}}

	return repo.findOneAndUpdate(ctx, query, update, note)
}

func updateBlockPayloadToDocument(payload *models.UpdateBlockPayload) bson.E {
	switch payload.Type {
	case notesv1.Block_TYPE_HEADING_1.String():
		return bson.E{Key: "heading", Value: payload.Heading}
	case notesv1.Block_TYPE_HEADING_2.String():
		return bson.E{Key: "heading", Value: payload.Heading}
	case notesv1.Block_TYPE_HEADING_3.String():
		return bson.E{Key: "heading", Value: payload.Heading}
	case notesv1.Block_TYPE_BULLET_POINT.String():
		return bson.E{Key: "bulletPoint", Value: payload.BulletPoint}
	case notesv1.Block_TYPE_NUMBERED_POINT.String():
		return bson.E{Key: "numberPoint", Value: payload.NumberPoint}
	case notesv1.Block_TYPE_PARAGRAPH.String():
		return bson.E{Key: "paragraph", Value: payload.Paragraph}
	case notesv1.Block_TYPE_MATH.String():
		return bson.E{Key: "math", Value: payload.Math}
	case notesv1.Block_TYPE_IMAGE.String():
		return bson.E{Key: "image", Value: payload.Image}
	case notesv1.Block_TYPE_CODE.String():
		return bson.E{Key: "code", Value: payload.Code}
	}
	return bson.E{Key: "", Value: nil}
}
