package main

import (
	"log"
	"os"

	"github.com/MelnikovNA/noolingo-api-gw/internal/app"
)

const configPath = "configs/config.yml"

func main() {
	err := app.Run(configPath)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
