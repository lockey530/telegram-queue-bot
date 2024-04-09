// helps to remove feedback logic from app logic, and allow easy tweaking in one place.
package userfeedback

import (
	"fmt"
)

// handlers
const (
	NonTextFeedback string = "I don't know what this is :( please send me text commands!"

	NonCommandFeedback string = "Please input a command which starts with '/', like /start"

	InvalidCommandFeedback string = "Sorry, I don't recognize your command :("

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

const (
	JoinQueueAlreadyJoined string = "You have already joined this queue!"
	JoinQueueSuccess       string = "Queue joined!"
	JoinQueueFailure       string = "You were unable to join the queue due to an unexpected error :("

	LeaveQueueNotJoined string = "It seems you have not joined this queue yet!"
	LeaveQueueSuccess   string = "Left the queue..."
	LeaveQueueFailure   string = "You were unable to leave the queue due to an unexpected error :("

	SeeQueueStateSuccess string = "Something went wrong when accessing the queue state :("
)

func HowLongFeedback(queueLength int) string {
	var info string
	if queueLength == 1 {
		info = fmt.Sprintf("%s %d %s", "is", queueLength, "group")
	} else {
		info = fmt.Sprintf("%s %d %s", "are", queueLength, "groups")
	}

	return fmt.Sprintf("There %s in the queue now.", info)
}

const (
	KickCommandInvalidArguments string = "input the handle of the person you are kicking, e.g. /kick @xyz"
	KickCommandInvalidUser      string = "user inputted was not in queue."
	KickCommandUserFeedback     string = "You have been removed from the queue."
	KickCommandAdminFeedback    string = "First person in queue kicked and notified"

	PingCommandUserFeedback string = "Hey, you are the first person in queue! get moving :D"
)
