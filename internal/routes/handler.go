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

func (h Handler) InitializeRoutes(r *mux.Router) {

	r.Handle("/scrape", h.Scrape()).Methods(http.MethodPost)
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
	msg, sc := getResponseStatus(res.Message.ErrorLog, res.Message.Count)
	hn, _ := os.Hostname()
	return models.ScrapeResponse{
		Data: res.Data,
		Message: models.Message{
			ErrorLog: res.Message.ErrorLog,
			HostName: hn,
			Status:   msg,
			Count:    res.Message.Count,
		},
	}, sc
}

func getResponseStatus(errs []models.ErrorLog, lengthOfResults int) (msg string, status int) {
	if len(errs) > 0 {
		msg = "ERROR"
		var s500, s400, s404, s206, s200 bool
		for _, e := range errs {
			code, _ := strconv.Atoi(e.StatusCode)
			switch {
			case code == 206:
				s206 = true
			case code == 404:
				if lengthOfResults > 0 {
					s206 = true
				} else {
					s404 = true
				}
			case code == 400:
				if lengthOfResults > 0 {
					s206 = true
				} else {
					s400 = true
				}
			case code >= 500:
				if lengthOfResults > 0 {
					s206 = true
				} else {
					s500 = true
				}
			default:
				s500 = true
			}
		}
		switch {
		case s206:
			status = http.StatusPartialContent
		case s500:
			status = http.StatusInternalServerError
		case s400:
			status = http.StatusBadRequest
		case s404:
			status = http.StatusNotFound
		case s200:
			status = http.StatusOK
		default:
			status = http.StatusInternalServerError
		}
		return msg, status
	}

	return "SUCCESS", http.StatusOK
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
			RootCause:  rootCause,
			StatusCode: strconv.Itoa(status),
			Trace:      err.Error(),
		})
	}
	return errLogs
}
