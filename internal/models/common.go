package models

type Data struct {
	Name          string   `json:"name,omitempty"`
	Categories    []string `json:"categories,omitempty"`
	Ratings       string   `json:"ratings,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	StreetAddress string   `json:"street_address,omitempty"`
	Locality      string   `json:"locality,omitempty"`
	Location      Location `json:"location,omitempty"`
	URL           string   `json:"url,omitempty"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
