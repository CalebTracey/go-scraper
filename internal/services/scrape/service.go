package scrape

import (
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/gocolly/colly"
)

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

//go:generate mockgen -destination=mockService.go -package=scrape . ServiceI
type ServiceI interface {
	ScrapeCommonData(scrapeUrl string) (dataList []models.Data, errs []error)
}

type Service struct {
	Collector *colly.Collector
}

func InitializeService(config *Config) (*Service, error) {
	s := &CollyScraper{
		TimeoutSeconds:        config.TimeoutSeconds,
		LoadingTimeoutSeconds: config.LoadingTimeoutSeconds,
		UserAgent:             config.UserAgent,
	}
	collector, err := s.Init()
	if err != nil {
		return &Service{}, err
	}
	return &Service{Collector: collector}, nil
}

func (s *Service) ScrapeCommonData(scrapeUrl string) (dataList []models.Data, errs []error) {
	s.Collector.OnHTML("div.v-card", func(h *colly.HTMLElement) {
		info := h.DOM
		data := models.Data{
			Name: info.Find(ypInfoSection).Find(ypBusinessName).Text(),
			Ratings: models.Ratings{
				BBBRating: info.Find(ypInfoSection).Find(ypBusinessRating).Text(),
				//TARatingURL: h.Request.AbsoluteURL(h.ChildAttr("a.ta-rating-wrapper", "href")),
			},
			YearsInBusiness: info.Find(ypInfoSection).Find(ypBadges).Find(ypYearsInBusiness).Find("div.number").Text(),
			Phone:           info.Find(ypInfoSection).Find(ypPhones).Text(),
			StreetAddress:   info.Find(ypInfoSection).Find(ypAddress).Find(ypStreetAddress).Text(),
			Locality:        info.Find(ypInfoSection).Find(ypAddress).Find(ypLocality).Text(),
			URL:             info.Find(ypInfoSection).Find(ypLinks).Find(ypResultSite).AttrOr("href", ""),
			DataUrl:         h.Request.AbsoluteURL(h.ChildAttr(ypBusinessName, "href")),
		}
		h.ForEach("div.categories a", func(i int, he1 *colly.HTMLElement) {
			data.Categories = append(data.Categories, he1.Text)
		})
		dataList = append(dataList, data)
	})

	err := s.Collector.Visit(scrapeUrl)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	s.Collector.Wait()
	return dataList, errs
}
