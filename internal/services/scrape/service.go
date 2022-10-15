package scrape

import (
	"context"
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type ServiceI interface {
	ScrapeData(ctx context.Context, scrapeUrl string) (dataList []models.Data, err error)
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

func (s *Service) ScrapeData(ctx context.Context, scrapeUrl string) (dataList []models.Data, err error) {
	collector := s.Collector
	subCollector := s.Collector
	collector.OnHTML("div.scrollable-pane", func(h *colly.HTMLElement) {
		h.ForEach("a.business-name", func(i int, he *colly.HTMLElement) {
			//url, found := he.DOM.Attr("href")
			url := he.Request.AbsoluteURL(he.Attr("href"))
			if url != "" {
				err = subCollector.Visit(url)
				if err != nil {
					log.Error(err)
				}
			}
		})
	})

	subCollector.OnHTML("main.container", func(h *colly.HTMLElement) {
		info := h.DOM
		//url, _ := info.Find("section.details-card").Find("p.website").Attr("href")
		url := h.Request.AbsoluteURL(info.Find("section.details-card").Find("p.website").AttrOr("href", ""))
		data := models.Data{
			Name: info.Find("article.business-card").Find("h1.dockable").Text(),
			//TODO update these for new url
			//Ratings:       info.Find("span.bbb-rating extra-rating").Text(),
			Phone:   info.Find("section.details-card").Find("p.phone").Text(),
			Address: info.Find("section.details-card").Find("p").Text(),
			URL:     url,
		}
		//TODO update for new url
		h.ForEach("div.categories a", func(i int, he *colly.HTMLElement) {
			data.Categories = append(data.Categories, he.Text)
		})

		dataList = append(dataList, data)
	})
	if err != nil {
		log.Error(err)
	}

	collector.OnHTML("div.pagination a.next", func(h *colly.HTMLElement) {
		pageLink := h.Request.AbsoluteURL(h.Attr("href"))

		if pageLink != "" {
			err = collector.Visit(pageLink)
			if err != nil {
				log.Error(err)
			}
		}
	})

	//if err != nil {
	//	return nil, err
	//}
	err = collector.Visit(scrapeUrl)
	if err != nil {
		return nil, err
	}

	collector.Wait()
	subCollector.Wait()
	return dataList, nil
}
