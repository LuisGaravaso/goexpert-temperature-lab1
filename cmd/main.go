// @title Weather API
// @version 1.0
// @description API for retrieving weather information by location
// @BasePath /

package main

import (
	"log"
	"os"

	_ "temperatures/docs"
	locationGateway "temperatures/internal/gateway/location"
	weatherGateway "temperatures/internal/gateway/weather"
	"temperatures/internal/infra/web"
	"temperatures/internal/infra/web/webserver"
)

func main() {
	locationGateway := &locationGateway.AwesomeAPILocationGateway{}
	weatherGateway := &weatherGateway.WeatherAPIGateway{}

	// Porta dinâmica para compatibilidade com Cloud Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// WebServer
	webWeatherHandler := web.NewWebWeatherHandler(locationGateway, weatherGateway)
	webServer := webserver.NewWebServer(":" + port)
	webServer.AddHandler("/temperature/{location}", webWeatherHandler.Get)

	log.Println("Starting web server on port", ":"+port)
	webServer.Start()
}
