package bot

import (
	"log"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/the-lx0980/wroxen-go/internal/search"
	"github.com/the-lx0980/wroxen-go/internal/storage"
	"github.com/the-lx0980/wroxen-go/internal/user"
)

var commandOrEmojiRe = regexp.MustCompile(`((^\/|^,|^!|^\.|^[\U0001F600-\U000E007F]).*)`)

func (w *Wroxen) handleMessage(msg *tgbotapi.Message) {
	if msg == nil || msg.Text == "" {
		return
	}
	if !msg.Chat.IsGroup() && !msg.Chat.IsSuperGroup() {
		return
	}
	if commandOrEmojiRe.MatchString(msg.Text) {
		return
	}
	if len(msg.Text) <= 2 {
		return
	}

	query := msg.Text

	// call user client to search messages in SEARCH_CHAT
	tdmsgs, err := w.User.SearchMessages(w.cfg.SearchID, query, 50)
	if err != nil {
		log.Println("SearchMessages error:", err)
		return
	}
	if len(tdmsgs) == 0 {
		// optionally send "no results"
		return
	}

	entries := search.FormatMessagesToEntries(tdmsgs)
	page := 1
	w.sendResultMessage(msg, query, entries, page, int64(msg.From.ID))
}

func (w *Wroxen) sendResultMessage(fromMessage *tgbotapi.Message, query string, movies []string, page int, requesterID int64) {
	total := len(movies)

	// pagination
	var moviesPage []string
	if total <= 10 {
		moviesPage = movies
	} else {
		start := (page - 1) * 10
		end := page * 10
		if start < 0 {
			start = 0
		}
		if end > total {
			end = total
		}
		moviesPage = movies[start:end]
	}

	text := generateResultMessage(query, moviesPage, page)

	// build keyboard
	hasPrev, prevData, hasNext, nextData := generateInlineKeyboardData(query, total, page)
	var row []tgbotapi.InlineKeyboardButton
	if hasPrev {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Previous Page", prevData))
	}
	if hasNext {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Next Page", nextData))
	}
	var markup tgbotapi.InlineKeyboardMarkup
	if len(row) > 0 {
		markup = tgbotapi.NewInlineKeyboardMarkup(row)
	}

	msg := tgbotapi.NewMessage(fromMessage.Chat.ID, text)
	msg.ParseMode = "HTML"
	if len(row) > 0 {
		msg.ReplyMarkup = markup
	}

	sent, err := w.bot.Send(msg)
	if err != nil {
		log.Println("send message error:", err)
		return
	}

	// save in cache
	storage.Save(query, storage.DBEntry{
		MessageID:   sent.MessageID,
		Movies:      movies,
		Page:        page,
		RequesterID: requesterID,
	})
}

func (w *Wroxen) handleCallback(cq *tgbotapi.CallbackQuery) {
	if cq == nil || cq.Data == "" {
		return
	}
	data := cq.Data
	if !(startsWith(data, "next_page:") || startsWith(data, "previous_page:")) {
		return
	}

	parts := splitN(data, ":", 3)
	if len(parts) != 3 {
		return
	}
	query := parts[1]
	pageStr := parts[2]
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return
	}

	entry, ok := storage.Get(query)
	if !ok {
		_ = w.bot.Request(tgbotapi.NewCallback(cq.ID, "यह मैसेज काफी पुराना हो चुका है।"))
		return
	}

	// restrict to requester
	if entry.RequesterID != 0 && cq.From.ID != int(entry.RequesterID) {
		_, _ = w.bot.Request(tgbotapi.NewCallbackWithAlert(cq.ID, "यह आपके लिए नही है!"))
		return
	}

	// prepare page
	total := len(entry.Movies)
	start := (page - 1) * 10
	end := page * 10
	if start < 0 {
		start = 0
	}
	if end > total {
		end = total
	}
	moviesPage := entry.Movies[start:end]

	text := generateResultMessage(query, moviesPage, page)
	hasPrev, prevData, hasNext, nextData := generateInlineKeyboardData(query, total, page)
	var row []tgbotapi.InlineKeyboardButton
	if hasPrev {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Previous Page", prevData))
	}
	if hasNext {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("Next Page", nextData))
	}
	var markup tgbotapi.InlineKeyboardMarkup
	if len(row) > 0 {
		markup = tgbotapi.NewInlineKeyboardMarkup(row)
	}

	edit := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, text)
	edit.ParseMode = "HTML"
	if len(row) > 0 {
		edit.ReplyMarkup = &markup
	}

	_, err = w.bot.Send(edit)
	if err != nil {
		log.Println("callback edit error:", err)
	}

	// update DB page
	storage.Save(query, storage.DBEntry{
		MessageID:   cq.Message.MessageID,
		Movies:      entry.Movies,
		Page:        page,
		RequesterID: entry.RequesterID,
	})

	_, _ = w.bot.Request(tgbotapi.NewCallback(cq.ID, ""))
}

// small helpers copied from earlier code
func startsWith(s, pref string) bool {
	return len(s) >= len(pref) && s[:len(pref)] == pref
}

func splitN(s, sep string, n int) []string {
	parts := make([]string, 0, n)
	i := 0
	for i < n-1 {
		j := indexOfNth(s, sep, 1)
		if j < 0 {
			break
		}
		parts = append(parts, s[:j])
		s = s[j+1:]
		i++
	}
	parts = append(parts, s)
	return parts
}

func indexOfNth(s, sep string, nth int) int {
	if nth <= 0 {
		return -1
	}
	idx := -1
	pos := 0
	for i := 0; i < nth; i++ {
		k := -1
		for t := pos; t+len(sep) <= len(s); t++ {
			if s[t:t+len(sep)] == sep {
				k = t
				break
			}
		}
		if k == -1 {
			return -1
		}
		idx = k
		pos = k + len(sep)
	}
	return idx
}
