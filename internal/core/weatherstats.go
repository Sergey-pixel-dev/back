package core

type CurrentData struct {
	LastDate string `json:"last_date"` // последний запись на сервера о данных
	Main     Main   `json:"main"`
}

// Main содержит информацию о погоде
type Main struct {
	Temp     float32 `json:"temp"`     // текущая температура
	TempMin  float32 `json:"temp_min"` // минимальная температура
	TempMax  float32 `json:"temp_max"` // максимальная температура
	Pressure int     `json:"pressure"` // давление
	Humidity int     `json:"humidity"` // влажность
}
