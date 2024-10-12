package metrics

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func MostStartMetrics(addr string) error {
	mux.NewRouter().Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(addr, nil)
}
