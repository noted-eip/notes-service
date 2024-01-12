package mongo

import (
	"context"
	"errors"
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
	blocks := &[]models.NoteBlock{}

	if len(payload.Blocks) > 0 {
		// Create "real" empty array for mongodb golang drivers
		for i := 0; i < len(payload.Blocks); i++ {
			(payload.Blocks[i]).Thread = &[]models.BlockComment{}
			(payload.Blocks[i]).Styles = &[]models.TextStyle{}
		}
		blocks = &payload.Blocks
	} else {
		// @note: fill an empty block if none was provided
		content := ""
		*blocks = append((*blocks), models.NoteBlock{
			ID:        repo.newUUID(),
			Type:      "TYPE_PARAGRAPH",
			Paragraph: &content,
			Thread:    &[]models.BlockComment{},
			Styles:    &[]models.TextStyle{},
		})
	}

	note := &models.Note{
		ID:                          repo.newUUID(),
		Title:                       payload.Title,
		AuthorAccountID:             accountID,
		GroupID:                     payload.GroupID,
		CreatedAt:                   now,
		ModifiedAt:                  nil,
		AnalyzedAt:                  nil,
		Keywords:                    []*models.Keyword{},
		Blocks:                      blocks,
		AccountsWithEditPermissions: []string{accountID},
		Quizs:                       &[]models.Quiz{},
		Lang:                        payload.Lang,
	}

	err := repo.insertOne(ctx, &note)
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

	if payload.Blocks != nil {
		for i := range *payload.Blocks {
			// @note: if the id is invalid or wrong format, it's a new block, then we create a new ID
			if len((*payload.Blocks)[i].ID) < 21 {
				(*payload.Blocks)[i].ID = repo.newUUID()
			} else {
				// @note: Reset the "thread" field to its current value to avoid modification
				block, err := repo.GetBlock(ctx, &models.OneBlockFilter{
					GroupID: filter.GroupID,
					NoteID:  filter.NoteID,
					BlockID: (*payload.Blocks)[i].ID,
				}, accountID)
				if err == nil {
					(*payload.Blocks)[i].Thread = block.Thread
					(*payload.Blocks)[i].Styles = block.Styles
				} else {
					repo.logger.Error("error while getting block in UpdateNote", zap.Error(err))
				}
			}
		}
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
	requiredFields := bson.D{{Key: "blocks", Value: 0}, {Key: "keywords", Value: 0}, {Key: "quizs", Value: 0}}
	opts := options.Find().SetProjection(requiredFields)

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
	payload.Block.ID = repo.newUUID()

	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		{Key: "authorAccountId", Value: accountID},
	}

	payload.Block.Thread = &[]models.BlockComment{} // Make non-null empty array
	payload.Block.Styles = &[]models.TextStyle{}    // Make non-null empty array

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "blocks", Value: bson.D{
				{Key: "$each", Value: bson.A{payload.Block}},
				{Key: "$position", Value: payload.Index},
			}},
		}},
	}

	err := repo.updateOne(ctx, query, update)
	if err != nil {
		return nil, err
	}

	block, err := repo.GetBlock(ctx,
		&models.OneBlockFilter{
			GroupID: filter.GroupID,
			NoteID:  filter.NoteID,
			BlockID: payload.Block.ID,
		},
		accountID,
	)
	if err != nil {
		return nil, err
	}

	return block, nil
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

// If filter is set to nil, every edit permissions of his will be deleted on the db
// If filter is not set to nil, GroupID is mandatory. To specify one note, fill NoteID
func (repo *notesRepository) RemoveEditPermissions(ctx context.Context, filter *models.OneNoteFilter, accountID string) error {
	query := bson.D{
		{Key: "accountsWithEditPermissions", Value: accountID}, // NOTE: MongoDB model logic is not "safe" here - Who/What can call this function is decided in the endpoint's logic
	}

	if filter != nil {
		if filter.GroupID != "" {
			query = append(query, bson.E{Key: "groupId", Value: filter.GroupID})
		} else {
			return errors.New("when removing edit permissions with a filter, please specify a group id")
		}
		if filter.NoteID != "" {
			query = append(query, bson.E{Key: "_id", Value: filter.NoteID})
		}
	}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "accountsWithEditPermissions", Value: accountID},
		}},
	}

	_, err := repo.updateMany(ctx, query, update)
	return err
}

