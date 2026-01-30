package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	Requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"path"},
	)
)

func Register() {
	prometheus.MustRegister(Requests)
}
