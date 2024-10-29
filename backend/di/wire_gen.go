// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/Lucas-Linhar3s/JobHub/database"
	"github.com/Lucas-Linhar3s/JobHub/modules/auth"
	"github.com/Lucas-Linhar3s/JobHub/modules/auth/application"
	"github.com/Lucas-Linhar3s/JobHub/modules/auth/interfaces"
	"github.com/Lucas-Linhar3s/JobHub/pkg/config"
	"github.com/Lucas-Linhar3s/JobHub/pkg/http/server"
	"github.com/Lucas-Linhar3s/JobHub/pkg/jwt"
	"github.com/Lucas-Linhar3s/JobHub/pkg/log"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitializeApp() (*App, func(), error) {
	serverServer := server.NewServer()
	viper := config.NewViper()
	configConfig := config.LoadAttributes(viper)
	logger := log.NewLog(configConfig)
	jwtJWT := jwt.NewJwt(configConfig)
	databaseDatabase := database.NewDatabase(configConfig, logger)
	authApp := application.NewAuthApp(logger, databaseDatabase, configConfig, jwtJWT)
	authHandler := interfaces.NewAuthHandler(authApp)
	module := auth.ModuleRegister(jwtJWT, logger, serverServer, authHandler)
	app := newApp(serverServer, configConfig, logger, module)
	return app, func() {
	}, nil
}

// wire.go:

type App struct {
	Server *server.Server
	Config *config.Config
	Logger *log.Logger
}

func newApp(server2 *server.Server, config2 *config.Config,
	logger *log.Logger,

	authModule *auth.Module,
) *App {
	return &App{
		Server: server2,
		Config: config2,
		Logger: logger,
	}
}

var globalSet = wire.NewSet(config.NewViper, config.LoadAttributes, jwt.NewJwt, log.NewLog, server.NewServer, database.NewDatabase, newApp)