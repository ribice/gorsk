package queries

import (
	"fmt"

	"github.com/ribice/gorsk/internal/auth"
)

// UserQueries queries to initialize User table
func UserQueries() (queries []string) {
	userInsert := "INSERT INTO public.users VALUES (1, now(),now(), NULL, 'Admin', 'Admin', 'admin', '%s', 'johndoe@mail.com', NULL, NULL, NULL, NULL, true, 1, 1, 1);"
	queries = append(queries, fmt.Sprintf(userInsert, auth.HashPassword("admin")))
	return queries
}
