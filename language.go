package main

// func (srv *notesService) generateNoteTagsToModelNote(note *models.Note) error {
// 	var fullNote string

// 	// for _, block := range note.Blocks {
// 	// 	if block.Type != uint32(notespb.Block_TYPE_CODE) && block.Type != uint32(notespb.Block_TYPE_IMAGE) {
// 	// 		fullNote += block.Content + "\n"
// 	// 	}
// 	// }

// 	keywords, err := srv.language.GetKeywordsFromTextInput(fullNote)
// 	if err != nil {
// 		return status.Error(codes.Internal, err.Error())
// 	}

// 	note.Keywords = keywords
// 	return nil
// }