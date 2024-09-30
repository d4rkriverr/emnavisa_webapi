package account

import (
	"emnavisa/webserver/infrastructure/kernel"
	"emnavisa/webserver/infrastructure/utils"
	"net/http"
)

func BuildAccountService(app *kernel.Application) {
	// Create our Handler
	handler := newHandler(NewService(app))
	midd := utils.NewAuthMiddleware(app.Database)
	// midd := utils.AuthMiddleware(app.Database, )
	// Create a sub router for this service
	// router := app.Router.Methods(http.MethodGet).Subrouter()

	// Register our service routes
	app.Router.HandleFunc("POST /api/v2/account/auth", handler.HandleUserLogin)
	// app.Router.HandleFunc("POST /api/v2/account/create", handler.HandleUserCreate)

	app.Router.HandleFunc("GET /api/v2/account/info", midd.Protect(http.HandlerFunc(handler.HandleUserInfo)))
}
