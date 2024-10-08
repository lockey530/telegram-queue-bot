package types

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AcceptedCommands struct {
	Command     string
	Description string
	Handler     func(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string)
}
