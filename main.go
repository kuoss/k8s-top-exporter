package main

import (
	"github.com/jmnote/k8s-top-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	collector, err := collector.NewCollector()
	if err != nil {
		log.Println("Cannot initialize a new collector.")
	}
	prometheus.Register(collector)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><head><title>kubectl top exporter</title></head><body><h1>kubectl top exporter</h1><p><a href='/metrics'>Metrics</a></p></body></html>`))
	})
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Listening on :9977")
	log.Fatal(http.ListenAndServe(":9977", nil))
}
