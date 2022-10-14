package models

type Data struct {
	Name          string   `json:"name,omitempty"`
	Categories    []string `json:"categories,omitempty"`
	Ratings       string   `json:"ratings,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	StreetAddress string   `json:"street_address,omitempty"`
	Locality      string   `json:"locality,omitempty"`
	URL           string   `json:"url,omitempty"`
}
