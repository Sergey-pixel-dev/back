package usecase

import (
	"meteo/internal/structs"
)

func (u *Usecase) InsertNewDataMeteo(data *structs.POSTDataMeteo) error {
	err := u.dbp.INSERTNewPOSTDataMeteo(data)
	return err
}

func (u *Usecase) GetCurrentDataMeteo() (*structs.CurrentData, error) {
	cur, err := u.dbp.SELECTCurrentData()
	return cur, err
}

func (u *Usecase) GetCurrentDayDataMeteo() (*structs.WeatherData, error) {
	data, err := u.dbp.SELECTCurrentDayData()
	return data, err
}
func (u *Usecase) GetHistoricalData(from string, to string) (*structs.WeatherData, error) {
	data, err := u.dbp.SELECTHistoricalData(from, to)
	return data, err
}
