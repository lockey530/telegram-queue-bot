package dbaccess

import (
	"fmt"

	types "github.com/josh1248/nusc-queue-bot/internal/types"
)

func AddDummy() {
	JoinQueue("test_user", 12345)
}

func JoinQueue(username string, chatID int64) error {
	tx := db.MustBegin()
	_, err := tx.Exec("INSERT INTO queue (user_handle, chat_id) VALUES ($1, $2);",
		username, chatID)
	if err != nil {
		return fmt.Errorf("insertion query failed to execute. %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction for insertion failed to commit. %v", err)
	}
	return nil
}

func CheckQueueContents() ([]types.QueueUser, error) {
	queue := []types.QueueUser{}
	if err := db.Select(&queue, `
	SELECT 
		queue_id, user_handle, chat_id, (joined_at AT TIME ZONE 'UTC' AT TIME ZONE 'UTC+8')
	FROM queue;`); err != nil {
		return nil, fmt.Errorf("failed to get queue state. %v", err)
	}

	return queue, nil
}

func CheckIfInQueue(userHandle string) (bool, error) {
	var isInQueue bool

	fmt.Println(userHandle)
	if err := db.Get(&isInQueue, "SELECT EXISTS (SELECT 1 FROM queue WHERE user_handle = $1);", userHandle); err != nil {
		return false, fmt.Errorf("failed to get queue state. %v", err)
	}

	return isInQueue, nil
}

func CheckQueueLength(userHandle string) (bool, int, error) {
	isInQueue, err := CheckIfInQueue(userHandle)
	if err != nil {
		return isInQueue, -1, fmt.Errorf("failed to get queue state. %v", err)
	}

	var queueLength int
	// https://wiki.postgresql.org/wiki/Count_estimate for the method which requires ANALYZE
	// but can be faster.
	if !isInQueue {
		// update this
		err = db.Get(&queueLength, "SELECT count(*) FROM queue;")
	} else {
		err = db.Get(&queueLength, `
			SELECT 
				count(*)
			FROM 
				queue
			WHERE
				joined_at <= (
					SELECT
						joined_at
					FROM
						queue
					WHERE
						user_handle = $1
				);
			`, userHandle)
	}

	return isInQueue, queueLength, err
}

func LeaveQueue(userHandle string) error {
	result, err := db.Exec("DELETE FROM queue WHERE user_handle = $1;", userHandle)
	if err != nil {
		return fmt.Errorf("failed to leave queue. %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to leave queue. %v", err)
	} else if affected == 0 {
		return fmt.Errorf("user not in queue")
	}

	return nil
}

func NotifyQueue(position int64) (chatID int64, err error) {
	user := types.QueueUser{}
	if err := db.Get(&user, "SELECT (chat_id) FROM queue ORDER BY joined_at OFFSET $1 LIMIT 1;", position-1); err != nil {
		return 0, fmt.Errorf("failed to get first user in queue: %v", err)
	}

	return user.ChatID, nil
}
