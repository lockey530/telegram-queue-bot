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
		Command:     "kick",
		Description: "remove the specified person from the queue.",
		Handler:     KickCommand,
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
		Command:     "addDummy12345",
		Description: "FOR TESTING ONLY",
		Handler:     AddDummyCommand,
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

func AddDummyCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	dbaccess.AddDummy()
	return "dummy user added"
}

func SeeQueueCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	queueUsers, err := dbaccess.CheckQueueContents()
	if err != nil {
		feedback = seeQueueStateFailure
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
	return stopQueueSuccess
}

func StartQueueCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	queuestatus.SetQueueOpen()
	return startQueueSuccess
}

// To update: should be variable based on whether you have joined the queue.
func HowLongCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	isInQueue, queueLength, err := dbaccess.CheckQueueLength(userMessage.SentFrom().UserName)
	if err != nil {
		feedback = joinQueueFailure
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
		feedback = pingFirstInQueueFailure + err.Error()
		log.Printf("Error sending message, %v\n", err)
		return feedback
	}

	msg := tgbotapi.NewMessage(chatID, firstInQueueFeedback)
	_, err = bot.Send(msg)
	if err != nil {
		feedback = pingFirstInQueueFailure + err.Error()
		log.Printf("Error sending message %v\n", err)
		return feedback
	}

	feedback = "First person in queue notified."
	return feedback
}

func KickCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	if len(userMessage.Message.Text) < 7 {
		feedback = kickCommandInvalidArguments
		return feedback
	}

	telegramHandle := userMessage.Message.Text[7:]
	chatID, err := dbaccess.KickPerson(telegramHandle)
	if err != nil {
		feedback = "You failed to kick the specified user: " + err.Error()
		log.Printf("Error sending message %v\n", err)
		return feedback
	}

	msg := tgbotapi.NewMessage(chatID, kickedFromQueueFeedback)
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
	msg = tgbotapi.NewMessage(nextPersonChatID, firstInQueueFeedback)
	bot.Send(msg)

	_, nextPersonChatID, err = dbaccess.GetPositionInQueue(2)
	if err != nil {
		return feedback
	}
	msg = tgbotapi.NewMessage(nextPersonChatID, secondInQueueFeedback)
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
		return removeFirstInQueueFailure + err.Error()
	}

	_, nextPersonChatID, _ := dbaccess.GetPositionInQueue(1)
	if nextPersonChatID != -1 {
		msg = tgbotapi.NewMessage(nextPersonChatID, firstInQueueFeedback)
		bot.Send(msg)
	}

	_, nextPersonChatID, _ = dbaccess.GetPositionInQueue(2)
	if nextPersonChatID != -1 {
		msg = tgbotapi.NewMessage(nextPersonChatID, secondInQueueFeedback)
		bot.Send(msg)
	}

	_, nextPersonChatID, _ = dbaccess.GetPositionInQueue(3)
	if nextPersonChatID != -1 {
		msg = tgbotapi.NewMessage(nextPersonChatID, thirdInQueueFeedback)
		bot.Send(msg)
	}

	feedback = removeFirstInQueueSuccess
	return feedback
}

func AddAdminCommand(userMessage tgbotapi.Update, bot *tgbotapi.BotAPI) (feedback string) {
	if len(userMessage.Message.Text) < 12 {
		feedback = addAdminInvalidArguments
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
		feedback = removeAdminInvalidArguments
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
