package routes

import (
	"encoding/json"
	"github.com/calebtracey/go-scraper/internal/facade"
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

var scrapePayload = `{
	"terms": "coffee shops",
    "city": "South Portland",
    "state": "ME",
    "sort": "rating"
}`

var hostName, _ = os.Hostname()

var scrapeResponse = models.ScrapeResponse{
	Data: []models.Data{
		{
			Name:          "Test name",
			Categories:    []string{"fake place"},
			StreetAddress: "123 Fake",
			Locality:      "Address rd",
		},
	},
	Message: models.Message{
		Status:   "SUCCESS",
		HostName: hostName,
		Count:    1,
	},
}

func TestHandler_Scrape(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFacade := facade.NewMockServiceI(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		Service     facade.ServiceI
		requestBody string
		wantCode    int
		wantResp    models.ScrapeResponse
	}{
		{
			name:        "Happy Path",
			Service:     mockFacade,
			requestBody: scrapePayload,
			wantCode:    200,
			wantResp:    scrapeResponse,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/scrape", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			h := Handler{
				Service: tt.Service,
			}

			mockFacade.EXPECT().GetData(gomock.Any(), gomock.Any()).Return(tt.wantResp)
			router := mux.NewRouter()
			h.InitializeRoutes(router)
			router.ServeHTTP(w, r)

			var res models.ScrapeResponse
			err := json.NewDecoder(w.Body).Decode(&res)
			if err != nil {
				t.Errorf("expected json to decode, got err: %v", err.Error())
			}
			if w.Code != tt.wantCode {
				t.Errorf("expected statusCode: %v, got %v", tt.wantCode, w.Code)
			}

			if !reflect.DeepEqual(res.Data, tt.wantResp.Data) {
				t.Errorf("Scrape() = %v, want %v", res.Data, tt.wantResp.Data)
			}
		})
	}
}
