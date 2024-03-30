package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/josh1248/nusc-queue-bot/internal/commandtypes"
)

// Command and description is hard-coded within the HelpFunction for circular dependencies.
// You will need to update the help function accordingly if you want to change this definition.
// If this becomes too troublesome, consider using the reflect package to store the handler
// function names under a string first before using reflect.ValueOf().Call()
var AvailableCommands = []commandtypes.AcceptedCommands{
	{
		Command:     "join",
		Description: "Join the virtual queue for the photobooth.",
		Handler:     JoinCommand,
	},
	{
		Command:     "leave",
		Description: "Leave the virtual queue for the photobooth.",
		Handler:     LeaveCommand,
	},
	{
		Command:     "howlong",
		Description: "Returns the expected time to wait in the queue",
		Handler:     HowLongCommand,
	},
	{
		Command:     "help",
		Description: "Explains the main functionalities of the bot.",
		Handler:     HelpCommand,
	},
	{
		Command:     "greet",
		Description: "The bot is friendly :)",
		Handler:     GreetCommand,
	},
	{
		Command:     "start",
		Description: "Explains the main functionalities of the bot.",
		Handler:     HelpCommand,
	},
}

func NonTextHandler(userMessage tgbotapi.Update) (feedback string) {
	feedback = "I don't know what this is :( please send me text commands!"
	return feedback
}

func NonCommandHandler(userMessage tgbotapi.Update) (feedback string) {
	feedback = "Please input a command which starts with '/', like /start"
	return feedback
}

func InvalidCommand(userMessage tgbotapi.Update) (feedback string) {
	feedback = "Sorry, I don't recognize your command :("
	return feedback
}

func HelpCommand(userMessage tgbotapi.Update) (feedback string) {
	return `
	Welcome to the queue bot~

	/join - join the photobooth queue!

	/leave - leave the photobooth queue if you have previously joined.

	(dk) /wait - need some time? place yourself 5 slots behind the queue (1-time only).
	
	/help or /start - see this message again.
	`
}

func GreetCommand(userMessage tgbotapi.Update) (feedback string) {
	feedback = fmt.Sprintf("Hi %s, hope your day went well!", userMessage.SentFrom().FirstName)
	return feedback
}

func JoinCommand(userMessage tgbotapi.Update) (feedback string) {
	feedback = "Joined the queue..."
	return feedback
}

func LeaveCommand(userMessage tgbotapi.Update) (feedback string) {
	feedback = "Left the queue..."
	return feedback
}

func HowLongCommand(userMessage tgbotapi.Update) (feedback string) {
	feedback = "Very, very long."
	return feedback
}