func (repo *notesRepository) CreateBlockComment(ctx context.Context, filter *models.OneBlockFilter, payload *models.BlockComment, accountID string) (*models.BlockComment, error) {
	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		{Key: "blocks.id", Value: filter.BlockID},
		// {Key: "$or", Value: bson.A{
		// 	bson.D{{Key: "accountsWithEditPermissions", Value: accountID}},
		// 	bson.D{{Key: "authorAccountId", Value: accountID}},
		// }},
	}

	commentID := repo.newUUID()
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "blocks.$.thread", Value: models.BlockComment{
				ID:              commentID,
				AuthorAccountID: payload.AuthorAccountID,
				Content:         payload.Content,
			}},
		}},
	}

	err := repo.updateOne(ctx, query, update)
	if err != nil {
		return nil, err
	}

	res, err := repo.GetBlock(ctx, &models.OneBlockFilter{
		GroupID: filter.GroupID,
		NoteID:  filter.NoteID,
		BlockID: filter.BlockID,
	}, accountID)

	return res.FindComment(commentID), err
}

func (repo *notesRepository) ListBlockComments(ctx context.Context, filter *models.OneBlockFilter, lo *models.ListOptions, accountID string) (*[]models.BlockComment, error) {
	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		{Key: "blocks.id", Value: filter.BlockID},
		// {Key: "$or", Value: bson.A{
		// 	bson.D{{Key: "accountsWithEditPermissions", Value: accountID}},
		// 	bson.D{{Key: "authorAccountId", Value: accountID}},
		// }},
	}

	requiredFields := bson.D{{Key: "blocks", Value: 1}}
	opts := options.FindOne().SetProjection(requiredFields)

	note := models.Note{}
	err := repo.findOne(ctx, query, &note, opts)
	if err != nil {
		return nil, err
	}

	return note.FindBlock(filter.BlockID).Thread, nil
}

func (repo *notesRepository) DeleteBlockComment(ctx context.Context, filter *models.OneBlockFilter, payload *models.BlockComment, accountID string) (*models.BlockComment, error) {
	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		{Key: "blocks.id", Value: filter.BlockID},
		// {Key: "$or", Value: bson.A{
		// 	bson.D{{Key: "accountsWithEditPermissions", Value: accountID}},
		// 	bson.D{{Key: "authorAccountId", Value: accountID}},
		// }},
	}

	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "blocks.$.thread", Value: bson.D{
				{Key: "id", Value: payload.ID},
				{Key: "authorAccountId", Value: accountID},
			}},
		}},
	}

	err := repo.updateOne(ctx, query, update)
	if err != nil {
		return nil, err
	}
	return payload, err
}

func (repo *notesRepository) StoreNewQuiz(ctx context.Context, filter *models.OneNoteFilter, payload *models.Quiz, accountID string) (*models.Quiz, error) {
	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		{Key: "$or", Value: bson.A{
			// bson.D{{Key: "accountsWithEditPermissions", Value: accountID}},
			bson.D{{Key: "authorAccountId", Value: accountID}},
		}},
	}

	payload.ID = repo.newUUID()
	payload.CreatedAt = time.Now()

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "quizs", Value: payload},
		}},
	}

	err := repo.updateOne(ctx, query, update)
	if err != nil {
		return nil, err
	}
	return payload, err
}

func (repo *notesRepository) ListQuizs(ctx context.Context, filter *models.OneNoteFilter, accountID string) (*[]models.Quiz, error) {
	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		// Right checks are made before, we can't access members from here
	}

	requiredFields := bson.D{{Key: "quizs", Value: 1}}
	opts := options.FindOne().SetProjection(requiredFields)

	note := models.Note{}
	err := repo.findOne(ctx, query, &note, opts)
	if err != nil {
		return nil, err
	}

	return note.Quizs, nil
}

// This function is used to put an expiration date on all Quizs after a server reboot (background services)
func (repo *notesRepository) ListQuizsCreatedDateInternal(ctx context.Context) (*[]models.Quiz, error) {
	unwind := bson.D{{
		Key: "$unwind", Value: "$quizs",
	}}

	projectIDandDate := bson.D{{
		Key: "$project", Value: bson.D{
			{Key: "id", Value: "$quizs.id"},
			{Key: "createdAt", Value: "$quizs.createdAt"},
		},
	}}

	res := &[]models.Quiz{}
	err := repo.aggregate(ctx, mongo.Pipeline{unwind, projectIDandDate}, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *notesRepository) DeleteQuiz(ctx context.Context, filter *models.OneNoteFilter, quizID string, accountID string) error {
	query := bson.D{
		{Key: "_id", Value: filter.NoteID},
		{Key: "groupId", Value: filter.GroupID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "accountsWithEditPermissions", Value: accountID}},
			bson.D{{Key: "authorAccountId", Value: accountID}},
		}},
	}

	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "quizs", Value: bson.D{
				{Key: "id", Value: quizID},
			}},
		}},
	}

	return repo.updateOne(ctx, query, update)
}

func (repo *notesRepository) DeleteQuizFromIDInternal(ctx context.Context, quizID string) error {
	query := bson.D{}

	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "quizs", Value: bson.D{
				{Key: "id", Value: quizID},
			}},
		}},
	}

	_, err := repo.updateMany(ctx, query, update)
	return err
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
