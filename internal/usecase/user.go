package usecase

import (
	"errors"
	"meteo/internal/structs"

	"golang.org/x/crypto/bcrypt"
)

func (u *Usecase) RegisterNewUser(email string, password string) error {
	us, err2 := u.dbp.SELECTLoginUser(email)
	if err2 != nil {
		return err2
	}
	if us != nil {
		u.logger.LogINFO("Users with" + email + "already exists")
		return errors.New("already exists")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.LogERROR("bcrypt GenerateFromPassword: " + err.Error())
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
		u.logger.LogINFO("No registered user with " + email)
		return errors.New("Incorrect email")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		u.logger.LogINFO("Wrong password " + email)
		return errors.New("Wrong password")
	}
	return nil

}
