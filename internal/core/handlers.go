package core

import (
	"fmt"
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
		fmt.Println("read json error:", err.Error()) // добавить log свой заместо print
		err2 := writeJSON(w, http.StatusInternalServerError, envelope{"error": "MethodNotAllowed, method: " + r.Method}, nil)
		if err2 != nil {
			fmt.Println(err.Error())
		}
		return
	}
	err = serv.dbProvider.INSERTNewPOSTDataMeteo(&POSTMeteoData)
	if err != nil {
		fmt.Println("error, INSERNewPOSTDataMeteo, db.Exec()", err.Error())
		writeJSON(w, http.StatusInternalServerError, envelope{"error": "db.Exec"}, nil)
		return
	}
	writeJSON(w, http.StatusCreated, envelope{"msg": "OK"}, nil)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	err := writeJSON(w, http.StatusMethodNotAllowed, envelope{"error": "MethodNotAllowed, method: " + r.Method}, headers)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	err := writeJSON(w, http.StatusNotFound, envelope{"error": "PathNotFound, path: " + r.URL.Path}, headers)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func CORSMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
