package botaccess

const (
	userHelpFeedback string = `
	Welcome to the NUSC queue bot!

	/join - join the photobooth queue.

	/leave - leave the photobooth queue.

	/howlong - check how many people are in front of you.

	/help - view the available functions for this bot.

	For more options, check out the 'Menu' button at the bottom left of this chat!
	`

	adminHelpFeedback string = `
	Bot controls (admins):

	/seequeue -  see who is in the queue now.

	/ping - remind the first person in the queue to come.

	/done - remove the first person from the queue once they have finished their photo-taking.
	
	/kick @handle - remove a person from the queue, e.g. /kick @abc

	/stopqueue - stop allowing people to join the queue.

	/startqueue - allow people to join the queue.

	/adminlist - see who has the ability to control the bot.

	/addadmin @handle - allow another person to control the bot, e.g. /addadmin @abc

	/removeadmin @handle - remove an admin.

	/help - view all available functions for this bot.

	Bot controls (users):

	/join - join the photobooth queue.

	/leave - leave the photobooth queue.

	/howlong - check how many people are in front of you.
	`

	// users
	nonTextFeedback string = "I don't know what this is :( please send me text commands!"

	nonCommandFeedback string = "Please input a command which starts with '/', like /start"

	invalidCommandFeedback string = "Sorry, I don't recognize your command :("

	joinQueueAlreadyJoined string = "You have already joined this queue!"
	joinQueueSuccess       string = "Joined the queue. (Check the queue with /howlong.)"
	joinQueueFailure       string = "You were unable to join the queue due to an unexpected error :("

	leaveQueueNotJoined string = "It seems you have not joined this queue yet!"
	leaveQueueSuccess   string = "Left the queue..."
	leaveQueueFailure   string = "You were unable to leave the queue due to an unexpected error :("

	// admins
	seeQueueStateSuccess string = "Something went wrong when accessing the queue state :("

	removeFirstInQueueSuccess string = "successfully removed first peson in queue. Removed: "
	removeFirstInQueueFailure string = "failed to remove first person in queue. Error: "

	kickCommandInvalidArguments string = "input the handle of the person you are kicking, e.g. /kick @xyz"
	kickCommandInvalidUser      string = "user inputted was not in queue."
	kickCommandUserFeedback     string = "You have been removed from the queue."
	kickCommandAdminFeedback    string = "First person in queue kicked and notified"

	// append reason behind this string.
	addAdminFailure string = "failed to add admin :("
	addAdminSuccess string = "added successfully!"
	// append reason behind this string.
	removeAdminFailure string = "failed to remove admin :("
	removeAdminSuccess string = "removed successfully!"

	pingCommandUserFeedback string = "Hey, you are the first person in queue! get moving :D"

	checkAdminListFailure string = "Unable to retrieve admins :("
)
