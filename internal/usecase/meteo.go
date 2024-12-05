package usecase

import (
	"errors"
	"meteo/internal/structs"

	"golang.org/x/crypto/bcrypt"
)

// ВАЖНО! ДОБАВИТЬ ЛОГГИРОВАНИЕ USECASE
// data
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

// user
// может зарегестрировать пользователя с двумя email одинаковыми
func (u *Usecase) RegisterNewUser(email string, password string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	NewUser := structs.User{ //подумать над абстрагированием user создать отдельно его методы, потом
		Email:    email,
		Password: string(hashPassword),
		IsActive: true,
		Role:     "user",
		APIKey:   "none",
	}
	err = u.dbp.INSERTNewUser(&NewUser)
	return err
}

func (u *Usecase) LoginUser(email string, password string) error {
	user, err := u.dbp.SELECTLoginUser(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("Incorrect email")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New("Incorrect password")
	}
	return nil

}
