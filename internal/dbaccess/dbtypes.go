package dbaccess

import "time"

type QueueUser struct {
	QueueID    uint64    `db:"queue_id" note:"Postgres-generated identifier"`
	UserHandle string    `db:"user_handle" note:"Refers to the Telegram handle"`
	ChatID     int64     `db:"chat_id" note:"Refers to the ID used to identify message chats"`
	Joined_at  time.Time `db:"joined_at" note:"Refers to the time the user joined a queue"`
}

const queueSchema string = `
	CREATE TABLE queue (
		queue_id				SERIAL PRIMARY KEY,
		user_handle		 		TEXT 	UNIQUE NOT NULL,
		chat_id					BIGSERIAL UNIQUE NOT NULL,
		joined_at 				TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`

type AdminUser struct {
	AdminID     uint64 `db:"admin_id" note:"Telegram handle chat ID?"`
	AdminHandle string `db:"admin_handle" note:"Refers to the Telegram handle"`
}

const adminSchema string = `
	CREATE TABLE admins (
		admin_id				SERIAL PRIMARY KEY,
		admin_handle		 	TEXT 	UNIQUE NOT NULL
	);
`

// https://stackoverflow.com/questions/20582500/how-to-check-if-a-table-exists-in-a-given-schema
const checkTableExistenceQuery string = `
	SELECT EXISTS (
		SELECT FROM pg_tables
		WHERE  	schemaname = 'public'
		AND    	tablename  = $1
		);
`
