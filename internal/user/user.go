package user

import (
	"log"
	"time"

	tdlib "github.com/zelenin/go-tdlib/client"

	"github.com/the-lx0980/wroxen-go/internal/config"
)

// User wraps a TDLib client
type User struct {
	Client *tdlib.Client
	cfg    *config.Config
}

func NewUser(cfg *config.Config) *User {
	tdlib.SetLogVerbosityLevel(1)

	authorizer := tdlib.ClientAuthorizer()

	client := tdlib.NewClient(tdlib.Config{
		APIID:              cfg.AppID,
		APIHash:            cfg.AppHash,
		SystemLanguageCode: "en",
		DeviceModel:        "Server",
		ApplicationVersion: "1.0",
		UseTestDC:          false,
		DatabaseDirectory:  cfg.TdlibDB,
		FileDirectory:      cfg.TdlibDB + "/files",
		IgnoreFileNames:    false,
		Authorizer:         authorizer,
	})

	// give some time to init
	time.Sleep(300 * time.Millisecond)
	log.Println("TDLib user client created")
	return &User{Client: client, cfg: cfg}
}

// Start returns the user id if present (nil error on success)
func (u *User) Start() (int64, error) {
	me := u.Client.GetMe()
	if me == nil {
		return 0, nil
	}
	log.Printf("User client logged in as: %s (id=%d)\n", me.Username, me.UserId)
	return int64(me.UserId), nil
}

func (u *User) Stop() {
	u.Client.Close()
	log.Println("User client stopped.")
}

// SearchMessages searches messages in specified chat using TDLib's SearchChatMessagesRequest
func (u *User) SearchMessages(chatID int64, query string, limit int32) ([]tdlib.Message, error) {
	req := &tdlib.SearchChatMessagesRequest{
		ChatId: chatID,
		Query:  query,
		Limit:  limit,
	}
	res, err := u.Client.SearchChatMessages(req)
	if err != nil {
		return nil, err
	}
	var msgs []tdlib.Message
	for _, m := range res.Messages {
		msgs = append(msgs, m)
	}
	return msgs, nil
}
