package interfaces

import (
	middleware "github.com/Lucas-Linhar3s/JobHub/middlewares"
	"github.com/Lucas-Linhar3s/JobHub/pkg/http/server"
	"github.com/Lucas-Linhar3s/JobHub/pkg/jwt"
	"github.com/Lucas-Linhar3s/JobHub/pkg/log"
	"github.com/gin-gonic/gin"
)

func Routes(
	jwt *jwt.JWT,
	logger *log.Logger,
	server *server.Server,
	handler *AuthHandler,
) []gin.IRoutes {
	group := server.Router.Group("/auth")
	group.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.NoStrictAuth(jwt, logger),
	)

	return []gin.IRoutes{
		group.POST("/", handler.RegisterUser),
		group.POST("/login", handler.LoginWithEmailAndPassword),
		group.GET("/login", handler.RedirectLoginOauth), // /login?oauth_provider={name}
		group.GET("/login/callback", handler.Callback),
	}
}
