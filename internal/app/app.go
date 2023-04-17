package app

import (
	"log"

	"github.com/AnatoliyBr/go-clean-project-example/internal/controller/http"
	"github.com/AnatoliyBr/go-clean-project-example/internal/repository/counterstore"
	"github.com/AnatoliyBr/go-clean-project-example/internal/usecase"
)

func Run() {
	// Repository
	cs := counterstore.NewCounterStore(map[string]int{"i": 0, "j": 0})

	// Usecase
	uc := usecase.NewUseCase(cs)

	// Controller
	config := http.NewConfig()
	s := http.NewHTTPServer(config, uc)

	go s.StartCounterManager()

	if err := s.StartHTTPServer(); err != nil {
		log.Fatal(err)
	}
}
