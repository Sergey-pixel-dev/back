package main

import (
	"fmt"
	"meteo/internal/api"
	"meteo/internal/config"
	"meteo/internal/provider"
	"meteo/internal/usecase"
	"os"

	mlg "meteo/internal/libs/logger"
	"meteo/internal/libs/ratelimiter"

	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.LoadConfig("../configs/meteo.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}

	serv_file_log, _ := os.OpenFile("../log/serv.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	db_file_log, _ := os.OpenFile("../log/db.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	uc_file_log, _ := os.OpenFile("../log/uc.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	logger_serv := mlg.NewLogger()
	logger_serv.SetDescriptor(serv_file_log)
	logger_db := mlg.NewLogger()
	logger_db.SetDescriptor(db_file_log)
	logger_uc := mlg.NewLogger()
	logger_uc.SetDescriptor(uc_file_log)

	rt := ratelimiter.NewRateLimiter(5, 5000)

	dbp := provider.NewDBProvider(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBname, logger_db)
	use := usecase.NewUsecase(dbp, logger_uc, rt)
	srv := api.NewServer(cfg.IP, cfg.Port, logger_serv, use)
	srv.Run()
}
