package helper

import (
	"strings"

	"github.com/jomei/notionapi"
)

func RichTextsToString(richTexts []notionapi.RichText) string {
	var sb strings.Builder
	for _, rt := range richTexts {
		sb.WriteString(rt.Text.Content)
	}
	return sb.String()
}
