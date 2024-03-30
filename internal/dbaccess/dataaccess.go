package dbaccess

// let the database auto-generate the queue_id for us.
const JoinQueueQuery string = `
	INSERT INTO queue (user_handle) VALUES ($1)
`

func JoinQueue(user string) error {
	tx := db.MustBegin()
	_, err := tx.Exec(JoinQueueQuery, user)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
