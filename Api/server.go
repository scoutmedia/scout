package Api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	config "scout/Config"
	downloader "scout/Downloader"
	models "scout/Models"
	scraper "scout/Scraper"
	task "scout/Task"

	"github.com/gocolly/colly/v2"
)

type Server struct {
	listenAddr string
	config     config.Config
	downloader *downloader.Downloader
}

func NewServer(config config.Config, downloader *downloader.Downloader) *Server {
	return &Server{
		listenAddr: config.Port,
		config:     config,
		downloader: downloader,
	}
}
func (s *Server) Start() error {
	http.HandleFunc("/", s.handleRequest)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	var media models.Media
	if r.Method != "POST" {
		writeJSON(w, http.StatusNotAcceptable, map[string]any{"error": "Invalid request"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&media); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": fmt.Sprintf("Error decoding incoming request data: %v", err)})
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"success": "success"})
	log.Printf("%s request recieved", media.Name)
	go s.process(media, s.config)
}

func (s *Server) process(media models.Media, config config.Config) {
	newTask := task.NewTask(media, config.Sources)
	scraper := scraper.NewScraper(colly.NewCollector(), newTask, config)
	scraper.Start(s.downloader)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
