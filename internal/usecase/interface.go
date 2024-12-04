package usecase

import "meteo/internal/structs"

type DatabaseProvider interface {
	INSERTNewPOSTDataMeteo(data *structs.POSTDataMeteo) error
	SELECTCurrentData() (*structs.CurrentData, error)
}
