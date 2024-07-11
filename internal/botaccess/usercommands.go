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

func SeeQueueCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	queueUsers, err := dbaccess.CheckQueueContents()
	if err != nil {
		feedback = "Something went wrong when accessing the queue... blame @joshtwo."
		log.Println(err)
	}

	userToStr := func(user types.QueueUser) string {
		return fmt.Sprintf("@%s %s\n", user.UserHandle, user.Joined_at.Format("15:04:05"))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Total queue length: %v \n\nName		Joined at\n ", len(queueUsers)))
	for _, user := range queueUsers {
		sb.WriteString(userToStr(user))
	}

	feedback = sb.String()
	return feedback
}

// To update: should be variable based on whether you have joined the queue.
func HowLongCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	isInQueue, queueLength, err := dbaccess.CheckQueueLength(userMessage.SentFrom().UserName)
	if err != nil {
		feedback = "Something went wrong when accessing the queue... blame @joshtwo."
		log.Println(err)
		return feedback
	}

	info := func(queueLength int) string {
		if queueLength == 1 {
			return fmt.Sprintf("%s %d %s", "is", queueLength, "group")
		} else {
			return fmt.Sprintf("%s %d %s", "are", queueLength, "groups")
		}
	}

	if isInQueue {
		feedback = fmt.Sprintf("There %s in front of you.", info(queueLength-1))
	} else {
		feedback = fmt.Sprintf("There %s in the queue now. (Join the queue with /join.)", info(queueLength))
	}
	return feedback
}

func KickCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	if len(userMessage.Message.Text) < 6 {
		feedback = "input the username to kick. Example: /kick @userABC"
		return feedback
	}

	telegramHandle := userMessage.Message.Text[6:]
	chatID, err := dbaccess.KickPerson(telegramHandle)
	if err != nil {
		feedback = "You failed to kick the first person: " + err.Error()
		log.Printf("Error sending message %v\n", err)
		return feedback
	}

	msg := tgbotapi.NewMessage(chatID, "You have been kicked from the queue.")
	_, err = bot.Send(msg)
	if err != nil {
		feedback = "You failed to kick " + telegramHandle + " : " + err.Error()
		log.Printf("Error sending message %v\n", err)
		return feedback
	}

	feedback = "Successfully kicked " + telegramHandle
	return feedback
}

// Indirectly called handlers
func PingCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	chatID, err := dbaccess.NotifyQueue(1)
	if err != nil {
		feedback = "You failed to kick the first person: " + err.Error()
		log.Printf("Error sending message, %v\n", err)
		return feedback
	}

	msg := tgbotapi.NewMessage(chatID, "Hey, you are the first person in queue! get moving :D")
	_, err = bot.Send(msg)
	if err != nil {
		feedback = "You failed to kick the first person: " + err.Error()
		log.Printf("Error sending message %v\n", err)
		return feedback
	}

	feedback = "First person in queue notified."
	return feedback
}
