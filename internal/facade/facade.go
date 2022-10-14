package facade

import (
	"context"
	config "github.com/calebtracey/config-yaml"
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/calebtracey/go-scraper/internal/services/scrape"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ServiceI interface {
	GetData(ctx context.Context, req models.ScrapeRequest) (res models.ScrapeResponse)
}

type Service struct {
	ScrapeService scrape.ServiceI
}

func NewService(appConfig *config.Config) (Service, error) {
	scrapeConfig := scrape.NewConfig()
	scrapeSvc, err := scrape.NewService(scrapeConfig)
	if err != nil {
		return Service{}, err
	}

	return Service{
		ScrapeService: scrapeSvc,
	}, nil
}

func (s Service) GetData(ctx context.Context, req models.ScrapeRequest) (res models.ScrapeResponse) {
	var m models.Message
	var err error

	res.Data, err = s.ScrapeService.ScrapeData(ctx, req)
	if err != nil {
		m.ErrorLog = errorLogs([]error{err}, "Failed to scrape data", http.StatusInternalServerError)
		m.Status = strconv.Itoa(http.StatusInternalServerError)
		res.Message = m
		return res
	}

	m.Status = strconv.Itoa(http.StatusOK)
	m.Count = len(res.Data)
	res.Message = m

	return res
}

func errorLogs(errors []error, rootCause string, status int) []models.ErrorLog {
	var errLogs []models.ErrorLog
	for _, err := range errors {
		log.Errorf("%v: %v", rootCause, err.Error())
		errLogs = append(errLogs, models.ErrorLog{
			RootCause: rootCause,
			Status:    strconv.Itoa(status),
			Trace:     err.Error(),
		})
	}
	return errLogs
}
