package usecase

import (
	mlg "meteo/internal/libs/logger"
	"meteo/internal/libs/ratelimiter"
)

type Usecase struct {
	dbp     DatabaseProvider
	logger  *mlg.MyLogger
	limiter *ratelimiter.RateLimiter
}

func NewUsecase(dbp DatabaseProvider, logger *mlg.MyLogger, limiter *ratelimiter.RateLimiter) *Usecase {
	return &Usecase{
		dbp:     dbp,
		logger:  logger,
		limiter: limiter,
	}
}
