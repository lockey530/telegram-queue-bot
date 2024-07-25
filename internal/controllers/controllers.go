package controllers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/josh1248/nusc-queue-bot/internal/botaccess"
	"github.com/josh1248/nusc-queue-bot/internal/dbaccess"
)

// receives a user command and sends a reply message.
func ReceiveCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if userMessage.Message == nil {
		return
	}

	reply := tgbotapi.NewMessage(userMessage.Message.Chat.ID, "")
	username := userMessage.SentFrom().UserName

	log.Printf("Processing command %s from @%s...\n", userMessage.Message.Text, username)

	if userMessage.Message.Document != nil {
		log.Printf("Invalid non-text message received from user @%s\n", username)
		reply.Text = botaccess.NonTextHandler(userMessage, bot)
		return
	} else if !userMessage.Message.IsCommand() {
		log.Printf("Invalid output of %s received from user @%s\n", userMessage.Message.Text, username)
		reply.Text = botaccess.NonCommandHandler(userMessage, bot)
		return
	}

	var commandFulfilled bool = false

	isAdmin, err := dbaccess.CheckIfAdmin(username)
	if err != nil {
		log.Println("error: " + err.Error())
	} else if isAdmin {
		for _, command := range botaccess.AdminCommands {
			if userMessage.Message.Command() == command.Command {
				reply.Text = command.Handler(userMessage, bot)
				commandFulfilled = true
				break
			}
		}
	}

	if !commandFulfilled {
		for _, command := range botaccess.UserCommands {
			if userMessage.Message.Command() == command.Command {
				reply.Text = command.Handler(userMessage, bot)
				break
			}
		}
	}

	if reply.Text == "" {
		reply.Text = botaccess.InvalidCommand(userMessage, bot)
	}

	_, err = bot.Send(reply)
	if err != nil {
		log.Printf("Error sending message %s\n", err)
	}
	log.Printf("Processed command %s from @%s.\n", userMessage.Message.Text, username)
}
