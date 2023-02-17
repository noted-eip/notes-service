package main

import (
	"notes-service/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *notesAPI) generateNoteTagsToModelNote(note *models.Note) error {
	var fullNote string

	for _, block := range note.Blocks {
		content, ok := GetBlockContent(&block)
		if ok {
			fullNote += content + "\n"
		}
	}

	keywords, err := srv.language.GetKeywordsFromTextInput(fullNote)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	note.Keywords = keywords
	return nil
}
