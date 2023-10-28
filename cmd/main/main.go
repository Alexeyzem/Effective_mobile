package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"local/EffectiveMobile/config"
	"local/EffectiveMobile/internal/app"
	"log"
)

func main() {
	cfg := config.NewConfig()
	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatal(err)
	}
	a := app.New(cfg)
	if _, err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
