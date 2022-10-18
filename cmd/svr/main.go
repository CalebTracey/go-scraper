package main

import (
	"github.com/NYTimes/gziphandler"
	config "github.com/calebtracey/config-yaml"
	"github.com/calebtracey/go-scraper/internal/facade"
	"github.com/calebtracey/go-scraper/internal/routes"
	"github.com/calebtracey/go-scraper/internal/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var (
	configPath = "config.yaml"
)

const Port = "6080"

func main() {
	defer panicQuit()
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}
	appConfig := config.NewFromFile(configPath)
	service, err := facade.NewService(appConfig)
	if err != nil {
		log.Panicln(err)
	}

	handler := routes.Handler{
		Service: service,
	}
	router := mux.NewRouter().StrictSlash(true)
	handler.InitializeRoutes(router)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Access-Control-Allow-Methods", "Access-Control-Allow-Origin", "X-Requested-With", "Authorization", "Content-Type", "X-Requested-With", "Bearer", "Origin"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	log.Fatal(services.ListenAndServe(Port, gziphandler.GzipHandler(c.Handler(router))))
}

func panicQuit() {
	if r := recover(); r != nil {
		log.Errorf("I panicked and am quitting: %v", r)
		log.Error("I should be alerting someone...")
	}
}
