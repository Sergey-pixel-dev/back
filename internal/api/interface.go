package api

import (
	"meteo/internal/libs/mytoken"
	"meteo/internal/structs"
)

type Usecase interface {
	InsertNewDataMeteo(data *structs.POSTDataMeteo) error
	GetCurrentDataMeteo(remoteAddr string, apikey string) (*structs.CurrentData, error)
	GetCurrentDayDataMeteo(remoteAddr string, apikey string) (*structs.WeatherData, error)
	GetHistoricalData(remoteAddr string, apikey string, from string, to string) (*structs.WeatherData, error)
	RegisterNewUser(email string, password string) (*mytoken.Token, *mytoken.Token, error)
	LoginUser(email string, password string) (*mytoken.Token, *mytoken.Token, error)
	GetUserInfo(tokenAccess *mytoken.Token) (*structs.User, error)
	RefreshToken(RToken *mytoken.Token) (*mytoken.Token, *mytoken.Token, error)
}
