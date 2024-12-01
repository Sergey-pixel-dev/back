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

type DatabaseRow struct {
	Id       int    `db:"id"`
	Date     string `db:"date"`
	Temp     int    `db:"temp"`
	Humidity int    `db:"hum"`
	Pressure int    `db:"pres"`
	Date_Esp string `db:"date_esp"`
}

func (dbp *DatabaseProvider) INSERTNewPOSTDataMeteo(MeteoData *POSTDataMeteo) error {
	_, err := dbp.db.Exec("insert into meteo (date, temp, hum, pres, date_esp) values (($1), ($2), ($3), ($4), ($5))",
		time.Now().Format("2006-01-02 15:04:05"), MeteoData.Temp, MeteoData.Hum, MeteoData.Press, MeteoData.Date)
	if err != nil {
		dbp.logger.LogERROR("Error insert postdatameteo: " + MeteoData.Date + ", " + MeteoData.Hum + ", " + MeteoData.Temp + ", " + MeteoData.Press + " error: " + err.Error())
	}
	return err
}

func (dbp *DatabaseProvider) SELECTCurrentData() (*CurrentData, error) {
	query_rows, err := dbp.db.Query("SELECT * FROM meteo WHERE DATE(date_esp) = (SELECT DATE(MAX(date_esp)) FROM meteo);")
	if err != nil {
		dbp.logger.LogERROR("Error select currentdata: " + err.Error())
		return nil, err
	}
	var MaxTemp int
	var MinTemp int
	MinTemp = 1000
	MaxTemp = -1000
	var SliceRows []DatabaseRow
	for query_rows.Next() {
		var row DatabaseRow
		if err := query_rows.Scan(&row.Id, &row.Date, &row.Temp, &row.Humidity, &row.Pressure, &row.Date_Esp); err != nil {
			dbp.logger.LogERROR("Error query_rows.Scan(): " + err.Error())
			continue
		}
		MaxTemp = max(MaxTemp, row.Temp)
		MinTemp = min(MinTemp, row.Temp)
		SliceRows = append(SliceRows, row)
	}
	if err = query_rows.Err(); err != nil {
		dbp.logger.LogERROR("Error query_rows.Err(): " + err.Error())
		return nil, err
	}
	LastRow := SliceRows[len(SliceRows)-1]
	CurData := CurrentData{
		LastDate: LastRow.Date_Esp,
		Main: Main{
			Temp:     float32(LastRow.Temp/10) + float32(LastRow.Temp%10)/10.0,
			TempMin:  float32(MinTemp/10) + float32(MinTemp%10)/10.0,
			TempMax:  float32(MaxTemp/10) + float32(MaxTemp%10)/10.0,
			Pressure: LastRow.Pressure,
			Humidity: LastRow.Humidity,
		},
	}
	return &CurData, nil

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
