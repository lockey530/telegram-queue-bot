package botaccess

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/josh1248/nusc-queue-bot/internal/dbaccess"
	"github.com/josh1248/nusc-queue-bot/internal/types"
)

// Command and description is hard-coded within the HelpFunction for circular dependencies.
// You will need to update the help function accordingly if you want to change this definition.
// If this becomes too troublesome, consider using the reflect package to store the handler
// function names under a string first before using reflect.ValueOf().Call()
var AdminCommands = []types.AcceptedCommands{
	{
		Command:     "seequeue",
		Description: "see who is in the queue now.",
		Handler:     SeeQueueCommand,
	},
	{
		Command:     "ping",
		Description: "send a reminder to the first person in queue.",
		Handler:     PingCommand,
	},
	{
		Command:     "done",
		Description: "remove the first person from the queue once they have finished their photo-taking.",
		Handler:     RemoveFirstInQueueCommand,
	},
	{
		Command:     "adminlist",
		Description: "see who has the ability to control the bot.",
		Handler:     CheckAdminListCommand,
	},
	{
		Command:     "addadmin",
		Description: "allow another person to control the bot, e.g. /addadmin @abc",
		Handler:     AddAdminCommand,
	},
	{
		Command:     "removeadmin",
		Description: "remove admin rights for some user, e.g. /removeadmin @abc",
		Handler:     RemoveAdminCommand,
	},
	{
		Command:     "help",
		Description: "Explains the main functionalities of the bot.",
		Handler:     AdminHelpCommand,
	},
	{
		Command:     "start",
		Description: "Explains the main functionalities of the bot.",
		Handler:     AdminHelpCommand,
	},
}

func AdminHelpCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	return adminHelpFeedback
}

func CheckAdminListCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	admins, err := dbaccess.SeeAdminList(userMessage.SentFrom().UserName)
	if err != nil {
		log.Println(err)
		return checkAdminListFailure
	} else {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Admin list. Total admins: %v\n", len(admins)))
		for _, admin := range admins {
			sb.WriteString("@" + admin + "\n")
		}
		return sb.String()
	}
}

func RemoveFirstInQueueCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	_, err := dbaccess.RemoveFirstInQueue(userMessage.SentFrom().UserName)
	if err != nil {
		log.Println(err)
		return removeFirstInQueueFailure
	} else {
		return removeFirstInQueueSuccess
	}
}

func AddAdminCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	if len(userMessage.Message.Text) < 12 {
		feedback = "input the username to add as an admin. Example: /addadmin @userABC"
		return feedback
	}
	telegramHandle := userMessage.Message.Text[11:]

	err := dbaccess.AddAdmin(telegramHandle, userMessage.SentFrom().UserName)
	if err != nil {
		log.Println(err)
		return addAdminFailure
	} else {
		return addAdminSuccess
	}
}

func RemoveAdminCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	if len(userMessage.Message.Text) < 15 {
		feedback = "input the username to remove as an admin. Example: /removeadmin @userABC"
		return feedback
	}
	telegramHandle := userMessage.Message.Text[14:]

	issue, err := dbaccess.RemoveAdmin(telegramHandle, userMessage.SentFrom().UserName)
	if err != nil {
		log.Println(err)
		if issue != "" {
			return issue
		}
		return removeAdminFailure
	} else {
		return removeAdminSuccess
	}
}
