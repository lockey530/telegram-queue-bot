package queuestatus

var queueOpen bool = true

func IsQueueOpen() bool {
	return queueOpen
}

func SetQueueClose() {
	queueOpen = false
}

func SetQueueOpen() {
	queueOpen = true
}
