package provider

import (
	"fmt"
	"meteo/internal/structs"
	"time"
)

type DatabaseRow struct {
	Id       int    `db:"id"`
	Date     string `db:"date"`
	Temp     int    `db:"temp"`
	Humidity int    `db:"hum"`
	Pressure int    `db:"pres"`
	Date_Esp string `db:"date_esp"`
}

func (dbp *DatabaseProvider) INSERTNewPOSTDataMeteo(MeteoData *structs.POSTDataMeteo) error {
	_, err := dbp.db.Exec("insert into meteo (date, temp, hum, pres, date_esp) values (($1), ($2), ($3), ($4), ($5))",
		time.Now().Format("2006-01-02 15:04:05"), MeteoData.Temp, MeteoData.Hum, MeteoData.Press, MeteoData.Date)
	if err != nil {
		dbp.logger.LogERROR("Error insert postdatameteo: " + MeteoData.Date + ", " + MeteoData.Hum + ", " + MeteoData.Temp + ", " + MeteoData.Press + " error: " + err.Error())
	}
	return err
}

func (dbp *DatabaseProvider) SELECTCurrentData() (*structs.CurrentData, error) {
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
	CurData := structs.CurrentData{
		LastDate: LastRow.Date_Esp,
		Main: structs.Main{
			Temp:     float32(LastRow.Temp/10) + float32(LastRow.Temp%10)/10.0,
			TempMin:  float32(MinTemp/10) + float32(MinTemp%10)/10.0,
			TempMax:  float32(MaxTemp/10) + float32(MaxTemp%10)/10.0,
			Pressure: LastRow.Pressure,
			Humidity: LastRow.Humidity,
		},
	}
	return &CurData, nil

}

func (dbp *DatabaseProvider) SELECTCurrentDayData() (*structs.WeatherData, error) {
	query_rows, err := dbp.db.Query("SELECT * FROM meteo WHERE DATE(date_esp) = (SELECT DATE(MAX(date_esp)) FROM meteo);")
	if err != nil {
		dbp.logger.LogERROR("Error select currentdata: " + err.Error())
		return nil, err
	}
	var items []structs.WeatherItem
	for query_rows.Next() {
		var row DatabaseRow
		if err := query_rows.Scan(&row.Id, &row.Date, &row.Temp, &row.Humidity, &row.Pressure, &row.Date_Esp); err != nil {
			dbp.logger.LogERROR("Error query_rows.Scan(): " + err.Error())
			continue
		}
		items = append(items, structs.WeatherItem{
			Date:     row.Date_Esp,
			Temp:     float32(row.Temp/10) + float32(row.Temp%10)/10.0,
			Humidity: row.Humidity,
			Pressure: row.Pressure,
		})
	}
	if err = query_rows.Err(); err != nil {
		dbp.logger.LogERROR("Error query_rows.Err(): " + err.Error())
		return nil, err
	}
	CurDayData := structs.WeatherData{
		LastDate: items[len(items)-1].Date,
		Data:     items,
	}
	return &CurDayData, nil

}

func (dbp *DatabaseProvider) SELECTHistoricalData(from string, to string) (*structs.WeatherData, error) {
	//query_rows, err := dbp.db.Query("SELECT * FROM meteo WHERE date_esp BETWEEN '2024-12-01 00:00:00' AND '2024-12-05 23:59:59';")
	fmt.Println(from)
	fmt.Println(to)
	query_rows, err := dbp.db.Query("SELECT * FROM meteo WHERE DATE(date_esp) BETWEEN $1 AND $2;", from, to)
	if err != nil {
		dbp.logger.LogERROR("Error select currentdata: " + err.Error())
		return nil, err
	}
	var LastTime string
	var items []structs.WeatherItem
	for query_rows.Next() {
		var row DatabaseRow
		if err := query_rows.Scan(&row.Id, &row.Date, &row.Temp, &row.Humidity, &row.Pressure, &row.Date_Esp); err != nil {
			dbp.logger.LogERROR("Error query_rows.Scan(): " + err.Error())
			continue
		}
		items = append(items, structs.WeatherItem{
			Date:     row.Date_Esp,
			Temp:     float32(row.Temp/10) + float32(row.Temp%10)/10.0,
			Humidity: row.Humidity,
			Pressure: row.Pressure,
		})
		LastTime = row.Date_Esp
	}
	if err = query_rows.Err(); err != nil {
		dbp.logger.LogERROR("Error query_rows.Err(): " + err.Error())
		return nil, err
	}
	data := structs.WeatherData{}
	if len(items) == 0 {
		data.LastDate = "-"
		data.Data = nil
	} else {
		data.LastDate = LastTime
		data.Data = items
	}
	return &data, nil

}
