package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/josh1248/nusc-queue-bot/internal/controller"
)

func main() {
	log.Println("Connecting to bot...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("Could not read token")
	} else if token == "(your bot API here)" {
		log.Fatalln(`You forgot to input your API token within the .env file! Setup:
			1. Duplicate the .envSETUP file.
			2. Change the file name/extension to ".env".
			3. Replace (your bot API here) with the Telegram API token given in @BotFather.`)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	log.Println("Successfully connected!")

	menuOptions := []tgbotapi.BotCommand{
		{
			Command:     "join",
			Description: "Join the virtual queue for the photobooth.",
		},
		{
			Command:     "leave",
			Description: "Leave the virtual queue for the photobooth.",
		},
		{
			Command:     "howlong",
			Description: "Returns the expected time to wait in the queue",
		},
		{
			Command:     "help",
			Description: "Explains the main functionalities of the bot.",
		},
		{
			Command:     "greet",
			Description: "The bot is friendly :)",
		},
		{
			Command:     "start",
			Description: "Explains the main functionalities of the bot.",
		},
	}
	menu := tgbotapi.NewSetMyCommands(menuOptions...)
	bot.Send(menu)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	log.Println("Listening for incoming messages...")

	for update := range updates {
		if update.Message == nil {
			continue
		}
		replyMessage := controller.ReceiveCommand(update)
		_, err := bot.Send(replyMessage)
		if err != nil {
			log.Printf("Error sending message %s\n", err)
		}
	}
}
