package botaccess

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/josh1248/nusc-queue-bot/internal/dbaccess"
	"github.com/josh1248/nusc-queue-bot/internal/queuestatus"
	"github.com/josh1248/nusc-queue-bot/internal/types"
)

// Command and description is hard-coded within the HelpFunction for circular dependencies.
// You will need to update the help function accordingly if you want to change this definition.
// If this becomes too troublesome, consider using the reflect package to store the handler
// function names under a string first before using reflect.ValueOf().Call()
var UserCommands = []types.AcceptedCommands{
	{
		Command:     "join",
		Description: "Join the virtual queue for the photobooth.",
		Handler:     JoinCommand,
	},
	{
		Command:     "leave",
		Description: "Leave the virtual queue for the photobooth.",
		Handler:     LeaveCommand,
	},
	{
		Command:     "howlong",
		Description: "Returns the expected time to wait in the queue",
		Handler:     HowLongCommand,
	},
	{
		Command:     "help",
		Description: "Explains the main functionalities of the bot.",
		Handler:     UserHelpCommand,
	},
	{
		Command:     "start",
		Description: "Explains the main functionalities of the bot.",
		Handler:     UserHelpCommand,
	},
}

func NonTextHandler(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	return nonTextFeedback
}

func NonCommandHandler(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	return nonCommandFeedback
}

func InvalidCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	return invalidCommandFeedback
}

func UserHelpCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	return userHelpFeedback
}

func JoinCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	if !queuestatus.IsQueueOpen() {
		return "sorry, queue closed!"
	}
	err := dbaccess.JoinQueue(userMessage.SentFrom().UserName, userMessage.SentFrom().ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			feedback = "You have already joined this queue!"
		} else {
			feedback = "You were unable to join the queue due to an unexpected error :("
		}
		log.Println(err)
	} else {
		return joinQueueSuccess
	}

	return feedback
}

func LeaveCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	err := dbaccess.LeaveQueue(userMessage.SentFrom().UserName)
	if err != nil {
		if strings.Contains(err.Error(), "user not in queue") {
			feedback = "It seems you have not joined this queue yet!"
		} else {
			feedback = "You were unable to leave the queue due to an unexpected error :("
		}
		log.Println(err)
	} else {
		feedback = "Left the queue..."
	}
	return feedback
}
