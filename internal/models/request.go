package models

type ScrapeRequest struct {
	Terms string `json:"terms,omitempty"`
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
	Sort  string `json:"sort"`
}
