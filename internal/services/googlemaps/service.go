package googlemaps

import (
	"context"
	"fmt"
	config "github.com/calebtracey/config-yaml"
	"github.com/calebtracey/go-scraper/internal/models"
	"googlemaps.github.io/maps"
	"os"
)

//go:generate mockgen -destination=mockService.go -package=googlemaps . ServiceI
type ServiceI interface {
	GeocodeLocationAddress(ctx context.Context, address string) (loc models.Location, err error)
}

type Service struct {
	Client *maps.Client
}

func InitializeService(config *config.Config) (*Service, error) {
	geoCodingSvc, err := config.GetServiceConfig("GEOCODING")
	if err != nil {
		return nil, err
	}
	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv(geoCodingSvc.ApiKeyEnvironmentVariable.Value)))
	if err != nil {
		return nil, err
	}
	return &Service{
		Client: client,
	}, nil
}

func (s Service) GeocodeLocationAddress(ctx context.Context, address string) (loc models.Location, err error) {
	r := &maps.GeocodingRequest{
		Address: address,
	}
	res, err := s.Client.Geocode(ctx, r)
	if err != nil || len(res) == 0 {
		return loc, fmt.Errorf("res Geocode err: %v", err)
	}

	loc.Lat = res[0].Geometry.Location.Lat
	loc.Lng = res[0].Geometry.Location.Lng

	return loc, nil
}
