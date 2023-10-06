package main

import (
	"fmt"
	"log"
	"scout/Api"
	config "scout/Config"
	"sync"
)

var wg sync.WaitGroup

func main() {
	// Initalize and load configuration
	config := config.NewConfig()
	config.Load()
	// Initialize api server
	server := Api.NewServer(config)
	fmt.Printf("%s %s server running on port%s\n", config.Name, config.Version, config.Port)
	log.Fatal(server.Start())
}
