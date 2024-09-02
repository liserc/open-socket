package prommetrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"net/http"
)

const commonPath = "/metrics"

var (
	baseCollector = []prometheus.Collector{
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	}
)

func Init(registry *prometheus.Registry, prometheusPort int, path string, handler http.Handler, cs ...prometheus.Collector) error {
	registry.MustRegister(cs...)
	srv := http.NewServeMux()
	srv.Handle(path, handler)
	return http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), srv)
}
