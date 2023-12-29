package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	defer func() {
		w.WriteHeader(http.StatusNoContent)
	}()

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	var metric Metric
	if err := json.Unmarshal(bodyBytes, &metric); err == nil {
		ih.metricsChan <- metric
	}
	var timing Timing
	if err := json.Unmarshal(bodyBytes, &timing); err == nil {
		requestTiming := timing.Timing
		countryHeader := r.Header.Get("CF-IPCountry")
		if requestTiming != "" && countryHeader != "" {
			log.Println(fmt.Sprintf("RequestTiming t=%s country=%s", requestTiming, countryHeader))
		}
	}
}
