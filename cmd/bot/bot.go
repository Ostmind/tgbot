package main

import (
	"log"

	"github.com/Ostmind/tgbot/internal/bot/bot"
	"github.com/Ostmind/tgbot/internal/bot/config"
	"github.com/Ostmind/tgbot/internal/storage/postgres"
	mainconfig "github.com/Ostmind/tgbot/internal/todo/config"
)

func main() {

	cfg := config.MustNew()

	db, err := postgres.New(mainconfig.DatabaseConfig(cfg.DB))
	if err != nil {
		log.Fatalf("couldn't establish db connection %s", err)
	}

	bot.RunBot(cfg.Bot, db)
}
