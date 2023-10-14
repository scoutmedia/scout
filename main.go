package main

import (
	"fmt"
	"log"
	"scout/Api"
	config "scout/Config"
	downloader "scout/Downloader"
)

func main() {
	// Initalize and load configuration
	config := config.NewConfig()
	config.Load()

	// Initiate downloader
	downloader := downloader.NewDownloader(config.DataDir)

	// Initialize api server
	server := Api.NewServer(config, downloader)
	fmt.Printf("%s %s server running on port%s\n", config.Name, config.Version, config.Port)
	log.Fatal(server.Start())
}
