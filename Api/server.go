package Api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scout/Models"
	scraper "scout/Scraper"
	"sync"

	"github.com/gocolly/colly/v2"
)

var wg sync.WaitGroup

type Server struct {
	listenAddr string
}

func NewServer(port string) *Server {
	return &Server{
		listenAddr: port,
	}
}
func (s *Server) Start() error {
	http.HandleFunc("/", s.handleRequest)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	var requestData Models.Request
	if r.Method != "POST" {
		writeJSON(w, http.StatusNotAcceptable, map[string]any{"error": "Invalid request"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": fmt.Sprintf("Error decoding incoming request data: %v", err)})
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"success": "success"})

	wg.Add(1)
	go func() { // rewrite this to support individual go routine per request sent
		defer wg.Done()
		processRequest(requestData.Data)
	}()
	wg.Wait()

}

func processRequest(data string) {
	scraper := scraper.NewScraper(colly.NewCollector())
	scraper.Init(data)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
