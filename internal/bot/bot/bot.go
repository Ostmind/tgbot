package bot

import (
	"log"

	"github.com/Ostmind/tgbot/internal/bot/auth"
	"github.com/Ostmind/tgbot/internal/bot/comands"
	"github.com/Ostmind/tgbot/internal/bot/config"
	"github.com/Ostmind/tgbot/internal/storage/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RunBot(cfg config.BotConfig, db *postgres.Storage) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(cfg.BotOffset)
	u.Timeout = cfg.BotTimeout
	updates := bot.GetUpdatesChan(u)

	authManager := auth.NewAuthManager()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "start":
			comands.HandleStart(bot, update)
		case "login":
			comands.HandleLogin(bot, update, authManager)
		case "verify":
			comands.HandleVerify(bot, update, authManager)
		case "secret":
			comands.HandleSecret(bot, update, authManager)
		case "mytodos":
			comands.HandleMyTodos(bot, update, db)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Use /start, /login, /verify, /secret.")
			bot.Send(msg)
		}
	}
}
