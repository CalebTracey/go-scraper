package models

type Data struct {
	Name       string   `json:"name,omitempty"`
	Categories []string `json:"categories,omitempty"`
	Ratings    string   `json:"ratings,omitempty"`
	Phone      string   `json:"phone,omitempty"`
	Address    string   `json:"address,omitempty"`
	Location   Location `json:"location,omitempty"`
	URL        string   `json:"url,omitempty"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
