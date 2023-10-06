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
	server := Api.NewServer(config.Port)
	fmt.Printf("%s server running on port%s\n", config.Name, config.Port)
	log.Fatal(server.Start())
}
