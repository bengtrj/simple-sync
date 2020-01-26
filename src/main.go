package main

import (
	"fmt"
	"log"

	"github.com/bengtrj/simple-sync/config"
)

func main() {
	// passd := os.Getenv("PASSWORD")

	c, err := config.Load()
	if err != nil {
		log.Fatal("error loading config")
	}
	fmt.Println(c.DesiredState.Servers[0].IP)
}
