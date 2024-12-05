package api

import (
	"meteo/internal/structs"
)

type Usecase interface {
	InsertNewDataMeteo(data *structs.POSTDataMeteo) error
	GetCurrentDataMeteo() (*structs.CurrentData, error)
	GetCurrentDayDataMeteo() (*structs.WeatherData, error)
	GetHistoricalData(from string, to string) (*structs.WeatherData, error)
	RegisterNewUser(email string, password string) error
	LoginUser(email string, password string) error
}
