package main

import (
	"log"
	"command-injection-server/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalln(err)
	}
}
