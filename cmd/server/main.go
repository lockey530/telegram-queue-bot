package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/josh1248/nusc-queue-bot/internal/controllers"
	"github.com/josh1248/nusc-queue-bot/internal/dbaccess"
	"github.com/josh1248/nusc-queue-bot/internal/handlers"
)

func main() {
	// Per the documentation, this function will not override an env variable that already exists.
	// Therefore, it is safe to run this code in both local and remote deployments.
	godotenv.Load()
	log.Println(os.LookupEnv("REMOTE_DEPLOY"))
	remoteDeploy, err := strconv.ParseBool(os.Getenv("REMOTE_DEPLOY"))
	if err != nil {
		log.Fatalln("environment variables improperly set up: ", err)
	}

	if remoteDeploy {
		log.Println("Remote deployment of app started.")
	} else {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln("Failed to read off local .env variables: ", err)
		}
	}

	log.Println("Connecting to bot...")
	/*
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln("Error loading .env file")
		}*/

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
		// prevent logs from exposing your bot token.
		if err.Error()[0:4] == "Post" {
			log.Fatalln(`Failed to connect to Telegram bot API. 
			Check that you have entered the API key correctly,
			and that you are connected to the internet.`)
		} else {
			log.Fatalf("Error creating bot: %v", err)
		}
	}
	log.Println("Successfully connected!")

	menuOptions := make([]tgbotapi.BotCommand, len(handlers.AvailableCommands))
	for i, command := range handlers.AvailableCommands {
		menuOptions[i] = tgbotapi.BotCommand{Command: command.Command, Description: command.Description}
	}
	menu := tgbotapi.NewSetMyCommands(menuOptions...)
	bot.Send(menu)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

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
