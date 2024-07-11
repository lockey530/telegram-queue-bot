package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/josh1248/nusc-queue-bot/internal/botaccess"
	"github.com/josh1248/nusc-queue-bot/internal/controllers"
	"github.com/josh1248/nusc-queue-bot/internal/dbaccess"
)

func main() {
	godotenv.Load()

	bot := botaccess.InitializeBotAPIConnection()

	tmp := os.Getenv("RESET_DB")
	if tmp == "" {
		log.Fatalln("Could not read RESET_DB info from .env")
	}

	toClearData, err := strconv.ParseBool(tmp)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Check done, CLEAR_DATA=%t", toClearData)
	dbaccess.EstablishDBConnection(toClearData)

	menuOptions := make([]tgbotapi.BotCommand, len(botaccess.UserCommands))
	for i, command := range botaccess.UserCommands {
		menuOptions[i] = tgbotapi.BotCommand{Command: command.Command, Description: command.Description}
	}
	menu := tgbotapi.NewSetMyCommands(menuOptions...)
	bot.Send(menu)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	log.Println("Listening for incoming messages...")

	for update := range updates {
		controllers.ReceiveCommand(update, bot)
	}
}
