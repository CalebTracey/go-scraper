package scrape

import (
	"context"
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/gocolly/colly"
)

type ServiceI interface {
	ScrapeData(ctx context.Context, scrapeUrl string) (dataList []models.Data, err error)
}

type Service struct {
	Collector *colly.Collector
}

func InitializeService(config *Config) (Service, error) {
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

func (s Service) ScrapeData(ctx context.Context, scrapeUrl string) (dataList []models.Data, err error) {
	s.Collector.OnHTML("div[class=info]", func(h *colly.HTMLElement) {
		info := h.DOM
		url, _ := info.Find("a.track-visit-website").Attr("href")
		data := models.Data{
			Name:          info.Find("a.business-name").Text(),
			Ratings:       info.Find("span.bbb-rating extra-rating").Text(),
			Phone:         info.Find("div.phones").Text(),
			StreetAddress: info.Find("div.street-address").Text(),
			Locality:      info.Find("div.locality").Text(),
			URL:           url,
		}
		h.ForEach("div.categories a", func(i int, he *colly.HTMLElement) {
			data.Categories = append(data.Categories, he.Text)
		})

		dataList = append(dataList, data)
	})

	//s.Collector.OnHTML("div.pagination a.next", func(h *colly.HTMLElement) {
	//	pageLink := h.Request.AbsoluteURL(h.Attr("href"))
	//
	//	if pageLink != "" {
	//		err = s.Collector.Visit(pageLink)
	//		if err != nil {
	//			return
	//		}
	//	}
	//})

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
