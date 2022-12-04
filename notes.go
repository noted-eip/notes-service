package main

import (
	"context"
	"fmt"
	"strconv"

	"notes-service/auth"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"
	"notes-service/validators"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notesService struct {
	notespb.UnimplementedNotesAPIServer

	auth      auth.Service
	logger    *zap.Logger
	repoNote  models.NotesRepository
	repoBlock models.BlocksRepository
}

var _ notespb.NotesAPIServer = &notesService{}

func (srv *notesService) CreateNote(ctx context.Context, in *notespb.CreateNoteRequest) (*notespb.CreateNoteResponse, error) {
	err := validators.ValidateCreateNoteRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	note, err := srv.repoNote.Create(ctx, &models.Note{AuthorId: in.Note.AuthorId, Title: in.Note.Title})

	if err != nil {
		srv.logger.Error("failed to create note", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not create note")
	}

	blocks := make([]models.Block, len(in.Note.Blocks))

	for index, block := range in.Note.Blocks {
		err := fillBlockContent(&blocks[index], block)
		if err != nil {
			srv.logger.Error("failed to create note", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "invalid content provided for block index : %d", index)
		}
		srv.repoBlock.Create(ctx, &models.Block{NoteId: note.ID, Type: uint32(in.Note.Blocks[index].Type), Index: uint32(index + 1), Content: blocks[index].Content})
	}

	return &notespb.CreateNoteResponse{
		Note: &notespb.Note{
			Id:       note.ID,
			AuthorId: note.AuthorId,
			Title:    note.Title,
			Blocks:   in.Note.Blocks,
		},
	}, nil
}

// QUE SI ELLE EST DANS TON GROUPE auth
func (srv *notesService) GetNote(ctx context.Context, in *notespb.GetNoteRequest) (*notespb.GetNoteResponse, error) {
	err := validators.ValidateGetNoteRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	note, err := srv.repoNote.Get(ctx, in.Id)
	if err != nil {
		srv.logger.Error("failed to get note", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "could not get note")
	}

	blocksTmp, err := srv.repoBlock.GetBlocks(ctx, note.ID)
	if err != nil {
		srv.logger.Error("failed to get blocks", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "invalid content provided for blocks form noteId : %d", note.ID)
	}

	//convert []models.block to []notespb.Block
	blocks := make([]*notespb.Block, len(blocksTmp))
	for index, block := range blocksTmp {
		blocks[index] = &notespb.Block{}
		err := fillContentFromModelToApi(block, block.Type, blocks[index])
		if err != nil {
			srv.logger.Error("failed to the content of a block", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "fail to get content from block Id : %s", block.ID)
		}
		blocks[index] = &notespb.Block{Id: block.ID, Type: notespb.Block_Type(block.Type), Data: blocks[index].Data}
	}
	noteToReturn := notespb.Note{Id: note.ID, AuthorId: note.AuthorId, Title: note.Title, Blocks: blocks}

	return &notespb.GetNoteResponse{Note: &noteToReturn}, nil
}

func (srv *notesService) UpdateNote(ctx context.Context, in *notespb.UpdateNoteRequest) (*notespb.UpdateNoteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateUpdateNoteRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//check si la note appartient a celui qui veut la modifier
	note, err := srv.repoNote.Get(ctx, in.Id)
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not update note")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	}

	//appeler deleteBlock avec le filtre note_id
	err = srv.repoBlock.DeleteBlocks(ctx, in.Id)
	if err != nil {
		srv.logger.Error("blocks weren't deleted : ", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not delete blocks")
	}

	//update juste les infos de la note et pas les blocks sinon
	noteId := in.Id
	err = srv.repoNote.Update(ctx, noteId, &models.Note{AuthorId: in.Note.AuthorId, Title: in.Note.Title})
	if err != nil {
		srv.logger.Error("failed to update note", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not update note")
	}

	//appeller createBlock en boucle pour tout les autres blocks
	blocks := make([]models.Block, len(in.Note.Blocks))
	for index, block := range in.Note.Blocks {
		err := fillBlockContent(&blocks[index], block)
		if err != nil {
			srv.logger.Error("failed to update blocks", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "invalid content provided for block index : %d", index)
		}
		srv.repoBlock.Create(ctx, &models.Block{NoteId: in.Id, Type: uint32(in.Note.Blocks[index].Type), Index: uint32(index + 1), Content: blocks[index].Content})
	}

	return nil, nil
}

func (srv *notesService) DeleteNote(ctx context.Context, in *notespb.DeleteNoteRequest) (*notespb.DeleteNoteResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateDeleteNoteRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//check si la note appartient a celui qui veut la modifier
	note, err := srv.repoNote.Get(ctx, in.Id)
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not update note")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	}

	//appeler deleteBlock avec le filtre note_id
	err = srv.repoBlock.DeleteBlocks(ctx, in.Id)
	if err != nil {
		srv.logger.Error("blocks weren't deleted : ", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not delete blocks")
	}

	err = srv.repoNote.Delete(ctx, in.Id)
	if err != nil {
		srv.logger.Error("failed to delete note", zap.Error(err))
		return nil, status.Error(codes.Internal, "could not delete note")
	}

	return nil, nil
}

// QUE SI ELLE EST DANS TON GROUPE auth
func (srv *notesService) ListNotes(ctx context.Context, in *notespb.ListNotesRequest) (*notespb.ListNotesResponse, error) {
	err := validators.ValidateListNoteRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	notes, err := srv.repoNote.List(ctx, in.AuthorId)
	if err != nil {
		srv.logger.Error("failed to get note", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "could not get note")
	}

	notesResponse := make([]*notespb.Note, len(*notes))
	for index, note := range *notes {
		notesResponse[index] = &notespb.Note{Id: note.ID, AuthorId: note.AuthorId, Title: note.Title}
	}

	return &notespb.ListNotesResponse{
		Notes: notesResponse,
	}, nil
}

func (srv *notesService) InsertBlock(ctx context.Context, in *notespb.InsertBlockRequest) (*notespb.InsertBlockResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateInsertBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//check si la note appartient a celui qui veut la modifier
	note, err := srv.repoNote.Get(ctx, strconv.Itoa(int(in.NoteId)))
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not update note")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	}

	var block = models.Block{}
	err = fillBlockContent(&block, in.Block)
	if err != nil {
		srv.logger.Error("failed to create block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid content provided for block index : %d", in.Index)
	}

	BlockId, err := srv.repoBlock.Create(ctx, &models.Block{NoteId: strconv.Itoa(int(in.NoteId)), Type: uint32(in.Block.Type), Index: in.Index, Content: block.Content})

	return &notespb.InsertBlockResponse{
		Block: &notespb.Block{
			Id:   *BlockId,
			Type: in.Block.Type,
			Data: in.Block.Data,
		},
	}, nil
}

func (srv *notesService) UpdateBlock(ctx context.Context, in *notespb.UpdateBlockRequest) (*notespb.UpdateBlockResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateUpdateBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//check si la note appartient a celui qui veut la modifier
	note, err := srv.repoNote.Get(ctx, in.Id)
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not update note")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	}

	var block = models.Block{}
	err = fillBlockContent(&block, in.Block)
	if err != nil {
		srv.logger.Error("failed to update block", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "invalid content provided for block id : %s", in.Id)
	}

	srv.repoBlock.Update(ctx, in.Id, &models.Block{ID: in.Id, Type: uint32(in.Block.Type), Index: in.Index, Content: block.Content})
	return &notespb.UpdateBlockResponse{
		Block: &notespb.Block{
			Id:   in.Id,
			Type: in.Block.Type,
			Data: in.Block.Data,
		},
	}, nil
}

func (srv *notesService) DeleteBlock(ctx context.Context, in *notespb.DeleteBlockRequest) (*notespb.DeleteBlockResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	err = validators.ValidateDeleteBlockRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//check si la note appartient a celui qui veut la modifier
	note, err := srv.repoNote.Get(ctx, in.Id)
	if err != nil {
		srv.logger.Error("Note not found in database", zap.Error(err))
		return nil, status.Error(codes.NotFound, "could not update note")
	}
	if token.UserID.String() != note.AuthorId {
		return nil, status.Error(codes.PermissionDenied, "This author has not the rights to update this note")
	}

	err = srv.repoBlock.DeleteBlock(ctx, in.Id)
	if err != nil {
		srv.logger.Error("block was not deleted : ", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not delete block")
	}

	return nil, nil
}

func fillBlockContent(block *models.Block, blockRequest *notespb.Block) error {
	switch op := blockRequest.Data.(type) {
	case *notespb.Block_Heading:
		block.Content = op.Heading
	case *notespb.Block_Paragraph:
		block.Content = op.Paragraph
	case *notespb.Block_NumberPoint:
		block.Content = op.NumberPoint
	case *notespb.Block_BulletPoint:
		block.Content = op.BulletPoint
	case *notespb.Block_Math:
		block.Content = op.Math
	//case *notespb.Block_Image_:
	//	block.Image.caption = &op.Image.Caption
	//	block.Image.url = &op.Image.Url
	//case *notespb.Block_Code_:
	//	block.Code.lang = &op.Code.Lang
	//	block.Code.Snippet = &op.Code.Snippet

	default:
		fmt.Println("No Data in this block")
		return status.Error(codes.Internal, "no data in this block")
	}
	return nil
}

func fillContentFromModelToApi(blockRequest *models.Block, contentType uint32, blockApi *notespb.Block) error {
	switch contentType {
	case 1:
		blockApi.Data = &notespb.Block_Heading{Heading: blockRequest.Content}
	case 2:
		blockApi.Data = &notespb.Block_Paragraph{Paragraph: blockRequest.Content}
	case 3:
		blockApi.Data = &notespb.Block_NumberPoint{NumberPoint: blockRequest.Content}
	case 4:
		blockApi.Data = &notespb.Block_BulletPoint{BulletPoint: blockRequest.Content}
	case 5:
		blockApi.Data = &notespb.Block_Math{Math: blockRequest.Content}
	//case 6:
	//	(*blockApi).Data = &notespb.Block_Image_{Image: {caption: blockRequest.Image.caption, url: blockRequest.Image.url}}
	//case 7:
	//	(*blockApi).Data = &notespb.Block_Code_{Code: {sinppet: blockRequest.Code.Snippet, lang: blockRequest.Code.Lang}}
	default:
		fmt.Println("No such content in this block")
		return status.Errorf(codes.Internal, "no such content in this block")
	}
	return nil
}

func (srv *notesService) authenticate(ctx context.Context) (*auth.Token, error) {
	token, err := srv.auth.TokenFromContext(ctx)
	if err != nil {
		srv.logger.Debug("failed to authenticate request", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return token, nil
}
