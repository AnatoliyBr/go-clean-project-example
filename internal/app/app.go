package app

import (
	"flag"
	"log"

	"github.com/AnatoliyBr/go-clean-project-example/internal/controller/http"
	"github.com/AnatoliyBr/go-clean-project-example/internal/repository/counterstore"
	"github.com/AnatoliyBr/go-clean-project-example/internal/usecase"
	"github.com/BurntSushi/toml"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/http.toml", "path to config file")
}

func Run() {
	// Repository
	cs := counterstore.NewCounterStore(map[string]int{"i": 0, "j": 0})

	// Usecase
	uc := usecase.NewUseCase(cs)

	// Controller
	flag.Parse()
	config := http.NewConfig()
	s := http.NewHTTPServer(config, uc)

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	go s.StartCounterManager()

	if err := s.StartHTTPServer(); err != nil {
		log.Fatal(err)
	}
}
