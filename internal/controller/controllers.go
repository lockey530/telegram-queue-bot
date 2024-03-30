package controller

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// receives a user message and returns with message to be sent.
func ReceiveCommand(userMessage tgbotapi.Update) (reply tgbotapi.MessageConfig) {
	reply = tgbotapi.NewMessage(userMessage.Message.Chat.ID, "")
	username := userMessage.SentFrom().UserName

	if userMessage.Message.Document != nil {
		log.Printf("Invalid non-text message received from user @%s\n", username)
		reply.Text = "Please send me text only :("
	} else if !userMessage.Message.IsCommand() {
		log.Printf("Invalid output of %s received from user @%s\n", userMessage.Message.Text, username)
		reply.Text = "Please input a command (e.g. /help, /join)."
	} else {
		log.Printf("Valid command %s recevied from @%s and processed\n", userMessage.Message.Text, username)
		reply.Text = fmt.Sprintf("Hi @%s, command received!", username)
	}

	return reply
}
