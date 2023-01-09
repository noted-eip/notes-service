package exports

import (
	"bytes"
	notespb "notes-service/protorepo/noted/notes/v1"

	"github.com/yuin/goldmark"

	pdf "github.com/stephenafamo/goldmark-pdf"
)

func NoteToPDF(n *notespb.Note) ([]byte, error) {
	markdownNote, err := NoteToMarkdown(n)
	var result bytes.Buffer

	if err != nil {
		return nil, err
	}

	// NOTE: Can be customized later on
	converter := goldmark.New(
		goldmark.WithRenderer(pdf.New()),
	)

	converter.Convert(markdownNote, &result)

	return result.Bytes(), nil
}
