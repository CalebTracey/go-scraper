package scrape

import (
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/gocolly/colly"
	"reflect"
	"testing"
)

func TestService_ScrapeCommonData(t *testing.T) {
	type fields struct {
		Collector *colly.Collector
	}
	type args struct {
		scrapeUrl string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantDataList []models.Data
		wantErrs     []error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Collector: tt.fields.Collector,
			}
			gotDataList, gotErrs := s.ScrapeCommonData(tt.args.scrapeUrl)
			if !reflect.DeepEqual(gotDataList, tt.wantDataList) {
				t.Errorf("ScrapeCommonData() gotDataList = %v, want %v", gotDataList, tt.wantDataList)
			}
			if !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("ScrapeCommonData() gotErrs = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}
