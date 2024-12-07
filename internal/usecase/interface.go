package usecase

import (
	"meteo/internal/structs"
)

type DatabaseProvider interface {
	INSERTNewPOSTDataMeteo(data *structs.POSTDataMeteo) error
	SELECTCurrentData() (*structs.CurrentData, error)
	SELECTCurrentDayData() (*structs.WeatherData, error)
	SELECTHistoricalData(from string, to string) (*structs.WeatherData, error)
	INSERTNewUser(user *structs.User) error
	SELECTLoginUser(email string) (*structs.User, error)
	SELECTUserByID(userID int) (*structs.User, error)
	SELECTApiKey(apikey string) bool
}
