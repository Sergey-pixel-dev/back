package api

import "meteo/internal/structs"

type Usecase interface {
	InsertNewDataMeteo(data *structs.POSTDataMeteo) error
	GetCurrentDataMeteo() (*structs.CurrentData, error)
	GetCurrentDayDataMeteo() (*structs.WeatherData, error)
}
