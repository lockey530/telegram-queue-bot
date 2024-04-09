package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/josh1248/nusc-queue-bot/internal/apitoken"
	"github.com/josh1248/nusc-queue-bot/internal/controllers"
	"github.com/josh1248/nusc-queue-bot/internal/dbaccess"
	"github.com/josh1248/nusc-queue-bot/internal/handlers"
)

func main() {
	bot := apitoken.GetBotAPIToken()

	menuOptions := make([]tgbotapi.BotCommand, len(handlers.AvailableCommands))
	for i, command := range handlers.AvailableCommands {
		menuOptions[i] = tgbotapi.BotCommand{Command: command.Command, Description: command.Description}
	}
	menu := tgbotapi.NewSetMyCommands(menuOptions...)
	bot.Send(menu)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// temporary - split and shift this.
	godotenv.Load()

	log.Println("Checking if db data is to be cleared...")
	tmp := os.Getenv("CLEAR_DATA")
	if tmp == "" {
		log.Fatalln("Could not read CLEAR_DATA info from .env")
	}

	toClearData, err := strconv.ParseBool(tmp)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Check done, CLEAR_DATA=%t", toClearData)
	dbaccess.EstablishDBConnection(toClearData)

	updates := bot.GetUpdatesChan(u)
	log.Println("Listening for incoming messages...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		replyMessage := controllers.ReceiveCommand(update, bot)
		_, err := bot.Send(replyMessage)
		if err != nil {
			log.Printf("Error sending message %s\n", err)
		}
	}
}
