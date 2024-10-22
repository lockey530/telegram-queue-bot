package botaccess

const (
	userHelpFeedback string = `
	Welcome to the Epoque queue bot! Join the queue and you will be able to check how many groups are in front of you! We will also notify you when its nearing your turn so please watch out for our message :)

	/join - join the photobooth queue.

	/leave - leave the photobooth queue.

	/howlong - check how many people are in front of you.

	/help - view the available functions for this bot.

	For more options, check out the 'Menu' button at the bottom left of this chat!

 	IMPORTANT: QUEUE POLICY
	- When you get notified that your turn is reaching soon, gather your group and head to the photo booth IMMEDIATELY
	- Once at the Photo Booth, tell the IC your tele handle to verify your group
	- If you don’t show up within 5 mins of it reaching your turn, you will be kicked out of the queue automatically, you will have to queue again. So make sure to turn on notif and check your messages if you are queueing for your group
	- Walk-ins will NOT be entertained
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

	/*
		user command specific messages
	*/
	nonTextFeedback string = "I don't know what this is :( please send me text commands!"

	nonCommandFeedback string = "Please input a command which starts with '/', like /start"

	invalidCommandFeedback string = "Sorry, I don't recognize your command :("

	joinQueueAlreadyJoined string = "You have already joined this queue!"
	joinQueueSuccess       string = "Joined the queue. (Check the queue with /howlong.)"
	joinQueueFailure       string = "You were unable to join the queue due to an unexpected error :("
	joinQueueClosed        string = "sorry, queue closed!"

	leaveQueueNotJoined string = "It seems you have not joined this queue yet!"
	leaveQueueSuccess   string = "Left the queue..."
	leaveQueueFailure   string = "You were unable to leave the queue due to an unexpected error :("

	/*
		admin command specific messages (feedback for admins)
	*/
	seeQueueStateFailure string = "Something went wrong when accessing the queue state :("

	removeFirstInQueueSuccess string = "successfully removed first peson in queue."
	removeFirstInQueueFailure string = "failed to remove first person in queue."

	kickCommandInvalidArguments string = "input the username to kick. Example: /kick @userABC"
	kickCommandInvalidUser      string = "user inputted was not in queue."
	kickCommandUserFeedback     string = "You have been removed from the queue."
	kickCommandAdminFeedback    string = "First person in queue kicked and notified"

	pingFirstInQueueSuccess string = "First person in queue notified."
	// to append with error reason
	pingFirstInQueueFailure string = "You failed to notify the first person: "

	// admin list
	addAdminInvalidArguments string = "input the username to add as an admin. Example: /addadmin @userABC"
	addAdminFailure          string = "failed to add admin :("
	addAdminSuccess          string = "added successfully!"

	removeAdminInvalidArguments string = "input the username to remove as an admin. Example: /removeadmin @userABC"
	removeAdminFailure          string = "failed to remove admin :("
	removeAdminSuccess          string = "removed successfully!"

	checkAdminListFailure string = "Unable to retrieve admins :("

	// queue (should always be a success - add failure messages otherwise)
	startQueueSuccess string = "queue has been opened for entry."
	stopQueueSuccess  string = "queue has been closed."

	/*
		admin command specific messages (feedback for users)
	*/
	firstInQueueFeedback  string = "It is your turn for the photobooth!"
	secondInQueueFeedback string = "You are the next person in queue - head down to the photobooth!"
	thirdInQueueFeedback  string = "2 groups left in front of you - please head down to the photobooth!"

	kickedFromQueueFeedback string = "You have been kicked from the queue."
)
