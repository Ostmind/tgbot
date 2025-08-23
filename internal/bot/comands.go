package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update, auth *AuthManager) {
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
