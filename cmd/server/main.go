package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/josh1248/nusc-queue-bot/internal/handlers"
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

	// 4: Update bot with autofill menu that informs user of the possible commands available
	registered := []tgbotapi.BotCommand{{Command: "s1", Description: "hello"}}
	menu := tgbotapi.NewSetMyCommands(registered...)
	bot.Send(menu)

	// Set up updates configuration
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Start polling for updates
	updates := bot.GetUpdatesChan(u)
	log.Println("Listening for incoming messages...")

	for update := range updates {

		// Check if the update contains a message
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message.IsCommand() {
			handlers.HandleCommand()
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
