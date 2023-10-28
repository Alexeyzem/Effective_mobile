package app

import (
	"local/EffectiveMobile/config"
	"local/EffectiveMobile/internal/logger"
	"local/EffectiveMobile/internal/service"
	storage "local/EffectiveMobile/internal/storage/postgres"
	"local/EffectiveMobile/internal/transport/http"
	v1 "local/EffectiveMobile/internal/transport/http/v1"
	postgres "local/EffectiveMobile/pkg/client"
)

type App struct {
	Config *config.Config
}

func New(config *config.Config) *App {
	return &App{Config: config}
}

func (a *App) Run() (func(), error) {
	log, cleanup, err := logger.New(a.Config.Logger)
	if err != nil {
		return func() {}, err
	}

	postgresClient, cleanup2, err := postgres.NewClient(a.Config.Postgres.URL)
	if err != nil {
		cleanup()
		return func() {}, err
	}

	peopleStorage := storage.NewPeoplePostgres(postgresClient)

	peopleService := service.NewPeople(peopleStorage)

	v1PeopleController := v1.NewPeopleConroller(peopleService, log, nil)
	v1Controller := v1.NewController(v1PeopleController)

	server := http.NewServer(a.Config.Server, v1Controller)

	cleanup3, err := server.Run()
	if err != nil {
		cleanup2()
		cleanup()
		return func() {}, err
	}

	return func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
