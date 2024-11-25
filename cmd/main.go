package main

import (
	"database/sql"
	"fmt"
	Router "github.com/Sergey-pixel-dev/router"
	_ "github.com/lib/pq"
	"log"
	"meteo/internal/core"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "bmstu"
)

func main() {
	router := Router.NewRouter()
	serv := core.NewServer("localhost", "8081", *router)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	serv.SetServerDBprovider(core.NewdbProvider(db))

	router.AddRoute(Router.NewRoute("POST", "/api/post", serv.POSTNewDataHandler))
	router.MethodNotAllowedHandler = http.HandlerFunc(core.MethodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(core.NotFoundHandler)
	router.AddMiddleware(core.CORSMiddleware)

	serv.Router = *router

	err = http.ListenAndServe(serv.Ip+":"+serv.Port, &serv.Router)
	if err != nil {
		log.Fatal(err)
	}

}
