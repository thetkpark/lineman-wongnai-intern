package main

import (
	"github.com/thetkpark/lineman-wongnai-intern/covid"
	"github.com/thetkpark/lineman-wongnai-intern/data"
	"github.com/thetkpark/lineman-wongnai-intern/router"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("unable to create logger: %s", err)
	}
	sugaredLogger := logger.Sugar()
	app := router.NewGinRouter("8080", sugaredLogger)

	covidDataStore := data.NewDiskCovidCasesDataStore("covid-case.json")
	covidHandler := covid.NewHandler(covidDataStore)

	app.Get("/", covidHandler.Summarize)
	app.ListenAndServe()()
}
