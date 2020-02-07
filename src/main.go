package main

import (
	"fmt"
	"gopkg.in/fsnotify.v1"
	"log"

	"github.com/bengtrj/simple-sync/command"
	"github.com/bengtrj/simple-sync/config"
)

func main() {
	run()
	watchConfig()
}

func run() {
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

	fmt.Println("\n**** Successfully synchronized the configuration across servers ****")
}

func watchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Could not watch file: %v", err)
	}
	defer watcher.Close()
	fmt.Println("\nWatching for config file changes...")
	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Write {
					run()
					fmt.Println("\nWatching for config file changes...")
				}
			case err := <-watcher.Errors:
				log.Fatalf("Could watching file: %v", err)
			}
		}
	}()

	if err := watcher.Add(config.DesiredStateFilePath); err != nil {
		log.Fatalf("Could not add watcher to file: %v", err)
	}

	<-done
}
