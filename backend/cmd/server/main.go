package main

import (
	"fmt"

	"github.com/Lucas-Linhar3s/JobHub/di"
	"go.uber.org/zap"
)

// @title Job Board Tech API
// @version 1.0.0
// @description Job Board Tech API è uma API para gerenciamento de vagas de emprego.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	app, cleanUp, err := di.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer cleanUp()
	app.Logger.Info("server start", zap.String("host", fmt.Sprintf("http://%s%s", app.Config.Http.Host, app.Config.Http.Port)))
	app.Logger.Info("docs addr", zap.String("addr", fmt.Sprintf("http://%s%s/swagger/index.html", app.Config.Http.Host, app.Config.Http.Port)))
	if err = app.Server.Run(app.Config); err != nil {
		app.Logger.Fatal("error in init server", zap.Error(err))
	}
}
