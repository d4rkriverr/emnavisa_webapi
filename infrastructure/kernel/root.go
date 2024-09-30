package kernel

import (
	"emnavisa/webserver/infrastructure/db"
	"emnavisa/webserver/infrastructure/utils"
	"net/http"
)

func Boot() (*Application, error) {

	router := http.NewServeMux()
	database, err := db.NewPostgresDB()
	if err != nil {
		return nil, err
	}

	return &Application{
		Server: &http.Server{
			Addr:    ":8080",
			Handler: utils.CORS(router),
		},
		Database: database,
		Router:   router,
	}, nil
}
