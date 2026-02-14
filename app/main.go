package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Region    string    `json:"region"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

type InfoResponse struct {
	Service     string            `json:"service"`
	Version     string            `json:"version"`
	Region      string            `json:"region"`
	Environment string            `json:"environment"`
	Uptime      string            `json:"uptime"`
	Metadata    map[string]string `json:"metadata"`
}

var startTime = time.Now()

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/readiness", handleReadiness)
	http.HandleFunc("/info", handleInfo)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Starting server on port %s", port)
	log.Printf("Region: %s", getRegion())
	log.Printf("Version: %s", getVersion())

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer recordMetrics(r.Method, "/", http.StatusOK, start)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Multi-Region GCP Deployment Demo",
		"region":  getRegion(),
		"version": getVersion(),
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer recordMetrics(r.Method, "/health", http.StatusOK, start)

	w.Header().Set("Content-Type", "application/json")
	response := HealthResponse{
		Status:    "healthy",
		Region:    getRegion(),
		Version:   getVersion(),
		Timestamp: time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer recordMetrics(r.Method, "/readiness", http.StatusOK, start)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer recordMetrics(r.Method, "/info", http.StatusOK, start)

	uptime := time.Since(startTime)
	
	w.Header().Set("Content-Type", "application/json")
	response := InfoResponse{
		Service:     "multi-region-demo",
		Version:     getVersion(),
		Region:      getRegion(),
		Environment: getEnvironment(),
		Uptime:      uptime.String(),
		Metadata: map[string]string{
			"go_version": "1.21",
			"build_date": getBuildDate(),
		},
	}
	json.NewEncoder(w).Encode(response)
}

func recordMetrics(method, endpoint string, status int, start time.Time) {
	duration := time.Since(start).Seconds()
	httpRequestsTotal.WithLabelValues(method, endpoint, fmt.Sprintf("%d", status)).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

func getRegion() string {
	region := os.Getenv("REGION")
	if region == "" {
		return "unknown"
	}
	return region
}

func getVersion() string {
	version := os.Getenv("VERSION")
	if version == "" {
		return "1.0.0"
	}
	return version
}

func getEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		return "development"
	}
	return env
}

func getBuildDate() string {
	date := os.Getenv("BUILD_DATE")
	if date == "" {
		return time.Now().Format("2006-01-02")
	}
	return date
}
