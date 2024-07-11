package botaccess

const (
	NonTextFeedback string = "I don't know what this is :( please send me text commands!"

	NonCommandFeedback string = "Please input a command which starts with '/', like /start"

	InvalidCommandFeedback string = "Sorry, I don't recognize your command :("

	JoinQueueAlreadyJoined string = "You have already joined this queue!"
	JoinQueueSuccess       string = "Queue joined!"
	JoinQueueFailure       string = "You were unable to join the queue due to an unexpected error :("

	LeaveQueueNotJoined string = "It seems you have not joined this queue yet!"
	LeaveQueueSuccess   string = "Left the queue..."
	LeaveQueueFailure   string = "You were unable to leave the queue due to an unexpected error :("

	SeeQueueStateSuccess string = "Something went wrong when accessing the queue state :("

	KickCommandInvalidArguments string = "input the handle of the person you are kicking, e.g. /kick @xyz"
	KickCommandInvalidUser      string = "user inputted was not in queue."
	KickCommandUserFeedback     string = "You have been removed from the queue."
	KickCommandAdminFeedback    string = "First person in queue kicked and notified"

	PingCommandUserFeedback string = "Hey, you are the first person in queue! get moving :D"
)
