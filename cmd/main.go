package main

import (
	"emnavisa/webserver/infrastructure/kernel"
	"emnavisa/webserver/registry/account"
	"emnavisa/webserver/registry/calllog"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	app, err := kernel.Boot()
	if err != nil {
		log.Fatalf("[X] - Cannot boot: %v", err)
	}

	account.BuildAccountService(app)
	calllog.BuildCallsService(app)

	go app.Run()
	kernel.WaitForShutdown(app)
}
