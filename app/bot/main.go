package main

import (
	"game-controller-bot/internal"
	"log"
)

var version = "0.0.0"

func main() {
	err := internal.StartBot(version)
	if err != nil {
		log.Fatal(err)
	}
}
