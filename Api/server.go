package Api

import (
	"encoding/json"
	"fmt"
	"net/http"
	config "scout/Config"
	downloader "scout/Downloader"
	models "scout/Models"
	scraper "scout/Scraper"
	task "scout/Task"
	"sync"

	"github.com/gocolly/colly/v2"
)

var wg sync.WaitGroup

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
	var requestData models.Request

	if r.Method != "POST" {
		writeJSON(w, http.StatusNotAcceptable, map[string]any{"error": "Invalid request"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": fmt.Sprintf("Error decoding incoming request data: %v", err)})
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"success": "success"})
	go s.process(requestData.Data, s.config)
}

func (s *Server) process(data string, config config.Config) {
	newTask := task.NewTask(data, config.Sources)
	scraper := scraper.NewScraper(colly.NewCollector(), newTask, config)
	scraper.Start(s.downloader)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
