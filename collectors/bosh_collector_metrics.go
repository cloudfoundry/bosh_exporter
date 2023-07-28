package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

type BoshCollectorMetrics struct {
	namespace   string
	environment string
	boshName    string
	boshUUID    string
}

func NewBoshCollectorMetrics(
	namespace string,
	environment string,
	boshName string,
	boshUUID string,
) *BoshCollectorMetrics {
	return &BoshCollectorMetrics{
		namespace:   namespace,
		environment: environment,
		boshName:    boshName,
		boshUUID:    boshUUID,
	}
}

func (m *BoshCollectorMetrics) NewLastBoshScrapeDurationSecondsMetric() prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "",
			Name:      "last_scrape_duration_seconds",
			Help:      "Duration of the last scrape from BOSH.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
	)
}

func (m *BoshCollectorMetrics) NewLastBoshScrapeTimestampMetric() prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "",
			Name:      "last_scrape_timestamp",
			Help:      "Number of seconds since 1970 since last scrape from BOSH.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
	)
}

func (m *BoshCollectorMetrics) NewLastBoshScrapeErrorMetric() prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "",
			Name:      "last_scrape_error",
			Help:      "Whether the last scrape of metrics from BOSH resulted in an error (1 for error, 0 for success).",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
	)
}

func (m *BoshCollectorMetrics) NewTotalBoshScrapeErrorsMetric() prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Subsystem: "",
			Name:      "scrape_errors_total",
			Help:      "Total number of times an error occurred scraping BOSH.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
	)
}

func (m *BoshCollectorMetrics) NewTotalBoshScrapesMetric() prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Subsystem: "",
			Name:      "scrapes_total",
			Help:      "Total number of times BOSH was scraped for metrics.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
	)
}
