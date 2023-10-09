package main

import (
	"fmt"
	"log"
	"scout/Api"
	config "scout/Config"
	downloader "scout/Downloader"
	"sync"
)

var wg sync.WaitGroup

func main() {
	// Initalize and load configuration
	config := config.NewConfig()
	config.Load()

	// Initiate downloader

	downloader := downloader.NewDownloader(config.DataDir)
	go downloader.Monitor(downloader.Client)
	// Initialize api server
	server := Api.NewServer(config, downloader)
	fmt.Printf("%s %s server running on port%s\n", config.Name, config.Version, config.Port)
	log.Fatal(server.Start())
	log.Println("here")
}
