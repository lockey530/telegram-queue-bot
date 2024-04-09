package apitoken

import (
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var initialized bool = false
var botAPI *tgbotapi.BotAPI

// This function should be initialized in the main function before use in subsequent areas,
// which can call this function to retrieve the bot token.
// implicit .env function arguments during setup: remote_deploy, bot_token.
func GetBotAPIToken() *tgbotapi.BotAPI {
	if initialized {
		return botAPI
	}

	var err error

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
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatalln("Failed to read off local .env variables: ", err)
		}
	}

	log.Println("Connecting to bot...")

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("Could not read token")
	} else if token == "(your bot API here)" {
		log.Fatalln(`You forgot to input your API token within the .env file! Setup:
			1. Duplicate the .envSETUP file.
			2. Change the file name/extension to ".env".
			3. Replace (your bot API here) with the Telegram API token given in @BotFather.`)
	}

	botAPI, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		if strings.Contains(err.Error(), token) {
			log.Fatalln(`Failed to connect to Telegram bot API. 
			Check that you have entered the API key correctly,
			and that you are connected to the internet.`)
		} else {
			log.Fatalf("Error creating bot: %v", err)
		}
	}
	log.Println("Successfully connected!")

	initialized = true
	return botAPI
}
