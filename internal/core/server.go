package core

import "database/sql"
import rtr "github.com/Sergey-pixel-dev/router"

type Server struct {
	Ip         string
	Port       string
	Router     rtr.Router
	DBProvider DatabaseProvider
	Logger     *MyLogger //указатель из-за mutex
}

func NewdbProvider(db *sql.DB) *DatabaseProvider {
	return &DatabaseProvider{db: db}
}
func (s *Server) SetServerDBprovider(dp *DatabaseProvider) {
	s.DBProvider = *dp
}

func NewServer(ip string, port string, handler rtr.Router, logger *MyLogger) *Server {
	return &Server{
		Ip:     ip,
		Port:   port,
		Router: handler,
		Logger: logger,
	}
}

func (serv *Server) Close() {
	serv.DBProvider.Close()
	serv.Logger.Close()
}
