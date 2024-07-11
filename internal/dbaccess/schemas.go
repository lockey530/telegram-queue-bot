package dbaccess

const (
	adminSchema string = `
		CREATE TABLE admins (
			admin_id				SERIAL PRIMARY KEY,
			admin_handle		 	TEXT 	UNIQUE NOT NULL
		);
	`

	queueSchema string = `
		CREATE TABLE queue (
			queue_id				SERIAL PRIMARY KEY,
			user_handle		 		TEXT 	UNIQUE NOT NULL,
			chat_id					BIGSERIAL UNIQUE NOT NULL,
			joined_at 				TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	indexQueueSchema string = `
		CREATE INDEX chat_id ON queue (chat_id);
	`
)
