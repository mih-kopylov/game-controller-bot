package main

import (
	"game-controller-bot/internal"
	"log"
)

func main() {
	err := internal.StartBot()
	if err != nil {
		log.Fatal(err)
	}
}
