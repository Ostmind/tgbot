package bot

import (
	"log"
	"os"

	"github.com/Ostmind/tgbot/internal/bot"
	"github.com/Ostmind/tgbot/internal/config"
	"github.com/Ostmind/tgbot/internal/storage/postgres"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		panic("TELEGRAM_TOKEN not set")
	}

	cfg := config.MustNew()

	db, err := postgres.New(cfg.DB)
	if err != nil {
		log.Fatalf("couldn't establish db connection %s", err)
	}

	bot.RunBot(token, db)
}
