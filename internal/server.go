package internal

import (
	"encoding/json"
	"net/http"
)

type ingestHandler struct {
	metricsChan chan<- Metric
}

type Server struct {
	port string
	mux  *http.ServeMux
}

// Server handles collecting metrics and emitting them to
// a channel that will be batch written through remote write
func NewServer(metricsChan chan<- Metric, config Config) *Server {
	mux := http.NewServeMux()
	ingestHandler := ingestHandler{metricsChan: metricsChan}
	mux.HandleFunc("/ingest", ingestHandler.ingest)
	return &Server{
		port: ":" + config.GetPort(),
		mux:  mux,
	}
}

func (s *Server) Start() {
	http.ListenAndServe(s.port, s.mux)
}

func (ih *ingestHandler) ingest(w http.ResponseWriter, r *http.Request) {
	var metric Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err == nil {
		ih.metricsChan <- metric
	}
	w.WriteHeader(http.StatusNoContent)
}
