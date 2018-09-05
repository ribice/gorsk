package queries

// RoleQueries queries to initialize Role table
func RoleQueries() (queries []string) {
	queries = append(queries, "INSERT INTO public.roles VALUES (1, 1, 'SUPER_ADMIN');")
	queries = append(queries, "INSERT INTO public.roles VALUES (2, 2, 'ADMIN');")
	queries = append(queries, "INSERT INTO public.roles VALUES (3, 3, 'COMPANY_ADMIN');")
	queries = append(queries, "INSERT INTO public.roles VALUES (4, 4, 'LOCATION_ADMIN');")
	queries = append(queries, "INSERT INTO public.roles VALUES (5, 5, 'USER');")
	return queries
}
