package main

import (
	"faculty/internal/app"
	"log"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
