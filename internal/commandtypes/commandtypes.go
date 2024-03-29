package commandtypes

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type CommandHandler struct {
	Command     string
	Description string
	// conducts action associated with function and returns a response providing feedback to the user
	Action func(update tgbotapi.Update) string
}
