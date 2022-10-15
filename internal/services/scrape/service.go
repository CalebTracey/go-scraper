package scrape

import (
	"context"
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/gocolly/colly"
	"strings"
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
	collector := s.Collector
	subCollector := s.Collector
	collector.OnHTML("div.scrollable-pane", func(h *colly.HTMLElement) {
		h.ForEach("a.business-name", func(i int, he *colly.HTMLElement) {
			url, found := he.DOM.Attr("href")
			if found {
				err = subCollector.Visit(strings.Join([]string{baseUrl, url}, ""))
				if err != nil {
					return
				}
			}
		})
	})

	subCollector.OnHTML("main.container", func(h *colly.HTMLElement) {
		info := h.DOM
		url, _ := info.Find("a.track-visit-website").Attr("href")
		data := models.Data{
			Name:          info.Find("article.business-card").Find("h1.dockable").Text(),
			Ratings:       info.Find("span.bbb-rating extra-rating").Text(),
			Phone:         info.Find("div.phones").Text(),
			StreetAddress: info.Find("div.street-address").Text(),
			Locality:      info.Find("div.locality").Text(),
			URL:           url,
		}
		h.ForEach("div.categories a", func(i int, he1 *colly.HTMLElement) {
			data.Categories = append(data.Categories, he1.Text)
		})

		dataList = append(dataList, data)
	})
	if err != nil {
		return
	}

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
	err = collector.Visit(scrapeUrl)
	if err != nil {
		return nil, err
	}

	collector.Wait()
	return dataList, nil
}
