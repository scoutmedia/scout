package main

import (
	"fmt"
	"log"
	"scout/Api"
	config "scout/Config"
	downloader "scout/Downloader"
	logger "scout/Logger"
)

func main() {
	// Initalize and load configuration
	config := config.NewConfig()
	config.Load()

	// Initiate Logger
	logger := logger.NewLogger()
	logger.Init(config.Sentry.Dsn)
	// Initiate downloader
	downloader := downloader.NewDownloader(config.DataDir, logger)
	// Initialize api server
	server := Api.NewServer(config, downloader)
	logger.Info("Sever Start", fmt.Sprintf("%s %s server running on port%s\n", config.Name, config.Version, config.Port))
	log.Fatal(server.Start())
}
