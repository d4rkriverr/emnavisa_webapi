package calllog

import (
	"emnavisa/webserver/infrastructure/kernel"
	"emnavisa/webserver/infrastructure/utils"
	"net/http"
)

func BuildCallsService(app *kernel.Application) {
	// Create our Handler
	handler := newHandler(NewService(app))
	midd := utils.NewAuthMiddleware(app.Database)

	app.Router.HandleFunc("GET /api/v2/calls/all", midd.Protect(http.HandlerFunc(handler.GetAllCalls)))
	app.Router.HandleFunc("POST /api/v2/calls/create", midd.Protect(http.HandlerFunc(handler.CreateCall)))
	app.Router.HandleFunc("POST /api/v2/calls/update", midd.Protect(http.HandlerFunc(handler.EditCallLog)))
}
