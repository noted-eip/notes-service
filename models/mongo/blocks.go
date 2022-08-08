package mongo

import (
	"context"
	"notes-service/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type block struct {
	ID      string  `json:"id" bson:"_id,omitempty"`
	NoteId  string  `json:"noteId" bson:"noteId,omitempty"`
	Type    uint32  `json:"type" bson:"type,omitempty"`
	Content *string `json:"content" bson:"content,omitempty"`
}

type blocksRepository struct {
	logger          *zap.Logger
	db              *mongo.Database
	notesRepository models.NotesRepository
}

func NewBlocksRepository(db *mongo.Database, logger *zap.Logger, notesRepository models.NotesRepository) models.BlocksRepository {
	return &blocksRepository{
		logger:          logger,
		db:              db,
		notesRepository: notesRepository,
	}
}

func (srv *blocksRepository) Create(ctx context.Context, blockRequest *models.BlockWithIndex) error {
	id, err := uuid.NewRandom()

	if err != nil {
		srv.logger.Error("failed to generate new random uuid", zap.Error(err))
		return status.Errorf(codes.Internal, "could not create account")
	}

	//get note
	filter := models.NoteFilter{AuthorId: blockRequest.NoteId}
	currentNote, err := srv.notesRepository.Get(ctx, &filter)
	if err != nil {
		srv.logger.Error("mongo get note failed", zap.Error(err), zap.String("note Id", blockRequest.NoteId))
		return status.Errorf(codes.Internal, "could not create block")
	}

	//ajouter ce block a la note
	insertedBlock := models.Block{ID: id.String(), NoteId: blockRequest.NoteId, Type: blockRequest.Type, Content: blockRequest.Content}
	if (len(currentNote.Blocks) - 1) >= int(blockRequest.Index) {
		if blockRequest.Index == 0 {
			currentNote.Blocks = append([]*models.Block{&insertedBlock}, currentNote.Blocks...)
		} else {
			fstPartNote := currentNote.Blocks[0 : blockRequest.Index-1]
			secPartNote := currentNote.Blocks[blockRequest.Index:len(currentNote.Blocks)]
			fstPartNote = append(fstPartNote, &insertedBlock)
			currentNote.Blocks = append(fstPartNote, secPartNote...)
		}
	} else {
		// the index doesn't exist so we add it at the end
		currentNote.Blocks = append(currentNote.Blocks, &insertedBlock)
	}

	noteToReturn := note{ID: id.String(), AuthorId: currentNote.AuthorId, Title: currentNote.Title, Blocks: currentNote.Blocks}

	_, err = srv.db.Collection("notes").InsertOne(ctx, noteToReturn)
	if err != nil {
		srv.logger.Error("mongo insert note failed", zap.Error(err), zap.String("block name", blockRequest.NoteId))
		return status.Errorf(codes.Internal, "could not create block")
	}
	return nil
}

func (srv *blocksRepository) Update(ctx context.Context, filter *models.BlockFilter, blockRequest *models.BlockWithIndex) (*models.NoteWithBlocks, error) {
	return nil, nil
}

func (srv *blocksRepository) Delete(ctx context.Context, filter *models.BlockFilter) (*models.NoteWithBlocks, error) {

	/*objID, err := primitive.ObjectIDFromHex(filter.BlockId)
	if err != nil {
		return status.Errorf(codes.Internal, "Error : during get note 1")
	}*/

	var note note
	err := srv.db.Collection("notes").FindOne(ctx, buildBlockQuery(filter)).Decode(&note)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not find note")
	}

	var exist = true
	var currentIndex = 0
	for currentIndex = 0; currentIndex < len(note.Blocks); currentIndex++ {
		if note.Blocks[currentIndex].ID == filter.BlockId {
			exist = true
			break
		}
	}

	if !exist {
		return nil, status.Errorf(codes.Internal, "No such block with requested id")
	}

	//suprimer l'ancien endroit ou il etait
	fstPartNote := note.Blocks[0:currentIndex]
	secPartNote := note.Blocks[currentIndex+1 : len(note.Blocks)]
	note.Blocks = append(fstPartNote, secPartNote...)

	//update note by removing the block
	uuid, err := uuid.Parse(note.ID)
	return &models.NoteWithBlocks{ID: uuid, AuthorId: note.AuthorId, Title: note.Title, Blocks: note.Blocks}, nil
}

func buildBlockQuery(filter *models.BlockFilter) bson.M {
	query := bson.M{}
	if filter.BlockId != "" {
		query["_id"] = filter.BlockId
	}
	if filter.NoteId != "" {
		query["_id"] = filter.NoteId
	}
	return query
}
