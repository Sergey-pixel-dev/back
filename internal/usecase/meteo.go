package usecase

import (
	"errors"
	"meteo/internal/structs"
	"strings"
)

// data
func (u *Usecase) InsertNewDataMeteo(data *structs.POSTDataMeteo) error {
	err := u.dbp.INSERTNewPOSTDataMeteo(data)
	return err
}

func (u *Usecase) GetCurrentDataMeteo(remoteAddr, apikey string) (*structs.CurrentData, error) {
	if strings.Contains(remoteAddr, "localhost") {
		if u.dbp.SELECTApiKey(apikey) {
			if u.limiter.Allow(apikey) {
				cur, err := u.dbp.SELECTCurrentData()
				return cur, err
			} else {
				//u.logger.LogWARN(remoteAddr + " has rate limit exceeded")
				return nil, errors.New("rate limit exceeded")
			}
		} else {
			return nil, errors.New("invalid api_key")
		}
	}
	cur, err := u.dbp.SELECTCurrentData()
	return cur, err

}

func (u *Usecase) GetCurrentDayDataMeteo(remoteAddr, apikey string) (*structs.WeatherData, error) {
	if strings.Contains(remoteAddr, "localhost") {
		if u.dbp.SELECTApiKey(apikey) {
			if u.limiter.Allow(apikey) {
				data, err := u.dbp.SELECTCurrentDayData()
				return data, err
			} else {
				return nil, errors.New("rate limit exceeded")
			}
		} else {
			return nil, errors.New("invalid api_key")
		}
	}
	data, err := u.dbp.SELECTCurrentDayData()
	return data, err
}
func (u *Usecase) GetHistoricalData(remoteAddr, apikey, from, to string) (*structs.WeatherData, error) {
	if strings.Contains(remoteAddr, "localhost") {
		if u.dbp.SELECTApiKey(apikey) {
			if u.limiter.Allow(apikey) {
				data, err := u.dbp.SELECTHistoricalData(from, to)
				return data, err
			} else {
				return nil, errors.New("rate limit exceeded")
			}
		} else {
			return nil, errors.New("invalid api_key")
		}
	}
	data, err := u.dbp.SELECTHistoricalData(from, to)
	return data, err
}
