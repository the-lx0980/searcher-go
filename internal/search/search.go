package search

import (
	"fmt"

	tdlib "github.com/zelenin/go-tdlib/client"
)

// Convert TDLib Message to formatted movie entry
func FormatMessagesToEntries(msgs []tdlib.Message) []string {
	out := make([]string, 0, len(msgs))
	for _, m := range msgs {
		var caption string
		// Extract caption/text from different message types
		switch content := m.Content.(type) {
		case *tdlib.MessageText:
			caption = content.Text.Text
		case *tdlib.MessagePhoto:
			caption = content.Caption.Text
		case *tdlib.MessageVideo:
			caption = content.Caption.Text
		case *tdlib.MessageDocument:
			caption = content.Caption.Text
		default:
			// fallback: empty
			caption = ""
		}

		// generate link (for supergroups/public groups this format often works)
		chatID := m.ChatId
		messageID := m.Id
		link := fmt.Sprintf("https://t.me/c/%d/%d", ToPublicChatID(chatID), messageID)

		name, year, q := ExtractMovieDetails(caption)
		entry := fmt.Sprintf("<b>%s (%s) %s</b>\n<b>Link:</b> %s", name, year, q, link)
		out = append(out, entry)
	}
	return out
}

func ToPublicChatID(chatID int64) int64 {
	if chatID < 0 {
		s := fmt.Sprintf("%d", chatID)
		s = stringsTrimPrefix(s, "-100")
		if s == "" {
			return chatID
		}
		var v int64
		fmt.Sscanf(s, "%d", &v)
		if v == 0 {
			return chatID
		}
		return v
	}
	return chatID
}

// small helper to trim prefix (avoid importing strings multiple times)
func stringsTrimPrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
