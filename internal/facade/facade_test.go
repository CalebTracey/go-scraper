package facade

import (
	"context"
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/calebtracey/go-scraper/internal/services/googlemaps"
	"github.com/calebtracey/go-scraper/internal/services/scrape"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestService_GetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockScrapeService := scrape.NewMockServiceI(ctrl)
	mockGeocodeService := googlemaps.NewMockServiceI(ctrl)

	type fields struct {
		ScrapeService  scrape.ServiceI
		GeocodeService googlemaps.ServiceI
	}
	type args struct {
		ctx context.Context
		req models.ScrapeRequest
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		wantRes             models.ScrapeResponse
		mockAddress         string
		mockScrapeRes       []models.Data
		mockScrapeErrs      []error
		mockGeocodeRes      models.Location
		mockGeocodeErr      error
		expectValidationErr bool
	}{
		{
			name: "Happy Path",
			fields: fields{
				ScrapeService:  mockScrapeService,
				GeocodeService: mockGeocodeService,
			},
			args: args{
				ctx: context.Background(),
				req: models.ScrapeRequest{
					Terms: "Coffee Shops",
					City:  "Portland",
					State: "ME",
					Sort:  "rating",
				},
			},
			wantRes: models.ScrapeResponse{
				Data: []models.Data{
					{
						Name:          "Test name",
						Categories:    []string{"fake place"},
						StreetAddress: "123 Fake",
						Locality:      "Address rd",
						Location: models.Location{
							Lat: 123,
							Lng: 321,
						},
					},
				},
				Message: models.Message{
					ErrorLog: nil,
					Count:    1,
				},
			},
			mockAddress: "123 Fake Address rd",
			mockScrapeRes: []models.Data{
				{
					Name:          "Test name",
					Categories:    []string{"fake place"},
					StreetAddress: "123 Fake",
					Locality:      "Address rd",
				},
			},
			mockGeocodeRes: models.Location{
				Lat: 123,
				Lng: 321,
			},
			mockScrapeErrs:      nil,
			mockGeocodeErr:      nil,
			expectValidationErr: false,
		},
		{
			name: "Sad Path: Validation error",
			fields: fields{
				ScrapeService:  mockScrapeService,
				GeocodeService: mockGeocodeService,
			},
			args: args{
				ctx: context.Background(),
				req: models.ScrapeRequest{
					Terms: "",
					City:  "",
					State: "",
					Sort:  "rating",
				},
			},
			wantRes: models.ScrapeResponse{
				Data: nil,
				Message: models.Message{
					ErrorLog: []models.ErrorLog{
						{
							"400",
							"search terms required",
							"Request error",
						},
						{
							"400",
							"city required",
							"Request error",
						},
						{
							"400",
							"state required",
							"Request error",
						},
					},
					Count: 0,
				},
			},
			mockScrapeErrs:      nil,
			mockGeocodeErr:      nil,
			expectValidationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				ScrapeService:  tt.fields.ScrapeService,
				GeocodeService: tt.fields.GeocodeService,
			}
			scrapeUrl := scrape.BuildScrapeUrl(tt.args.req)
			if !tt.expectValidationErr {
				mockScrapeService.EXPECT().ScrapeCommonData(scrapeUrl).Return(tt.mockScrapeRes, tt.mockScrapeErrs)
				mockGeocodeService.EXPECT().GeocodeLocationAddress(tt.args.ctx, tt.mockAddress).Return(tt.mockGeocodeRes, tt.mockGeocodeErr)
			}
			if gotRes := s.GetData(tt.args.ctx, tt.args.req); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("GetData() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
