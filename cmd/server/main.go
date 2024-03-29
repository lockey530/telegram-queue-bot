package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/josh1248/nusc-queue-bot/internal/botcommands"
	"github.com/josh1248/nusc-queue-bot/internal/importtest"
)

func main() {

	// Setup procedure

	// 1: connect with API
	log.Println("Connecting to bot...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	// 2: Check existence of env key
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("Could not read token")
	} else if token == "(your bot API here)" {
		log.Fatalln(`You forgot to input your API token within the .env file! Setup:
			1. Duplicate the .envSETUP file.
			2. Change the file name to ".env".
			3. Replace (your bot API here) with the Telegram API token given in @BotFather.`)
	}

	// 3: Check that API key is legitimate
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	log.Println("Successfully connected!")

	// Temp skeleton functionality implementation
	// abstract into an "internal" folder once program becomes complex enough, and repeated patterns are observed.

	// Set up updates configuration
	updates := tgbotapi.NewUpdate(0)
	updates.Timeout = 60

	// Start polling for updates
	updatesChannel, err := bot.GetUpdatesChan(updates)
	if err != nil {
		log.Fatalf("Error getting updates channel: %v", err)
	}
	log.Println("Listening for incoming messages...")

	// Process incoming updates
	for update := range updatesChannel {
		log.Println(botcommands.Hi)
		log.Println(importtest.X())

		// Check if the update contains a message
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "hi":
				msg.Text = "Hello!"
			case "bye":
				msg.Text = "Goodbye!"
			case "name":
				// Check if the command has an argument
				if len(update.Message.CommandArguments()) > 0 {
					msg.Text = "Your name is " + update.Message.CommandArguments() + "."
				} else {
					msg.Text = "Please provide your name with the command, e.g., /name John."
				}
			default:
				msg.Text = "Sorry, I don't understand your command =("
			}
		}

		// Send the message
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}
