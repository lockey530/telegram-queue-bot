package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
	feedback = "here are my commands: (To be written)"
	return feedback
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
