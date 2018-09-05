package queries

// DBSetupQueries queries to create/initialize/enable any extension or schema in the database
func DBSetupQueries() (queries []string) {
	queries = append(queries, `CREATE EXTENSION "uuid-ossp";`)
	return queries
}
