package botaccess

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var botAPI *tgbotapi.BotAPI

// implicit .env arguments: BOT_TOKEN.
func InitializeBotAPIConnection() *tgbotapi.BotAPI {
	// Per the documentation, this function will not override an env variable that already exists.
	// Therefore, it is safe to run this code in both local and remote deployments.
	godotenv.Load()

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

	var err error
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

	return botAPI
}

func GetBotAPIConnection() (*tgbotapi.BotAPI, error) {
	if botAPI == nil {
		return nil, fmt.Errorf("retrieval of un-initialized bot API")
	}

	return botAPI, nil
}

// handlers
const (
	HelpFeepback string = `
	Welcome to the queue bot~

	/join - join the photobooth queue!

	/leave - leave the photobooth queue if you have previously joined.

	/wait - (Not supported yet) need some time? place yourself 5 slots behind the queue (1-time only).

	/help or /start - see this message again.

	For more options, check out the 'Menu' button at the bottom left of this chat!
	`
)

func GreetFeedback(name string) string {
	return fmt.Sprintf("Hi %s, hope your day is going well :)", name)
}

func HowLongFeedback(queueLength int) string {
	var info string
	if queueLength == 1 {
		info = fmt.Sprintf("%s %d %s", "is", queueLength, "group")
	} else {
		info = fmt.Sprintf("%s %d %s", "are", queueLength, "groups")
	}

	return fmt.Sprintf("There %s in the queue now.", info)
}
