package scrape

import (
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/gocolly/colly"
)

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
		dataUrl := h.Request.AbsoluteURL(h.ChildAttr("a.business-name", "href"))
		data := models.Data{
			Name:            info.Find("div.info-section").Find("a.business-name").Text(),
			Ratings:         info.Find("div.info-section").Find("span.bbb-rating").Text(),
			YearsInBusiness: info.Find("div.info-section").Find("div.badges").Find("div.years-in-business").Find("div.number").Text(),
			Phone:           info.Find("div.info-section").Find("div.phones").Text(),
			StreetAddress:   info.Find("div.info-section").Find("div.adr").Find("div.street-address").Text(),
			Locality:        info.Find("div.info-section").Find("div.adr").Find("div.locality").Text(),
			URL:             info.Find("div.info-section").Find("div.links").Find("a.track-visit-website").AttrOr("href", ""),
			DataUrl:         dataUrl,
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
