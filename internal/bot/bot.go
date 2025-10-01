package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/the-lx0980/wroxen-go/internal/config"
	"github.com/the-lx0980/wroxen-go/internal/user"
)

type Wroxen struct {
	bot  *tgbotapi.BotAPI
	User *user.User
	cfg  *config.Config
}

func NewWroxen(cfg *config.Config, u *user.User) (*Wroxen, error) {
	b, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	b.Debug = false
	log.Printf("Authorized on account %s", b.Self.UserName)

	return &Wroxen{
		bot:  b,
		User: u,
		cfg:  cfg,
	}, nil
}

func (w *Wroxen) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := w.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go w.handleMessage(update.Message)
		}
		if update.CallbackQuery != nil {
			go w.handleCallback(update.CallbackQuery)
		}
	}
}

func (w *Wroxen) Stop() {
	w.User.Stop()
	log.Println("Bot stopped. Bye.")
}

// small helpers for building messages (same as Python)
func generateResultMessage(query string, movies []string, page int) string {
	start := (page-1)*10 + 1
	text := "Here are the results for <b>" + query + "</b>:\n\n"
	for i, m := range movies {
		text += strconv.Itoa(start+i) + ". " + m + "\n\n"
	}
	return text
}

func generateInlineKeyboardData(query string, totalResults, currentPage int) (bool, string, bool, string) {
	hasPrev := false
	prev := ""
	hasNext := false
	next := ""
	if currentPage > 1 {
		hasPrev = true
		prev = "previous_page:" + query + ":" + strconv.Itoa(currentPage-1)
	}
	if totalResults > currentPage*10 {
		hasNext = true
		next = "next_page:" + query + ":" + strconv.Itoa(currentPage+1)
	}
	return hasPrev, prev, hasNext, next
}
