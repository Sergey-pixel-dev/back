package core

import (
	"database/sql"
	"fmt"
	"time"
)

type DatabaseProvider struct {
	db     *sql.DB
	logger *MyLogger
}

func (dbp *DatabaseProvider) INSERTNewPOSTDataMeteo(MeteoData *POSTDataMeteo) error {
	_, err := dbp.db.Exec("insert into meteo (date, temp, hum, pres, date_esp) values (($1), ($2), ($3), ($4), ($5))",
		time.Now().Format("2006-01-02 15:04:05"), MeteoData.Temp, MeteoData.Hum, MeteoData.Press, MeteoData.Date)
	if err != nil {
		dbp.logger.LogERROR("Error insert postdatameteo: " + MeteoData.Date + ", " + MeteoData.Hum + ", " + MeteoData.Temp + ", " + MeteoData.Press)
	}
	return err
}
func (dbp *DatabaseProvider) DBProviderInit(host string, port int, user, password, dbname string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		dbp.logger.LogFATAL(err.Error())
	}
	dbp.db = db
}

func NewDBProvider(logger *MyLogger) *DatabaseProvider {
	return &DatabaseProvider{db: nil, logger: logger}
}

func (dbp *DatabaseProvider) Close() {
	err := dbp.db.Close()
	dbp.logger.Close()
	if err != nil {
		dbp.logger.LogERROR("dbp close: " + err.Error())
	}
}
