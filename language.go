package main

import (
	"notes-service/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *notesAPI) generateNoteTagsToModelNote(note *models.Note) error {
	var fullNote string

	for _, block := range note.Blocks {
		content, ok := getContent(&block)
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

func getContent(block *models.NoteBlock) (string, bool) {
	switch block.Type {
	case "heading":
		return *block.Heading, true
	case "paragraph":
		return *block.Paragraph, true
	case "math":
		return *block.Math, true
	case "bulletpoint":
		return *block.BulletPoint, true
	case "numberpoint":
		return *block.NumberPoint, true
	default:
		return "", false
	}
}
