package routes

import (
	"encoding/json"
	"github.com/calebtracey/go-scraper/internal/facade"
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Handler struct {
	Service facade.ServiceI
}

func (h Handler) InitializeRoutes() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Handle("/scrape", h.Scrape()).Methods(http.MethodPost)

	return r
}

func (h Handler) Scrape() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res models.ScrapeResponse
		var req models.ScrapeRequest
		startTime := time.Now()

		defer func() {
			res.Message.TimeTaken = time.Since(startTime).String()
			res, status := setScrapeResponse(res)
			_ = json.NewEncoder(writeHeader(w, status)).Encode(res)
		}()

		reqBody, readErr := io.ReadAll(r.Body)

		if readErr != nil {
			res.Message.ErrorLog = errorLogs([]error{readErr}, "Unable to read request body", http.StatusBadRequest)
			return
		}
		err := json.Unmarshal(reqBody, &req)
		if err != nil {
			res.Message.ErrorLog = errorLogs([]error{err}, "Unable to parse request", http.StatusBadRequest)
			return
		}

		res = h.Service.GetData(r.Context(), req)

	}
}

func setScrapeResponse(res models.ScrapeResponse) (models.ScrapeResponse, int) {
	status, _ := strconv.Atoi(res.Message.Status)
	hn, _ := os.Hostname()
	return models.ScrapeResponse{
		Data: res.Data,
		Message: models.Message{
			ErrorLog: res.Message.ErrorLog,
			HostName: hn,
			Status:   res.Message.Status,
			Count:    res.Message.Count,
		},
	}, status
}

func writeHeader(w http.ResponseWriter, code int) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	return w
}

func errorLogs(errors []error, rootCause string, status int) []models.ErrorLog {
	var errLogs []models.ErrorLog
	for _, err := range errors {
		errLogs = append(errLogs, models.ErrorLog{
			RootCause: rootCause,
			Status:    strconv.Itoa(status),
			Trace:     err.Error(),
		})
	}
	return errLogs
}
