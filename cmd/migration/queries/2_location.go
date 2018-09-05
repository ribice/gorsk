package queries

// LocationQueries queries to initialize Location table
func LocationQueries() (queries []string) {
	queries = append(queries, "INSERT INTO public.locations VALUES (1, now(), now(), NULL, 'admin_location', true, 'admin_address', 1);")
	return queries
}
