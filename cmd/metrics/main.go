package main

import (
	"cubari-metrics/internal"
)

func main() {
	config := internal.NewConfigFromEnv()
	metricsChan := make(chan internal.Metric)
	sink := internal.NewSink(metricsChan, config)
	server := internal.NewServer(metricsChan, config)

	// If we were doing this "right", we'd handle graceful shutdowns to flush metrics
	// to the remote before the service shuts down (ie. signal.Notify). However, these
	// metrics should be high volume and lossy so it shouldn't be that important
	go server.Start()
	go sink.Start()

	for {
	}
}
