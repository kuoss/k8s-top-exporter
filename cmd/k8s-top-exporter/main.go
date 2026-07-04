package main

import (
	"log"
	"net/http"

	"github.com/jmnote/k8s-top-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	metricCollector, err := collector.NewCollector()
	if err != nil {
		log.Fatalf("initialize collector: %v", err)
	}

	prometheus.MustRegister(metricCollector)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html><head><title>k8s top exporter</title></head><body><h1>k8s top exporter</h1><p><a href='/metrics'>Metrics</a></p></body></html>`))
	})
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("Listening on :9977")
	log.Fatal(http.ListenAndServe(":9977", mux))
}
