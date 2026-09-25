package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "iplatform",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests.",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "iplatform",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Histogram of response latency (seconds).",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path", "status"})
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Не считаем собственный endpoint Prometheus
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// defer — чтобы метрика записалась даже если хендлер паникует
		// и панику поймает Recoverer ниже по стеку
		defer func() {
			duration := time.Since(start).Seconds()
			path := chi.RouteContext(r.Context()).RoutePattern()
			if path == "" {
				path = "unmatched"
			}
			status := strconv.Itoa(ww.status)

			httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
			httpRequestDuration.WithLabelValues(r.Method, path, status).Observe(duration)
		}()

		next.ServeHTTP(ww, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
