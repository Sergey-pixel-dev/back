package main

import (
	"meteo/internal/core"
	"net/http"
	"os"

	Router "github.com/Sergey-pixel-dev/router"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "bmstu"
)

func main() {
	serv_file_log, _ := os.OpenFile("../log/serv.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	db_file_log, _ := os.OpenFile("../log/db.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	logger_serv := core.NewLogger()
	logger_serv.SetDescriptor(serv_file_log)
	logger_db := core.NewLogger()
	logger_db.SetDescriptor(db_file_log)

	router := Router.NewRouter()
	serv := core.NewServer("localhost", "8081", *router, logger_serv)
	defer serv.Close()

	dbp := core.NewDBProvider(logger_db)
	dbp.DBProviderInit(host, port, user, password, dbname)
	serv.SetServerDBprovider(dbp)

	router.AddRoute(Router.NewRoute("POST", "/api/post", serv.POSTNewDataHandler))
	router.MethodNotAllowedHandler = http.HandlerFunc(core.MethodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(core.NotFoundHandler)
	router.AddMiddleware(core.CORSMiddleware)
	serv.Router = *router

	err := http.ListenAndServe(serv.Ip+":"+serv.Port, &serv.Router)
	if err != nil {
		serv.Logger.LogFATAL("ListenAndServe error : " + err.Error())
	}

}
