package usecase

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"meteo/internal/libs/mytoken"
	"meteo/internal/structs"
	"time"
)

func (u *Usecase) RegisterNewUser(email string, password string) (*mytoken.Token, *mytoken.Token, error) {
	u.logger.LogINFO("Attempt to register new user: " + email)
	us, err2 := u.dbp.SELECTLoginUser(email)
	if err2 != nil {
		return nil, nil, err2
	}
	if us != nil {
		u.logger.LogINFO("Users with " + email + " already exists")
		return nil, nil, errors.New("already exists")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.LogERROR("bcrypt GenerateFromPassword: " + err.Error())
		return nil, nil, err
	}
	NewUser := structs.User{ //подумать над абстрагированием user создать отдельно его методы, потом
		Email:    email,
		Password: string(hashPassword),
		IsActive: true,
		Role:     "user",
		APIKey:   uuid.New().String(),
	}
	err = u.dbp.INSERTNewUser(&NewUser)
	if err != nil {
		return nil, nil, err
	}
	accessToken, _ := mytoken.NewToken(map[string]interface{}{"alg": "HS256", "typ": "JWT"},
		map[string]interface{}{"userid": NewUser.ID, "exp": time.Now().Add(time.Hour).Unix()},
		"meteo",
	)
	refreshToken, _ := mytoken.NewToken(map[string]interface{}{"alg": "HS256", "typ": "JWT"},
		map[string]interface{}{"userid": NewUser.ID, "exp": time.Now().Add(72 * time.Hour).Unix()},
		"meteo",
	)
	u.logger.LogINFO("New user " + email + " has been successfully registered")
	return accessToken, refreshToken, nil

}

func (u *Usecase) LoginUser(email string, password string) (*mytoken.Token, *mytoken.Token, error) {
	u.logger.LogINFO("Attempt to login from user: " + email)
	user, err := u.dbp.SELECTLoginUser(email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.New("Incorrect email")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil, errors.New("Wrong password")
	}
	accessToken, _ := mytoken.NewToken(map[string]interface{}{"alg": "HS256", "typ": "JWT"},
		map[string]interface{}{"userid": user.ID, "exp": time.Now().Add(time.Hour).Unix()},
		"meteo",
	)
	refreshToken, _ := mytoken.NewToken(map[string]interface{}{"alg": "HS256", "typ": "JWT"},
		map[string]interface{}{"userid": user.ID, "exp": time.Now().Add(24 * time.Hour).Unix()},
		"meteo",
	)
	u.logger.LogINFO("User " + email + " logged in successfully")
	return accessToken, refreshToken, nil

}

func (u *Usecase) GetUserInfo(tokenAccess *mytoken.Token) (*structs.User, error) {
	if !tokenAccess.VerifyToken(func(payload map[string]interface{}) bool {
		exp, ok := payload["exp"].(float64)
		if !ok {
			return false
		}
		return int64(exp) > time.Now().Unix()
	}, "meteo") {
		return nil, errors.New("invalid token")
	}
	user, err := u.dbp.SELECTUserByID(int(tokenAccess.Payload["userid"].(float64)))
	return user, err
}
