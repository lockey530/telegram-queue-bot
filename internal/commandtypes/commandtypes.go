package commandtypes

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// Represents the valid commands that the bot handles
type BotCommand string

// Feedback given to the user after a command is ran
type BotFeedback string

type BotCommandHandler struct {
	Command     BotCommand
	Description string
	Feedback    func(update tgbotapi.Update) BotFeedback
}
