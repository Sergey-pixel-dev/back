package api

import (
	"meteo/internal/helper"
	//mylog "meteo/internal/libs/logger"
	"meteo/internal/libs/mytoken"
	"meteo/internal/structs"
	"net/http"
)

func (serv *Server) POSTRegisterNewUser(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("POSTRegisterNewUser, POST IP: " + r.RemoteAddr)
	usjson := structs.UserJSON{}
	helper.ReadJSON(w, r, &usjson)
	accessToken, refreshToken, err := serv.uc.RegisterNewUser(usjson.Email, usjson.Password)
	if err != nil {
		if err.Error() == "already exists" {
			helper.WriteJSON(w, http.StatusOK, helper.Envelope{"error": "already exists"}, nil)
		} else {
			helper.WriteJSON(w, http.StatusInternalServerError, helper.Envelope{"error": "internal error"}, nil)
		}
		return
	}
	//helper.WriteJSON(w, http.StatusOK, helper.Envelope{"msg": "registration succesfull"}, nil)
	accessToken.SendToken(w) //сам заголовок 200 поставит
	refreshToken.SendCookieToken("refresh_token", "/user/login/refresh", w)
}

func (serv *Server) POSTLoginUser(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("POSTLoginUser, POST IP: " + r.RemoteAddr)
	usjson := structs.UserJSON{}
	helper.ReadJSON(w, r, &usjson)
	accessToken, refreshToken, err := serv.uc.LoginUser(usjson.Email, usjson.Password)
	if err == nil {
		accessToken.SendToken(w)
		refreshToken.SendCookieToken("refresh_token", "/user/login/refresh", w)
		return
	}
	if err.Error() == "Incorrect email" || err.Error() == "Wrong password" {
		serv.Logger.LogINFO("POSTLoginUser: incorrect password or email")
		helper.WriteJSON(w, http.StatusUnauthorized, helper.Envelope{"error": "incorrect password or email"}, nil)
		return
	}
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, helper.Envelope{"error": "internal error"}, nil)
		return
	}

}

func (serv *Server) GETUserInfo(w http.ResponseWriter, r *http.Request) {
	rawToken := r.Header.Get("authorization")[7:]
	if rawToken == "" {
		helper.WriteJSON(w, 401, helper.Envelope{"error": "incorrect token"}, nil)
		return
	}
	tokenAccess, err := mytoken.GetToken(rawToken)
	if err != nil {
		helper.WriteJSON(w, 401, helper.Envelope{"error": "incorrect token"}, nil)
		return
	}
	user, err2 := serv.uc.GetUserInfo(tokenAccess)
	if err2 != nil {
		helper.WriteJSON(w, 401, helper.Envelope{"error": "invalid token"}, nil)
		return
	}
	helper.WriteJSON(w, 200, helper.Envelope{"email": user.Email, "api_key": user.APIKey}, nil)
}

func (serv *Server) corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)
}
func (serv *Server) corsMiddlewire(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
