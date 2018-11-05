package gorsk

// AccessRole represents access role type
type AccessRole int

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessRole = 100

	// AdminRole has admin specific permissions
	AdminRole AccessRole = 110

	// CompanyAdminRole can edit company specific things
	CompanyAdminRole AccessRole = 120

	// LocationAdminRole can edit location specific things
	LocationAdminRole AccessRole = 130

	// UserRole is a standard user
	UserRole AccessRole = 200
)

// Role model
type Role struct {
	ID          AccessRole `json:"id"`
	AccessLevel AccessRole `json:"access_level"`
	Name        string     `json:"name"`
}
