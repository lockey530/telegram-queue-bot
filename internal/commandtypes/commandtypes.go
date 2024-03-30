package commandtypes

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/josh1248/nusc-queue-bot/internal/handlers"
)

type AcceptedCommands struct {
	Command     string
	Description string
	Handler     func(userMessage tgbotapi.Update) (feedback string)
}

var AvailableCommands = []AcceptedCommands{
	{
		Command:     "join",
		Description: "Join the virtual queue for the photobooth.",
		Handler:     handlers.JoinCommand,
	},
	{
		Command:     "leave",
		Description: "Leave the virtual queue for the photobooth.",
		Handler:     handlers.LeaveCommand,
	},
	{
		Command:     "howlong",
		Description: "Returns the expected time to wait in the queue",
		Handler:     handlers.HowLongCommand,
	},
	{
		Command:     "help",
		Description: "Explains the main functionalities of the bot.",
		Handler:     handlers.HelpCommand,
	},
	{
		Command:     "greet",
		Description: "The bot is friendly :)",
		Handler:     handlers.GreetCommand,
	},
	{
		Command:     "start",
		Description: "Explains the main functionalities of the bot.",
		Handler:     handlers.HelpCommand,
	},
}
