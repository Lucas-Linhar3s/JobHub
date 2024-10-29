//go:build wireinject
// +build wireinject

package di

import (
	"github.com/Lucas-Linhar3s/JobHub/database"
	"github.com/Lucas-Linhar3s/JobHub/modules/auth"
	"github.com/Lucas-Linhar3s/JobHub/pkg/config"
	"github.com/Lucas-Linhar3s/JobHub/pkg/http/server"
	"github.com/Lucas-Linhar3s/JobHub/pkg/jwt"
	"github.com/Lucas-Linhar3s/JobHub/pkg/log"
	"github.com/google/wire"
)

type App struct {
	Server *server.Server
	Config *config.Config
	Logger *log.Logger
}

func newApp(
	server *server.Server,
	config *config.Config,
	logger *log.Logger,

	authModule *auth.Module,
) *App {
	return &App{
		Server: server,
		Config: config,
		Logger: logger,
	}
}

var globalSet = wire.NewSet(
	config.NewViper,
	config.LoadAttributes,
	jwt.NewJwt,
	log.NewLog,
	server.NewServer,
	database.NewDatabase,
	newApp,
)

func InitializeApp() (*App, func(), error) {
	panic(wire.Build(
		globalSet,
		auth.AuthModuleSet,
	))
}
