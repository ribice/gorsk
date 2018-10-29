package gorsk

// Location represents company location model
type Location struct {
	Base
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Address string `json:"address"`

	CompanyID int `json:"company_id"`
}
