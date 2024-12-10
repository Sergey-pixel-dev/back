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
			helper.WriteJSON(w, http.StatusUnauthorized, helper.Envelope{"error": "already exists"}, nil)
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
		refreshToken.SendCookieToken("refresh_token", "/user/login/refresh", w)
		accessToken.SendToken(w)
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
	rawToken := r.Header.Get("authorization")
	if len(rawToken) < 7 {
		helper.WriteJSON(w, 401, helper.Envelope{"error": "incorrect token"}, nil)
		return
	}
	rawToken = rawToken[7:]
	tokenAccess, err := mytoken.ParseToken(rawToken)
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

func (serv *Server) POSTChangePassword(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("POSTChangePassword, POST IP: " + r.RemoteAddr)
	usjson := structs.ChangePasswordJSON{}
	helper.ReadJSON(w, r, &usjson)
	rawToken := r.Header.Get("authorization")
	if len(rawToken) < 7 {
		helper.WriteJSON(w, 401, helper.Envelope{"error": "incorrect token"}, nil)
		return
	}
	rawToken = rawToken[7:]
	tokenAccess, err := mytoken.ParseToken(rawToken)
	if err != nil {
		helper.WriteJSON(w, 401, helper.Envelope{"error": "incorrect token"}, nil)
		return
	}
	err = serv.uc.ChangePassword(usjson.OldPass, usjson.NewPass, tokenAccess)
	if err == nil {
		helper.WriteJSON(w, 200, helper.Envelope{"msg": "password has been changed"}, nil)
		return
	}
	if err.Error() == "incorrect password" {
		helper.WriteJSON(w, 401, helper.Envelope{"error": "incorrect password"}, nil)
		return
	}
	if err.Error() == "invalid token" {
		helper.WriteJSON(w, 401, helper.Envelope{"error": "invalid token"}, nil)
		return
	}
	helper.WriteJSON(w, 401, helper.Envelope{"error": "internal error"}, nil)

}

func (serv *Server) RefreshTokenHanlder(w http.ResponseWriter, r *http.Request) {
	cookieToken, err := r.Cookie("refresh_token")
	if err != nil {
		helper.WriteJSON(w, 401, helper.Envelope{"error": "incorrect token23"}, nil)
		return
	}
	refreshToken, err := mytoken.GetCookieToken(cookieToken)
	if err != nil {
		helper.WriteJSON(w, 401, helper.Envelope{"error": "incorrect token12342"}, nil)
		return
	}
	accessToken, refreshToken, err := serv.uc.RefreshToken(refreshToken)
	if err == nil {
		accessToken.SendToken(w)
		refreshToken.SendCookieToken("refresh_token", "/user/login/refresh", w)
		return
	} else {
		helper.WriteJSON(w, 401, helper.Envelope{"error": "invalid token"}, nil)
		return
	}
}

func (serv *Server) corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Cookie, Content-Type, Authorization, Set-Cookie")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)
}
func (serv *Server) corsMiddlewire(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Cookie, Content-Type, Authorization, Set-Cookie")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
