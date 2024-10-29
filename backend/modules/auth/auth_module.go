package auth

import (
	"github.com/Lucas-Linhar3s/JobHub/modules/auth/application"
	"github.com/Lucas-Linhar3s/JobHub/modules/auth/interfaces"
	"github.com/Lucas-Linhar3s/JobHub/pkg/http/server"
	"github.com/Lucas-Linhar3s/JobHub/pkg/jwt"
	"github.com/Lucas-Linhar3s/JobHub/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// AuthModuleSet is a variable that represents the auth module set
var AuthModuleSet = wire.NewSet(
	application.NewAuthApp,
	interfaces.NewAuthHandler,
	ModuleRegister,
)

// Module is a struct that represents the auth module
type Module struct {
	Routes []gin.IRoutes
}

// ModuleRegister is a function that registers the auth module
func ModuleRegister(jwt *jwt.JWT, logger *log.Logger, server *server.Server, handler *interfaces.AuthHandler) *Module {
	return &Module{
		Routes: interfaces.Routes(jwt, logger, server, handler),
	}
}
