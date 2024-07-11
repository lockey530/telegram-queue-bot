package botaccess

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/josh1248/nusc-queue-bot/internal/commandtypes"
	"github.com/josh1248/nusc-queue-bot/internal/dbaccess"
)

// Command and description is hard-coded within the HelpFunction for circular dependencies.
// You will need to update the help function accordingly if you want to change this definition.
// If this becomes too troublesome, consider using the reflect package to store the handler
// function names under a string first before using reflect.ValueOf().Call()
var AvailableCommands = []commandtypes.AcceptedCommands{
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
		Command:     "seequeue",
		Description: "(admins only - check this): see the contents of the queue now.",
		Handler:     SeeQueueCommand,
	},
	{
		Command:     "ping",
		Description: "(admins only - check this): send a reminder to the first person in queue.",
		Handler:     PingCommand,
	},
	{
		Command:     "kick",
		Description: "(admins only - check this): remove the person at some set position",
		Handler:     KickCommand,
	},
	{
		Command:     "howlong",
		Description: "Returns the expected time to wait in the queue",
		Handler:     HowLongCommand,
	},
	{
		Command:     "help",
		Description: "Explains the main functionalities of the bot.",
		Handler:     HelpCommand,
	},
	{
		Command:     "start",
		Description: "Explains the main functionalities of the bot.",
		Handler:     HelpCommand,
	},
}

func NonTextHandler(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	return NonTextFeedback
}

const nonCommandFeedback string = "Please input a command which starts with '/', like /start"

func NonCommandHandler(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	return nonCommandFeedback
}

const invalidCommandFeedback string = "Sorry, I don't recognize your command :("

func InvalidCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	return invalidCommandFeedback
}

const helpFeepback string = `
Welcome to the queue bot~

/join - join the photobooth queue!

/leave - leave the photobooth queue if you have previously joined.

/howlong - check how many people are in front of you.

For more options, check out the 'Menu' button at the bottom left of this chat!
`

func HelpCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	return helpFeepback
}

func JoinCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	err := dbaccess.JoinQueue(userMessage)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			feedback = "You have already joined this queue!"
		} else {
			feedback = "You were unable to join the queue due to an unexpected error :("
		}
		log.Println(err)
	} else {
		feedback = "Joined the queue..."
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
	feedback, err := dbaccess.CheckQueueContents()
	if err != nil {
		feedback = "Something went wrong when accessing the queue... blame @joshtwo."
		log.Println(err)
	}
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

	var info string
	if queueLength == 1 {
		info = fmt.Sprintf("%s %d %s", "is", queueLength, "group")
	} else {
		info = fmt.Sprintf("%s %d %s", "are", queueLength, "groups")
	}

	if isInQueue {
		feedback = fmt.Sprintf("There %s in front of you now.", info)
	} else {
		feedback = fmt.Sprintf("There %s total in the queue now.", info)
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
