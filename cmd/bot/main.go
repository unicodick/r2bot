package main

import (
	"log"

	"github.com/unicodick/r2bot/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	application.Run()
}
