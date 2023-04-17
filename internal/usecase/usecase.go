package usecase

type Repository interface {
	Set(string, int) int
	Get(string) int
	Inc(string) int
	Dec(string) int
}

type CounterUseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *CounterUseCase {
	return &CounterUseCase{
		repo: repo,
	}
}

func (uc *CounterUseCase) Set(name string, val int) int {
	return uc.repo.Set(name, val)
}

func (uc *CounterUseCase) Get(name string) int {
	return uc.repo.Get(name)
}

func (uc *CounterUseCase) Inc(name string) int {
	return uc.repo.Inc(name)
}

func (uc *CounterUseCase) Dec(name string) int {
	return uc.repo.Dec(name)
}
