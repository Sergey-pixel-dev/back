package core

import (
	"time"
)

func (dbp *DatabaseProvider) INSERTNewPOSTDataMeteo(MeteoData *POSTDataMeteo) error {
	_, err := dbp.db.Exec("insert into meteo (date, temp, hum, pres, date_esp) values (($1), ($2), ($3), ($4), ($5))",
		time.Now().Format("2006-01-02 15:04:05"), MeteoData.Temp, MeteoData.Hum, MeteoData.Press, MeteoData.Date)
	return err
}
