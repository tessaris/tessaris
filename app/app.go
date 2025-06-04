package app

import (
	"github.com/tessaris/tessaris/config"
	"github.com/tessaris/tessaris/router"
)

type App struct {
	Config *config.Config
	Routes router.Routes
}

func New(cfg *config.Config, routes router.Routes) *App {
	return &App{cfg, routes}
}
