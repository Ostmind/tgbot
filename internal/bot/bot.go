package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RunBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	authManager := NewAuthManager()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "start":
			HandleStart(bot, update, authManager)
		case "login":
			HandleLogin(bot, update, authManager)
		case "verify":
			HandleVerify(bot, update, authManager)
		case "secret":
			HandleSecret(bot, update, authManager)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Use /start, /login, /verify, /secret.")
			bot.Send(msg)
		}
	}
}
