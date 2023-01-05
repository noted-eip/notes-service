package exports

import (
	notespb "notes-service/protorepo/noted/notes/v1"
	"strconv"
	"strings"
)

func headingBlockToMarkdown(b *notespb.Block) (string, error) {
	heading := b.GetHeading()
	typeName := b.Type.String()
	typeNameSplitted := strings.Split(typeName, "_")
	print(typeName, "\n")
	importance, err := strconv.Atoi(typeNameSplitted[len(typeNameSplitted)-1])

	if err != nil {
		return "", err
	}

	result := strings.Repeat("#", importance) + " " + heading

	sanitizeNewLines(&result)

	return result, nil
}

func codeBlockToMarkdown(b *notespb.Block) string {
	codeData := b.GetCode()
	result := "```" + codeData.Lang + "\n" + codeData.Snippet
	sanitizeNewLines(&result)
	result = result + "```\n"
	return result
}

func imageBlockToMarkdown(b *notespb.Block) string {
	imageData := b.GetImage()
	return "![](" + imageData.Url + " " + imageData.Caption + ")\n"
}

func sanitizeNewLines(str *string) {
	*str = strings.ReplaceAll(*str, "\r\n", "\n")
	if (*str)[len(*str)-1] != '\n' {
		*str = *str + "\n"
	}
}

func NoteToMarkdown(n *notespb.Note) ([]byte, error) {
	result := ""
	var err error = nil
	var converted string = ""

	for _, block := range n.Blocks {
		switch op := block.Data.(type) {
		case *notespb.Block_Heading:
			converted, err = headingBlockToMarkdown(block)
		case *notespb.Block_Paragraph: // NOTE: Is already formatted as Markdown.
			sanitizeNewLines(&op.Paragraph)
			converted = op.Paragraph
		case *notespb.Block_NumberPoint: // NOTE: Is already formatted as Markdown ?
			sanitizeNewLines(&op.NumberPoint)
			converted = op.NumberPoint + "\n"
		case *notespb.Block_BulletPoint: // NOTE: Is already formatted as Markdown ?
			sanitizeNewLines(&op.BulletPoint)
			converted = op.BulletPoint + "\n"
		case *notespb.Block_Math:
			sanitizeNewLines(&op.Math)
			converted = op.Math
		case *notespb.Block_Code_:
			converted = codeBlockToMarkdown(block)
		case *notespb.Block_Image_:
			converted = imageBlockToMarkdown(block)
		}
		if err != nil {
			return nil, err
		}
		result = result + converted
	}

	return []byte(result), nil
}
