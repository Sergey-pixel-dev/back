package api

import (
	"net/http"
	"strconv"

	mlg "meteo/internal/logger"

	rtr "github.com/Sergey-pixel-dev/router"
)

type Server struct {
	Ip     string
	Port   int
	Router *rtr.Router
	uc     Usecase
	Logger *mlg.MyLogger
}

func NewServer(ip string, port int, logger *mlg.MyLogger, uc Usecase) *Server {
	router := rtr.NewRouter()
	api := &Server{
		Ip:     ip,
		Port:   port,
		Router: router,
		uc:     uc,
		Logger: logger,
	}

	router.AddRoute(rtr.NewRoute("POST", "/data/post/", api.POSTNewDataHandler))
	router.AddRoute(rtr.NewRoute("GET", "/data/current/", api.GETCurrentDataHandler))
	router.AddRoute(rtr.NewRoute("GET", "/data/currentday/", api.GETCurrentDayDataHandler))
	router.AddRoute(rtr.NewRoute("GET", "/data/statistics", api.GETHistoricalDataHandler))
	router.AddRoute(rtr.NewRoute("POST", "/user/register", api.POSTRegisterNewUser))
	router.AddRoute(rtr.NewRoute("POST", "/user/login", api.POSTLoginUser))

	router.MethodNotAllowedHandler = http.HandlerFunc(api.MethodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(api.NotFoundHandler)
	router.AddMiddleware(api.CORSMiddleware)
	return api
}

func (serv *Server) Run() error {
	err := http.ListenAndServe(serv.Ip+":"+strconv.Itoa(serv.Port), serv.Router)
	if err != nil {
		serv.Logger.LogERROR("Run(): " + err.Error())
		return err
	}
	return nil
}

func (serv *Server) Close() {
	//serv.uc.Close()
	serv.Logger.Close()
}
