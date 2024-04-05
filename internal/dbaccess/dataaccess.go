package dbaccess

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func JoinQueue(user tgbotapi.Update) error {
	tx := db.MustBegin()
	_, err := tx.Exec("INSERT INTO queue (user_handle, chat_id) VALUES ($1, $2)",
		user.SentFrom().UserName, user.SentFrom().ID)

	if err != nil {
		return fmt.Errorf("insertion query failed to execute. %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction for insertion failed to commit. %v", err)
	}

	return nil
}

// wrapper or monad time?
func CheckQueue() (string, error) {
	queue := []QueueUser{}
	if err := db.Select(&queue, "SELECT * FROM queue"); err != nil {
		return "", fmt.Errorf("failed to get queue state. %v", err)
	}

	return fmt.Sprintf("%v", queue), nil
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
	user := QueueUser{}
	// change this if select query is no longer guaranteed to return 1 row.
	if err := db.Get(&user, "SELECT (chat_id) FROM queue ORDER BY joined_at"); err != nil {
		return 0, fmt.Errorf("failed to get first user in queue: %v", err)
	}

	return user.ChatID, nil
}

func KickPerson(position int64) (chatID int64, err error) {
	// this will need to be adjusted after implementing the waiting feature.
	_, err = db.Exec(`
		DELETE FROM queue 
		WHERE chat_id = (
			SELECT chat_id
			FROM queue
			ORDER BY joined_at
			LIMIT 1
		);
	`)

	if err != nil {
		return 0, fmt.Errorf("failed to kick %vth position from queue. %v", position, err)
	}

	return chatID, nil
}
