package provider

import (
	"database/sql"
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

//data

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
	var LastRow *DatabaseRow
	for query_rows.Next() {
		var row DatabaseRow
		if err := query_rows.Scan(&row.Id, &row.Date, &row.Temp, &row.Humidity, &row.Pressure, &row.Date_Esp); err != nil {
			dbp.logger.LogERROR("Error query_rows.Scan(): " + err.Error())
			continue
		}
		MaxTemp = max(MaxTemp, row.Temp)
		MinTemp = min(MinTemp, row.Temp)
		LastRow = &row
	}
	if err = query_rows.Err(); err != nil {
		dbp.logger.LogERROR("Error query_rows.Err(): " + err.Error())
		return nil, err
	}

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

	var items []structs.WeatherItem //возможно, (из-за того, что часто добавляются данные, и мы по ним не итерируемся) лучше использовать list
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

	query_rows, err := dbp.db.Query("SELECT * FROM meteo WHERE DATE(date_esp) BETWEEN $1 AND $2;", from, to)
	if err != nil {
		dbp.logger.LogERROR("Error select currentdata: " + err.Error())
		return nil, err
	}
	var LastTime string
	var items []structs.WeatherItem //возможно, (из-за того, что часто добавляются данные, и мы по ним не итерируемся) лучше использовать list
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

// user
func (dbp *DatabaseProvider) SELECTApiKey(apikey string) bool {
	row := dbp.db.QueryRow(`select id from users where api_key = ($1);`, apikey)
	err := row.Scan()
	return err != sql.ErrNoRows
}

func (dbp *DatabaseProvider) INSERTNewUser(user *structs.User) error {
	err := dbp.db.QueryRow(`INSERT INTO users (email, password, is_active, role, api_key) 
VALUES (($1), ($2), ($3), ($4), ($5)) RETURNING id;`, user.Email, user.Password, user.IsActive, user.Role, user.APIKey).Scan(&user.ID)
	if err != nil {
		dbp.logger.LogERROR("Error insertnewuser " + err.Error())
		return err
	}
	return nil
}

func (dbp *DatabaseProvider) SELECTLoginUser(email string) (*structs.User, error) {
	row := dbp.db.QueryRow(`select * from users where email = ($1);`, email)
	us := structs.User{}
	err := row.Scan(&us.ID, &us.Email, &us.Password, &us.IsActive, &us.Role, &us.APIKey, &us.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil //переписать
	} else if err != nil {
		dbp.logger.LogERROR("Error SELECTLoginUser: " + err.Error())
		return nil, err
	}
	return &us, nil
}

func (dbp *DatabaseProvider) SELECTUserByID(userID int) (*structs.User, error) {
	row := dbp.db.QueryRow(`select * from users where id = ($1);`, userID)
	us := structs.User{}
	err := row.Scan(&us.ID, &us.Email, &us.Password, &us.IsActive, &us.Role, &us.APIKey, &us.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil //переписать
	} else if err != nil {
		dbp.logger.LogERROR("Error SELECTUserByID: " + err.Error())
		return nil, err
	}
	return &us, nil
}

func (dbp *DatabaseProvider) UPDATEPassword(userID int, password string) error {
	_, err := dbp.db.Exec(`update users set password = ($1) where id = ($2)`, password, userID)
	if err != nil {
		dbp.logger.LogERRORIfExists(err, "Error UPDATEPassword: ")
	}
	return err
}
