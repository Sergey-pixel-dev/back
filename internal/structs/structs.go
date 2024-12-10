package structs

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

type WeatherData struct {
	LastDate string        `json:"last_date"` // Последний временной штамп
	Data     []WeatherItem `json:"data"`      // Массив данных о погоде
}

type WeatherItem struct {
	Date     string  `json:"date"`
	Temp     float32 `json:"temp"`
	Pressure int     `json:"pressure"`
	Humidity int     `json:"humidity"`
}

type POSTDataMeteo struct {
	Date  string `json:"date"`
	Temp  string `json:"temp"`
	Hum   string `json:"hum"`
	Press string `json:"pres"`
}

type UserJSON struct {
	Email    string `json:"email"` // Email - электронная почта
	Password string `json:"password"`
}

type ChangePasswordJSON struct {
	NewPass string `json:"new_password"` // Email - электронная почта
	OldPass string `json:"old_password"`
}

type User struct {
	ID        int    `psql:"id"`         // ID - уникальный идентификатор пользователя
	Email     string `psql:"email"`      // Email - электронная почта
	Password  string `psql:"password"`   // Password - хэш пароля
	IsActive  bool   `psql:"is_active"`  // IsActive - статус активности пользователя
	Role      string `psql:"role"`       // Role - роль пользователя
	APIKey    string `psql:"api_key"`    // APIKey - ключ API
	CreatedAt string `psql:"created_at"` // CreatedAt - дата и время создания записи
}
