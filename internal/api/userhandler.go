package api

import (
	"meteo/internal/helper"
	"meteo/internal/structs"
	"net/http"
)

func (serv *Server) POSTRegisterNewUser(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("POSTRegisterNewUser, POST IP: " + r.RemoteAddr)
	usjson := structs.UserJSON{}
	helper.ReadJSON(w, r, &usjson)
	accessToken, refreshToken, err := serv.uc.RegisterNewUser(usjson.Email, usjson.Password)
	if err.Error() == "already exists" {
		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"error": "already exists"}, nil)
		return
	}
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, helper.Envelope{"error": "internal error"}, nil)
	}
	accessToken.SendToken(w)
	refreshToken.SendCookieToken("refresh_token", "/user/login/refresh", w)
	helper.WriteJSON(w, http.StatusOK, helper.Envelope{"msg": "registration succesfull"}, nil)
}

func (serv *Server) POSTLoginUser(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("POSTLoginUser, POST IP: " + r.RemoteAddr)
	usjson := structs.UserJSON{}
	helper.ReadJSON(w, r, &usjson)
	err := serv.uc.LoginUser(usjson.Email, usjson.Password)
	if err == nil {
		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"msg": "login is succesful"}, nil)
		return
	}
	if err.Error() == "Incorrect email" || err.Error() == "Incorrect password" {
		serv.Logger.LogINFO("POSTLoginUser: incorrect password or email")
		helper.WriteJSON(w, http.StatusUnauthorized, helper.Envelope{"error": "incorrect password or email"}, nil)
		return
	}
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, helper.Envelope{"error": "internal error"}, nil)
		return
	}

}
