package bot

import (
	"context"
	"fmt"
	"github.com/Ostmind/tgbot/internal/storage/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Use /login to authenticate.")
	bot.Send(msg)
}

func HandleLogin(bot *tgbotapi.BotAPI, update tgbotapi.Update, auth *AuthManager) {
	chatID := update.Message.Chat.ID
	code := auth.GenerateCode(chatID)
	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("Your 2FA code is: %s\nUse /verify <code> to authenticate.", code))
	bot.Send(msg)
}

func HandleVerify(bot *tgbotapi.BotAPI, update tgbotapi.Update, auth *AuthManager) {
	chatID := update.Message.Chat.ID
	args := update.Message.CommandArguments()
	if auth.VerifyCode(chatID, args) {
		msg := tgbotapi.NewMessage(chatID, "Authentication success! You can access commands now.")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Invalid or expired code. Use /login to get a new code.")
		bot.Send(msg)
	}
}

func HandleSecret(bot *tgbotapi.BotAPI, update tgbotapi.Update, auth *AuthManager) {
	chatID := update.Message.Chat.ID
	if !auth.IsAuthenticated(chatID) {
		msg := tgbotapi.NewMessage(chatID, "You need to authenticate first. Use /login.")
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(chatID, "Secret content only for authenticated users!")
	bot.Send(msg)
}

func HandleMyTodos(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *postgres.Storage) {
	chatID := update.Message.Chat.ID
	telegramUserID := fmt.Sprintf("%d", update.Message.From.ID)

	todos, err := db.GetUserTodos(context.Background(), telegramUserID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Ошибка при получении задач.")
		bot.Send(msg)
		return
	}

	if len(todos) == 0 {
		msg := tgbotapi.NewMessage(chatID, "У вас нет задач.")
		bot.Send(msg)
		return
	}

	var sb strings.Builder
	for i, todo := range todos {
		status := "❌"
		if todo.Completed == "true" || todo.Completed == "1" {
			status = "✅"
		}
		sb.WriteString(fmt.Sprintf("%d. %s - %s %s\n", i+1, todo.Title, todo.Description, status))
	}

	reply := sb.String()
	msg := tgbotapi.NewMessage(chatID, reply)
	bot.Send(msg)
}
