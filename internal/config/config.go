package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppID     int
	AppHash   string
	BotToken  string
	SearchID  int64
	TdlibDB   string
}

var cfg *Config

func init() {
	// try load .env if present
	_ = godotenv.Load()

	appIDStr := os.Getenv("APP_ID")
	if appIDStr == "" {
		log.Fatal("APP_ID is required")
	}
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		log.Fatalf("invalid APP_ID: %v", err)
	}

	appHash := os.Getenv("APP_HASH")
	if appHash == "" {
		log.Fatal("APP_HASH is required")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN is required")
	}

	searchStr := os.Getenv("SEARCH_CHAT_ID")
	if searchStr == "" {
		log.Fatal("SEARCH_CHAT_ID is required")
	}
	searchID, err := strconv.ParseInt(searchStr, 10, 64)
	if err != nil {
		log.Fatalf("invalid SEARCH_CHAT_ID: %v", err)
	}

	tdlibDB := os.Getenv("TDLIB_DB_DIR")
	if tdlibDB == "" {
		tdlibDB = "/data/tdlib-db"
	}

	cfg = &Config{
		AppID:    appID,
		AppHash:  appHash,
		BotToken: botToken,
		SearchID: searchID,
		TdlibDB:  tdlibDB,
	}
}

func Get() *Config {
	return cfg
}
