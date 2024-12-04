package usecase

import "meteo/internal/structs"

func (u *Usecase) InsertNewDataMeteo(data *structs.POSTDataMeteo) error {
	err := u.dbp.INSERTNewPOSTDataMeteo(data)
	return err
}

func (u *Usecase) GetCurrentDataMeteo() (*structs.CurrentData, error) {
	cur, err := u.dbp.SELECTCurrentData()
	return cur, err
}
