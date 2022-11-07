package main

import (
	"example-grpc-auth/config"
	"example-grpc-auth/server"
	"log"
	"os"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	app := server.NewApp()

	if err := app.Run(os.Getenv("APP_PORT")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
