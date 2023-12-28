package internal

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	"github.com/m3db/prometheus_remote_client_golang/promremote"
)

const (
	remoteWriteInterval = 5 * time.Second
)

// Sink that collects metrics and pushes the metrics
// to a remote-write endpoint every n seconds from the channel
type Metric struct {
	Property string `json:"p"`
	Metadata string `json:"m"`
}

type Sink struct {
	metricsMap  map[string]map[string]float64
	metricsChan <-chan Metric
	promClient  promremote.Client
	authString  string
}

func NewSink(metricsChan <-chan Metric, config Config) *Sink {
	metricsMap := map[string]map[string]float64{
		// Pre-initialize the top-level properties with an allow-list so that we can better
		// scope the metrics we care about (since the ingest endpoint will be publicly accessible)
		// Realistically, we should also pre-initialize the metadata but we'll let it go for now
		// and keep an eye on cardinality.
		"proxy.cubari.moe": {},
		"cubari.moe":       {},
	}

	client, err := promremote.NewClient(promremote.NewConfig(
		promremote.WriteURLOption(config.GetRemoteWriteURL()),
	))
	if err != nil {
		log.Fatal("unable to construct client", err)
	}

	authString := base64.StdEncoding.EncodeToString([]byte(config.GetRemoteWriteUsername() + ":" + config.GetRemoteWritePassword()))

	return &Sink{
		metricsMap:  metricsMap,
		metricsChan: metricsChan,
		promClient:  client,
		authString:  authString,
	}
}

func (s *Sink) Start() {
	timer := time.NewTicker(remoteWriteInterval)
	for {
		select {
		// Handle metrics
		case metric := <-s.metricsChan:
			s.handleMetric(metric)

		// Export metrics
		case <-timer.C:
			s.exportMetrics()
		}
	}
}

func (s *Sink) handleMetric(metric Metric) {
	propMetadata, exists := s.metricsMap[metric.Property]
	if !exists {
		log.Println("Property not found", metric.Property)
		return
	}
	propMetadata[metric.Metadata] += 1
}

func (s *Sink) exportMetrics() {
	exportTimestamp := time.Now()
	var timeSeriesList []promremote.TimeSeries

	for property, metadatas := range s.metricsMap {
		for metadata, count := range metadatas {
			mts := promremote.TimeSeries{
				Labels: []promremote.Label{
					{
						Name:  "__name__",
						Value: "cubari_property_metadata_total",
					},
					{
						Name:  "property",
						Value: property,
					},
					{
						Name:  "metadata",
						Value: metadata,
					},
				},
				Datapoint: promremote.Datapoint{
					Timestamp: exportTimestamp,
					Value:     count,
				},
			}
			timeSeriesList = append(timeSeriesList, mts)
		}
	}

	ctx := context.Background()
	_, err := s.promClient.WriteTimeSeries(ctx, timeSeriesList, promremote.WriteOptions{
		Headers: map[string]string{
			"Authorization": "Basic " + s.authString,
		},
	})

	if err != nil {
		log.Println("Failed emitting metrics to remote", err)
	}
}
