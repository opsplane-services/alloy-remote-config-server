package config

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Counter for total config requests
	configRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "alloy_config_server_requests_total",
			Help: "Total number of config requests",
		},
		[]string{"template", "status"},
	)

	// Histogram for render duration
	renderDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "alloy_config_server_render_duration_seconds",
			Help:    "Time taken to render configuration templates",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"template"},
	)

	// Gauge for active templates
	activeTemplates = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "alloy_config_server_templates_loaded",
			Help: "Number of templates currently loaded",
		},
	)

	// Counter for template renders per template
	templateRendersTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "alloy_config_server_template_renders_total",
			Help: "Total number of renders per template",
		},
		[]string{"template"},
	)

	// Gauge for last render timestamp
	lastRenderTimestamp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "alloy_config_server_last_render_timestamp",
			Help: "Unix timestamp of last successful render",
		},
		[]string{"template", "id"},
	)
)
