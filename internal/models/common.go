package models

type Data struct {
	Name            string   `json:"name,omitempty"`
	Categories      []string `json:"categories,omitempty"`
	Ratings         Ratings
	YearsInBusiness string   `json:"years_in_business,omitempty"`
	Phone           string   `json:"phone,omitempty"`
	StreetAddress   string   `json:"street_address,omitempty"`
	Locality        string   `json:"locality,omitempty"`
	Location        Location `json:"location,omitempty"`
	URL             string   `json:"url,omitempty"`
	DataUrl         string   `json:"data_url,omitempty"`
}

type Ratings struct {
	BBBRating   string `json:"bbb_rating,omitempty"`
	TARatingURL string `json:"ta_rating_url,omitempty"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
