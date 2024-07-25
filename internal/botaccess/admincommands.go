package botaccess

import (
	"fmt"
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
		Command:     "stopqueue",
		Description: "Stop admitting perople into the queue.",
		Handler:     StopQueueCommand,
	},
	{
		Command:     "startqueue",
		Description: "Start admitting perople into the queue.",
		Handler:     StartQueueCommand,
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
	sb.WriteString(fmt.Sprintf("People in queue: %v \nQueue open: %v \n\nName		Joined at\n ", len(queueUsers), queuestatus.IsQueueOpen()))
	for _, user := range queueUsers {
		sb.WriteString(userToStr(user))
	}

	feedback = sb.String()
	return feedback
}

func StopQueueCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	queuestatus.SetQueueClose()
	return "queue successfully stopped."
}

func StartQueueCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	queuestatus.SetQueueOpen()
	return "queue successfully opened."
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

func PingCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	chatID, err := dbaccess.NotifyQueue(1)
	if err != nil {
		feedback = "You failed to kick the first person: " + err.Error()
		log.Printf("Error sending message, %v\n", err)
		return feedback
	}

	msg := tgbotapi.NewMessage(chatID, "It's your turn for the photobooth!")
	_, err = bot.Send(msg)
	if err != nil {
		feedback = "You failed to kick the first person: " + err.Error()
		log.Printf("Error sending message %v\n", err)
		return feedback
	}

	feedback = "First person in queue notified."
	return feedback
}

func KickCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	if len(userMessage.Message.Text) < 7 {
		feedback = "input the username to kick. Example: /kick @userABC"
		return feedback
	}

	telegramHandle := userMessage.Message.Text[7:]
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

	_, nextPersonChatID, err := dbaccess.GetPositionInQueue(1)
	if err != nil {
		return feedback
	}
	msg = tgbotapi.NewMessage(nextPersonChatID, "It's your turn for the photobooth!")
	bot.Send(msg)

	_, nextPersonChatID, err = dbaccess.GetPositionInQueue(2)
	if err != nil {
		return feedback
	}
	msg = tgbotapi.NewMessage(nextPersonChatID, "You are the next person in queue - prepare to head down to the photobooth!")
	bot.Send(msg)

	return feedback

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
	removed, err := dbaccess.RemoveFirstInQueue()
	if err != nil {
		log.Printf("Failed to remove first person in queue %v\n", err)
		return fmt.Sprintf("%v\n", err)
	}

	msg := tgbotapi.NewMessage(removed, "Beep boop, thank you for coming =)")
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message %v\n", err)
		return fmt.Sprintf("Error sending message %v\n", err)
	}

	return "Successfully removed first user in queue."
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
