package mongo

import (
	"context"
	"notes-service/models"
	"time"

	notesv1 "notes-service/protorepo/noted/notes/v1"

	"github.com/jaevor/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type notesRepository struct {
	repository
}

func (repo *notesRepository) GrantNoteEditPermission(ctx context.Context, filter *models.OneNoteFilter, AccountID string, recipientAccountID string) error {
	note := &models.Note{}
	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		{Key: "authorAccountId", Value: AccountID},
	}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "accountsWithEditPermissions", Value: recipientAccountID},
		}},
	}
	err := repo.findOneAndUpdate(ctx, query, update, note)
	if err != nil {
		return err
	}

	return nil
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
	for i := range payload.Blocks {
		payload.Blocks[i].ID = repo.newUUID()
	}

	now := time.Now()
	note := &models.Note{
		ID:                          repo.newUUID(),
		Title:                       payload.Title,
		AuthorAccountID:             accountID,
		GroupID:                     payload.GroupID,
		CreatedAt:                   now,
		ModifiedAt:                  nil,
		AnalyzedAt:                  nil,
		Keywords:                    []*models.Keyword{},
		Blocks:                      payload.Blocks,
		AccountsWithEditPermissions: []string{accountID},
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
		{Key: "groupId", Value: filter.GroupID},
	}

	err := repo.findOne(ctx, query, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (repo *notesRepository) UpdateNotesInternal(ctx context.Context, filter *models.ManyNotesFilter, payload interface{}) (*models.Note, error) {
	note := &models.Note{}
	query := bson.D{
		{Key: "groupId", Value: filter.GroupID},
		{Key: "authorAccountId", Value: filter.AuthorAccountID},
	}
	update := bson.D{
		{Key: "$set", Value: payload},
		{Key: "$set", Value: bson.D{
			{Key: "modifiedAt", Value: time.Now()},
		}}}

	_, err := repo.updateMany(ctx, query, update)
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
		// {Key: "authorAccountId", Value: accountID}, // NOTE: Removed to manage notes permissions
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
		{Key: "authorAccountId", Value: accountID},
	}

	return repo.deleteOne(ctx, query)
}

func (repo *notesRepository) DeleteNotes(ctx context.Context, filter *models.ManyNotesFilter) error {
	query := bson.D{}
	if filter != nil {
		if filter.AuthorAccountID != "" {
			query = append(query, bson.E{Key: "authorAccountId", Value: filter.AuthorAccountID})
		}
		if filter.GroupID != "" {
			query = append(query, bson.E{Key: "groupId", Value: filter.GroupID})
		}
	}

	return repo.deleteMany(ctx, query)
}

func (repo *notesRepository) ListNotesInternal(ctx context.Context, filter *models.ManyNotesFilter, lo *models.ListOptions) ([]*models.Note, error) {
	notes := make([]*models.Note, 0)

	query := bson.D{}
	if filter != nil {
		if filter.AuthorAccountID != "" {
			query = append(query, bson.E{Key: "authorAccountId", Value: filter.AuthorAccountID})
		}
		if filter.GroupID != "" {
			query = append(query, bson.E{Key: "groupId", Value: filter.GroupID})
		}
	}
	requieredFields := bson.D{{Key: "blocks", Value: 0}, {Key: "keywords", Value: 0}}
	opts := options.Find().SetProjection(requieredFields)

	err := repo.find(ctx, query, &notes, lo, opts)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (repo *notesRepository) ListAllNotesInternal(ctx context.Context, filter *models.ManyNotesFilter) ([]*models.Note, error) {
	notes := make([]*models.Note, 0)

	query := bson.D{}
	if filter != nil {
		if filter.AuthorAccountID != "" {
			query = append(query, bson.E{Key: "authorAccountId", Value: filter.AuthorAccountID})
		}
		if filter.GroupID != "" {
			query = append(query, bson.E{Key: "groupId", Value: filter.GroupID})
		}
	}

	err := repo.findAll(ctx, query, &notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (repo *notesRepository) InsertBlock(ctx context.Context, filter *models.OneNoteFilter, payload *models.InsertNoteBlockPayload, accountID string) (*models.NoteBlock, error) {
	note := &models.Note{}
	payload.Block.ID = repo.newUUID()

	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		{Key: "authorAccountId", Value: accountID},
	}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "blocks", Value: bson.D{
				{Key: "$each", Value: bson.A{payload.Block}},
				{Key: "$position", Value: payload.Index},
			}},
		}},
	}

	err := repo.findOneAndUpdate(ctx, query, update, note)
	if err != nil {
		return nil, err
	}

	return note.FindBlock(payload.Block.ID), nil
}

func (repo *notesRepository) UpdateBlock(ctx context.Context, filter *models.OneBlockFilter, payload *models.UpdateBlockPayload, accountID string) (*models.NoteBlock, error) {
	note := &models.Note{}

	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		{Key: "blocks.id", Value: filter.BlockID},
		{Key: "authorAccountId", Value: accountID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "blocks.$.type", Value: payload.Block.Type},
			updateBlockPayloadToDocument(payload),
		}},
	}

	err := repo.findOneAndUpdate(ctx, query, update, note)
	if err != nil {
		return nil, err
	}

	return note.FindBlock(filter.BlockID), nil
}

func (repo *notesRepository) GetBlock(ctx context.Context, filter *models.OneBlockFilter, accountID string) (*models.NoteBlock, error) {
	note := &models.Note{}
	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
	}

	err := repo.findOne(ctx, query, note)
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
		{Key: "authorAccountId", Value: accountID},
	}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "blocks", Value: bson.D{
				{Key: "id", Value: filter.BlockID},
			}},
		}}}

	return repo.findOneAndUpdate(ctx, query, update, note)
}

func updateBlockPayloadToDocument(payload *models.UpdateBlockPayload) bson.E {
	switch payload.Block.Type {
	case notesv1.Block_TYPE_HEADING_1.String():
		return bson.E{Key: "blocks.$.heading", Value: payload.Block.Heading}
	case notesv1.Block_TYPE_HEADING_2.String():
		return bson.E{Key: "blocks.$.heading", Value: payload.Block.Heading}
	case notesv1.Block_TYPE_HEADING_3.String():
		return bson.E{Key: "blocks.$.heading", Value: payload.Block.Heading}
	case notesv1.Block_TYPE_BULLET_POINT.String():
		return bson.E{Key: "blocks.$.bulletPoint", Value: payload.Block.BulletPoint}
	case notesv1.Block_TYPE_NUMBER_POINT.String():
		return bson.E{Key: "blocks.$.numberPoint", Value: payload.Block.NumberPoint}
	case notesv1.Block_TYPE_PARAGRAPH.String():
		return bson.E{Key: "blocks.$.paragraph", Value: payload.Block.Paragraph}
	case notesv1.Block_TYPE_MATH.String():
		return bson.E{Key: "blocks.$.math", Value: payload.Block.Math}
	case notesv1.Block_TYPE_IMAGE.String():
		return bson.E{Key: "blocks.$.image", Value: payload.Block.Image}
	case notesv1.Block_TYPE_CODE.String():
		return bson.E{Key: "blocks.$.code", Value: payload.Block.Code}
	}
	return bson.E{Key: "", Value: nil}
}
