package queries

// CompanyQueries queries to initialize Company table
func CompanyQueries() (queries []string) {
	queries = append(queries, "INSERT INTO public.companies VALUES (1, now(), now(), NULL, 'admin_company', true);")
	return queries
}
