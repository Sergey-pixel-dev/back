package usecase

import mlg "meteo/internal/libs/logger"

type Usecase struct {
	dbp    DatabaseProvider
	logger *mlg.MyLogger
}

func NewUsecase(dbp DatabaseProvider, logger *mlg.MyLogger) *Usecase {
	return &Usecase{
		dbp:    dbp,
		logger: logger,
	}
}
