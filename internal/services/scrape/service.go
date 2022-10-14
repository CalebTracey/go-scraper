package scrape

import (
	"context"
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/gocolly/colly"
)

type ServiceI interface {
	ScrapeData(ctx context.Context, req models.ScrapeRequest) (dataList []models.Data, err error)
}

type Service struct {
	Collector *colly.Collector
}

func NewService(config *Config) (Service, error) {
	s := &CollyScraper{
		TimeoutSeconds:        config.TimeoutSeconds,
		LoadingTimeoutSeconds: config.LoadingTimeoutSeconds,
		UserAgent:             config.UserAgent,
	}
	collector, err := s.Init()
	if err != nil {
		return Service{}, err
	}
	return Service{Collector: collector}, nil
}

func (s Service) ScrapeData(ctx context.Context, req models.ScrapeRequest) (dataList []models.Data, err error) {
	scrapeUrl := buildScrapeUrl(req)
	s.Collector.OnHTML("div[class=info]", func(h *colly.HTMLElement) {
		data := models.Data{
			Name:          h.ChildText("a.business-name"),
			Ratings:       h.ChildText("span.bbb-rating extra-rating"),
			Phone:         h.ChildText("div.phones"),
			StreetAddress: h.ChildText("div.street-address"),
			Locality:      h.ChildText("div.locality"),
			URL:           h.ChildAttr("a.track-visit-website", "href"),
		}
		h.ForEach("div.categories a", func(i int, he *colly.HTMLElement) {
			data.Categories = append(data.Categories, he.Text)
		})

		dataList = append(dataList, data)
	})

	s.Collector.OnHTML("div.pagination a.next", func(h *colly.HTMLElement) {
		pageLink := h.Request.AbsoluteURL(h.Attr("href"))

		if pageLink != "" {
			err = s.Collector.Visit(pageLink)
			if err != nil {
				return
			}
		}
	})

	if err != nil {
		return nil, err
	}
	err = s.Collector.Visit(scrapeUrl)
	if err != nil {
		return nil, err
	}

	s.Collector.Wait()
	return dataList, nil
}
