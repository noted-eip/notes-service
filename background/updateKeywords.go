package background

import (
	"context"
	"notes-service/models"

	notespb "notes-service/protorepo/noted/notes/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *Service) UpdateKeywordsByNoteId(noteId string) error {
	//get la note & les blocks
	note, err := srv.repoNote.Get(context.TODO(), noteId)
	if err != nil {
		srv.logger.Error("failed to get note", zap.Error(err))
		return status.Error(codes.NotFound, "could not get note.")
	}
	blocks, err := srv.repoBlock.GetBlocks(context.TODO(), note.ID)
	if err != nil {
		srv.logger.Error("failed to get blocks", zap.Error(err))
		return status.Errorf(codes.NotFound, "invalid content provided for blocks form noteId : %s", note.ID)
	}
	for index, block := range blocks {
		newSize := len(note.Blocks) + 1
		note.Blocks = make([]models.Block, newSize)
		note.Blocks[index] = *block
	}
	//gen les keywords
	err = srv.generateNoteTagsToModelNote(note)
	if err != nil {
		srv.logger.Error("failed to gen keywords", zap.Error(err))
		return status.Errorf(codes.Internal, "failed to gen keywords for noteId : %s", note.ID)
	}

	//update la note
	newNote := models.NotePayload{ID: note.ID, AuthorId: note.AuthorId, Title: note.Title, Blocks: note.Blocks, Keywords: note.Keywords}
	err = srv.repoNote.Update(context.TODO(), noteId, &newNote)
	if err != nil {
		srv.logger.Error("failed upate note with keywords", zap.Error(err))
		return status.Errorf(codes.Internal, "failed upate note with keywords for noteId : %s", note.ID)
	}

	/*note, err = srv.repoNote.Get(context.TODO(), noteId)
	if err != nil {
		srv.logger.Error("failed to get note", zap.Error(err))
		return status.Error(codes.NotFound, "could not get note.")
	}*/
	return nil
}

func (srv *Service) generateNoteTagsToModelNote(note *models.Note) error {
	var fullNote string

	for _, block := range note.Blocks {
		if block.Type != uint32(notespb.Block_TYPE_CODE) && block.Type != uint32(notespb.Block_TYPE_IMAGE) {
			fullNote += block.Content + "\n"
		}
	}

	keywords, err := srv.language.GetKeywordsFromTextInput(fullNote)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	note.Keywords = *keywords
	return nil
}
