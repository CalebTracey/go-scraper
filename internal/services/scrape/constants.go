package scrape

// yellowpages elements
const (
	ypInfoSection     = "div.info-section"
	ypBusinessName    = "a.business-name"
	ypBusinessRating  = "span.bbb-rating"
	ypBadges          = "div.badges"
	ypYearsInBusiness = "div.years-in-business"
	ypPhones          = "div.phones"
	ypAddress         = "div.adr"
	ypStreetAddress   = "div.street-address"
	ypLocality        = "div.locality"
	ypLinks           = "div.links"
	ypResultSite      = "a.track-visit-website"

	ypAverageRating = "&&s=average_rating"
)

// sorting fields
const (
	ReqSortResultsByName     = "name"
	ReqSortResultsByDistance = "distance"
	ReqSortResultsByRating   = "rating"

	sortResultsByName     = "&s=name"
	sortResultsByDistance = "&s=distance"
	sortResultsByRating   = "&s=average_rating"
)
