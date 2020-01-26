package main

import (
	"log"

	"github.com/bengtrj/simple-sync/command"
	"github.com/bengtrj/simple-sync/config"
)

func main() {
	// passd := os.Getenv("PASSWORD")

	c, err := config.Load()
	if err != nil {
		log.Fatal("error loading config")
	}
	err = command.Sync(c)
	if err != nil {
		log.Fatalf("error syncing %v", err)
	}
}
