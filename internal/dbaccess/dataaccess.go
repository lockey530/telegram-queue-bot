package dbaccess

import "fmt"

// let the database auto-generate the queue_id for us.
const JoinQueueQuery string = `
	INSERT INTO queue (user_handle) VALUES ($1)
`

func JoinQueue(user string) error {
	tx := db.MustBegin()
	if _, err := tx.Exec(JoinQueueQuery, user); err != nil {
		return fmt.Errorf("insertion query failed to execute. %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction for insertion failed to commit. %v", err)
	}

	return nil
}
