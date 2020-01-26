package main

import (
	"log"

	"github.com/bengtrj/simple-sync/command"
	"github.com/bengtrj/simple-sync/config"
)

func main() {

	c, err := config.Load()
	if err != nil {
		log.Fatalf("An error occurred trying to parse the config: %v", err)
	}

	err = command.Sync(c)
	if err != nil {
		log.Fatalf("An error occurred trying to syncronize: %v", err)
	}

	err = config.SetKnownState(c.DesiredState)
	if err != nil {
		log.Fatalf("An error occurred trying update state file: %v", err)
	}

}
