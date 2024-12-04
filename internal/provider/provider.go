package provider

import (
	"database/sql"
	"fmt"

	mlg "meteo/internal/logger"
)

type DatabaseProvider struct {
	db     *sql.DB
	logger *mlg.MyLogger
}

func NewDBProvider(host string, port int, user, password, dbname string, logger *mlg.MyLogger) *DatabaseProvider {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.LogERROR("dbp open error: " + err.Error())
	}
	return &DatabaseProvider{db: db, logger: logger}
}

func (dbp *DatabaseProvider) Close() {
	err := dbp.db.Close()
	dbp.logger.Close()
	if err != nil {
		dbp.logger.LogERROR("dbp close: " + err.Error())
	}
}
