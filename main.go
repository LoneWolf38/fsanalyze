package main

import (
	"log"

	"fsanalyze/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
