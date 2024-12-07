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
	api_key := r.URL.Query().Get("api_key")
	CurData, err := serv.uc.GetCurrentDataMeteo(r.RemoteAddr, api_key)
	if err == nil {
		helper.WriteJSON(w, http.StatusOK, CurData, nil)
		return
	}
	if err.Error() == "rate limit exceeded" {
		helper.WriteJSON(w, http.StatusTooManyRequests, helper.Envelope{"msg": "rate limit exceeded"}, nil)
		return
	}
	if err.Error() == "invalid api_key" {
		helper.WriteJSON(w, http.StatusUnauthorized, helper.Envelope{"msg": "invalid api_key"}, nil)
		return
	}
	helper.WriteJSON(w, http.StatusBadRequest, helper.Envelope{"error": "db error"}, nil)
	return

}
func (serv *Server) GETCurrentDayDataHandler(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("GETCurrentDayDataHandler, GET IP: " + r.RemoteAddr)
	api_key := r.URL.Query().Get("api_key")
	data, err := serv.uc.GetCurrentDayDataMeteo(r.RemoteAddr, api_key)
	if err == nil {
		helper.WriteJSON(w, http.StatusOK, data, nil)
		return
	}
	if err.Error() == "rate limit exceeded" {
		helper.WriteJSON(w, http.StatusTooManyRequests, helper.Envelope{"msg": "rate limit exceeded"}, nil)
		return
	}
	if err.Error() == "invalid api_key" {
		helper.WriteJSON(w, http.StatusUnauthorized, helper.Envelope{"msg": "invalid api_key"}, nil)
		return
	}
	helper.WriteJSON(w, http.StatusBadRequest, helper.Envelope{"error": "db error"}, nil)
	return
}
func (serv *Server) GETHistoricalDataHandler(w http.ResponseWriter, r *http.Request) {
	serv.Logger.LogINFO("GETHistoricalDataHandler, GET IP: " + r.RemoteAddr)
	query := r.URL.Query()
	firstDateStr := query.Get("first_date")
	lastDateStr := query.Get("last_date")
	api_key := r.URL.Query().Get("api_key")
	//проверит DateStr
	data, err := serv.uc.GetHistoricalData(r.RemoteAddr, api_key, firstDateStr, lastDateStr)
	if err == nil {
		helper.WriteJSON(w, http.StatusOK, data, nil)
		return
	}
	if err.Error() == "rate limit exceeded" {
		helper.WriteJSON(w, http.StatusTooManyRequests, helper.Envelope{"msg": "rate limit exceeded"}, nil)
		return
	}
	if err.Error() == "invalid api_key" {
		helper.WriteJSON(w, http.StatusUnauthorized, helper.Envelope{"msg": "invalid api_key"}, nil)
		return
	}
	helper.WriteJSON(w, http.StatusBadRequest, helper.Envelope{"error": "db error"}, nil)
	return

}

//other

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
