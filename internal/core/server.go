package core

import "database/sql"
import rtr "github.com/Sergey-pixel-dev/router"

type Server struct {
	Ip         string
	Port       string
	Router     rtr.Router
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

func NewdbProvider(db *sql.DB) *DatabaseProvider {
	return &DatabaseProvider{db: db}
}
func (s *Server) SetServerDBprovider(dp *DatabaseProvider) {
	s.dbProvider = *dp
}

func NewServer(ip string, port string, handler rtr.Router) *Server {
	return &Server{
		Ip:     ip,
		Port:   port,
		Router: handler,
	}
}
