package dbaccess

type QueueUser struct {
	queueID    uint64 `db:"queue_id" note:"Postgres-generated identifier"`
	userHandle string `db:"user_handle" note:"Refers to the Telegram handle"`
}

const queueSchema string = `
	CREATE TABLE queue (
		queue_id				SERIAL PRIMARY KEY,
		user_handle		 		TEXT 	UNIQUE NOT NULL
	);
`

type AdminUser struct {
	adminID     uint64 `db:"admin_id" note:"Postgres-generated identifier"`
	adminHandle string `db:"admin_handle" note:"Refers to the Telegram handle"`
}

const adminSchema string = `
	CREATE TABLE admins (
		admin_id				SERIAL PRIMARY KEY,
		admin_handle		 	TEXT 	UNIQUE NOT NULL
	);
`

// https://stackoverflow.com/questions/20582500/how-to-check-if-a-table-exists-in-a-given-schema
const checkExistenceQuery string = `
	SELECT EXISTS (
		SELECT FROM pg_tables
		WHERE  	schemaname = 'public'
		AND    	tablename  = $1
		);
`
