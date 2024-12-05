package api

import (
	"meteo/internal/helper"
	"meteo/internal/structs"
	"net/http"
)

/*
POST / HTTP/1.1
Host: yourserver.com
Content-Type: application/json
Content-Length: 54

 /*{"date":"121124","temp":"243","hum":"99","pres":"759"}*/
/*{"date":"150405020106","temp":"243","hum":"99","pres":"759"}*/
/*{"date":"2006-01-02 15:04:05","temp":"243","hum":"99","pres":"759"} -expected
 */

// data
func (serv *Server) POSTNewDataHandler(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("POSTNewDataHanlder, POST IP: " + r.RemoteAddr)
	POSTMeteoData := structs.POSTDataMeteo{}
	err := helper.ReadJSON(w, r, &POSTMeteoData)
	if err != nil {
		serv.Logger.LogERROR("POSTNewDataHandler, read json error:" + err.Error())
		helper.WriteJSON(w, http.StatusBadRequest, helper.Envelope{"error": "incorrect json"}, nil)
		return
	}
	err = serv.uc.InsertNewDataMeteo(&POSTMeteoData)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, helper.Envelope{"error": "db.Exec"}, nil)
		return
	}
	helper.WriteJSON(w, http.StatusCreated, helper.Envelope{"msg": "OK"}, nil)
}

func (serv *Server) GETCurrentDataHandler(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("GETCurrentDataHandler, GET IP: " + r.RemoteAddr)
	CurData, err := serv.uc.GetCurrentDataMeteo()
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, helper.Envelope{"error": "db error"}, nil)
		return
	}
	helper.WriteJSON(w, http.StatusOK, CurData, nil)
}
func (serv *Server) GETCurrentDayDataHandler(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("GETCurrentDayDataHandler, GET IP: " + r.RemoteAddr)
	data, err := serv.uc.GetCurrentDayDataMeteo()
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, helper.Envelope{"error": "server error"}, nil)
		return
	}
	helper.WriteJSON(w, http.StatusOK, data, nil)
}
func (serv *Server) GETHistoricalDataHandler(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("GETHistoricalDataHandler, GET IP: " + r.RemoteAddr)
	query := r.URL.Query()
	firstDateStr := query.Get("first_date")
	lastDateStr := query.Get("last_date")
	//проверит DateStr
	data, err := serv.uc.GetHistoricalData(firstDateStr, lastDateStr)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, helper.Envelope{"error": "internal error"}, nil)
		return
	}
	helper.WriteJSON(w, http.StatusOK, data, nil)

}

// user
func (serv *Server) POSTRegisterNewUser(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("POSTRegisterNewUser, POST IP: " + r.RemoteAddr)
	usjson := structs.UserJSON{}
	helper.ReadJSON(w, r, &usjson)
	err := serv.uc.RegisterNewUser(usjson.Email, usjson.Password)
	if err != nil {
		serv.Logger.LogERROR("Error usecase Register new user: " + err.Error())
		helper.WriteJSON(w, http.StatusInternalServerError, helper.Envelope{"error": "internal error"}, nil)
	}
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

//other

func (serv *Server) CORSMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func (serv *Server) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	helper.WriteJSON(w, http.StatusMethodNotAllowed, helper.Envelope{"error": "MethodNotAllowed, method: " + r.Method}, headers)
}

func (serv *Server) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	helper.WriteJSON(w, http.StatusNotFound, helper.Envelope{"error": "PathNotFound, path: " + r.URL.Path}, headers)

}
