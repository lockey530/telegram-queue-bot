package dbaccess

// https://stackoverflow.com/questions/20582500/how-to-check-if-a-table-exists-in-a-given-schema
const checkTableExistenceQuery string = `
	SELECT EXISTS (
		SELECT FROM pg_tables
		WHERE  	schemaname = 'public'
		AND    	tablename  = $1
		);
`
