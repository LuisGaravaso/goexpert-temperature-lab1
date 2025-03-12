package web

import (
	"encoding/json"
	"net/http"
	locationGateway "temperatures/internal/gateway/location"
	weatherGateway "temperatures/internal/gateway/weather"
	usecase "temperatures/internal/usecase/get_weather"

	"github.com/go-chi/chi/v5"
)

type Error422 struct {
	Error string `json:"error" example:"invalid location"`
}

type Error404 struct {
	Error string `json:"error" example:"location not found"`
}

type Error500 struct {
	Error string `json:"error" example:"internal server error"`
}

type WebWeatherHandler struct {
	LocationGateway locationGateway.LocationGateway
	WeatherGateway  weatherGateway.WeatherGateway
}

func NewWebWeatherHandler(
	locationGateway locationGateway.LocationGateway,
	weatherGateway weatherGateway.WeatherGateway,
) *WebWeatherHandler {
	return &WebWeatherHandler{
		LocationGateway: locationGateway,
		WeatherGateway:  weatherGateway,
	}
}

// Get godoc
// @Summary Get Weather by Location
// @Description Returns weather info by location (CEP or lat,lng)
// @Tags Weather
// @Accept json
// @Produce json
// @Param location path string true "Location (CEP or lat,lng)"
// @Success 200 {object} usecase.GetWeatherOutputDTO
// @Failure 404 {object} Error404
// @Failure 422 {object} Error422
// @Failure 500 {object} Error500
// @Router /temperature/{location} [get]
func (h *WebWeatherHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the coordinates from the URL
	location := chi.URLParam(r, "location")

	var dto usecase.GetWeatherInputDTO
	dto.Location = location

	// Create the use case
	useCase := usecase.NewGetWeatherUseCase(h.LocationGateway, h.WeatherGateway)
	output, err := useCase.Execute(dto)

	// If the location is invalid, return a 422
	if err != nil && err.Error() == "Must be in the format 01001001 for CEP or -23.55028,-46.63389 for Coordinates" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid location"})
		return
	}

	// If the location is not found, return a 404
	if output.Coordinates == "" {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "location not found"})
		return
	}

	// Return the weather information
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusOK)
}
