package core

import (
	"net/http"
)

/*
POST / HTTP/1.1
Host: yourserver.com
Content-Type: application/json
Content-Length: 54

 /*{"date":"121124","temp":"243","hum":"99","pres":"759"}*/

/*{"date":"2006-01-02 15:04:05","temp":"243","hum":"99","pres":"759"} -expected
 */

func (serv *Server) POSTNewDataHandler(w http.ResponseWriter, r *http.Request) {
	POSTMeteoData := POSTDataMeteo{}
	err := readJSON(w, r, &POSTMeteoData)
	if err != nil {
		serv.Logger.LogERROR("POSTNewDataHandler, read json error:" + err.Error())
		writeJSON(w, http.StatusBadRequest, envelope{"error": "incorrect json"}, nil)
		return
	}
	err = serv.DBProvider.INSERTNewPOSTDataMeteo(&POSTMeteoData)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, envelope{"error": "db.Exec"}, nil)
		return
	}
	writeJSON(w, http.StatusCreated, envelope{"msg": "OK"}, nil)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	writeJSON(w, http.StatusMethodNotAllowed, envelope{"error": "MethodNotAllowed, method: " + r.Method}, headers)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	writeJSON(w, http.StatusNotFound, envelope{"error": "PathNotFound, path: " + r.URL.Path}, headers)

}

func CORSMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
