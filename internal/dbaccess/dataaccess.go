package dbaccess

import "fmt"

func JoinQueue(user string) error {
	tx := db.MustBegin()
	if _, err := tx.Exec("INSERT INTO queue (user_handle) VALUES ($1)", user); err != nil {
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
	result, err := db.Exec(`DELETE FROM queue WHERE user_handle = $1;`, userHandle)
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
