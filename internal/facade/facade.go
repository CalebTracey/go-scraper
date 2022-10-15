package facade

import (
	"context"
	config "github.com/calebtracey/config-yaml"
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/calebtracey/go-scraper/internal/services/googlemaps"
	"github.com/calebtracey/go-scraper/internal/services/scrape"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"strings"
)

type ServiceI interface {
	GetData(ctx context.Context, req models.ScrapeRequest) (res models.ScrapeResponse)
}

type Service struct {
	ScrapeService  scrape.ServiceI
	GeocodeService googlemaps.ServiceI
}

func NewService(appConfig *config.Config) (Service, error) {
	scrapeConfig := scrape.NewConfig()
	scrapeSvc, err := scrape.InitializeService(scrapeConfig)
	if err != nil {
		return Service{}, err
	}
	geocodeSvc, err := googlemaps.InitializeService(appConfig)
	if err != nil {
		return Service{}, err
	}

	return Service{
		ScrapeService:  scrapeSvc,
		GeocodeService: geocodeSvc,
	}, nil
}

func (s Service) GetData(ctx context.Context, req models.ScrapeRequest) (res models.ScrapeResponse) {
	var m models.Message
	var g errgroup.Group
	scrapeUrl := scrape.BuildScrapeUrl(req)
	dataList, err := s.ScrapeService.ScrapeData(ctx, scrapeUrl)
	if err != nil {
		m.ErrorLog = errorLogs([]error{err}, "Failed to scrape data", http.StatusInternalServerError)
		m.Status = strconv.Itoa(http.StatusInternalServerError)
		res.Message = m
		return res
	}

	for idx := range dataList {
		i := idx
		g.Go(func() error {
			address := strings.Join([]string{dataList[i].StreetAddress, dataList[i].Locality}, " ")
			loc, locErr := s.GeocodeService.GeocodeLocationAddress(ctx, address)
			if locErr != nil {
				return locErr
			}
			dataList[i].Location.Lat = loc.Lat
			dataList[i].Location.Lng = loc.Lng
			return nil
		})
	}

	if gErr := g.Wait(); gErr != nil {
		m.ErrorLog = errorLogs([]error{gErr}, "Failed to get geocode location", http.StatusInternalServerError)
		m.Status = strconv.Itoa(http.StatusInternalServerError)
		res.Message = m
		return res
	}

	res.Data = dataList
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
