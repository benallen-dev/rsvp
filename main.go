package main

import (
	"log"
)

func main() {
	cfg := LoadConfig()
	if err := Start(cfg); err != nil {
		log.Fatal(err)
	}
}
