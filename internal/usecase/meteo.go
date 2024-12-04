package usecase

type Usecase struct {
	dbp DatabaseProvider
}

func NewUsecase(dbp DatabaseProvider) *Usecase {
	return &Usecase{
		dbp: dbp,
	}
}
