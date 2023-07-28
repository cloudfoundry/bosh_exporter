package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ServiceDiscoveryCollectorMetrics struct {
	namespace   string
	environment string
	boshName    string
	boshUUID    string
}

func NewServiceDiscoveryCollectorMetrics(
	namespace string,
	environment string,
	boshName string,
	boshUUID string,
) *ServiceDiscoveryCollectorMetrics {
	return &ServiceDiscoveryCollectorMetrics{
		namespace:   namespace,
		environment: environment,
		boshName:    boshName,
		boshUUID:    boshUUID,
	}
}
func (m *ServiceDiscoveryCollectorMetrics) NewLastServiceDiscoveryScrapeTimestampMetric() prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "",
			Name:      "last_service_discovery_scrape_timestamp",
			Help:      "Number of seconds since 1970 since last scrape of Service Discovery from BOSH.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
	)
}

func (m *ServiceDiscoveryCollectorMetrics) NewLastServiceDiscoveryScrapeDurationSecondsMetric() prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "",
			Name:      "last_service_discovery_scrape_duration_seconds",
			Help:      "Duration of the last scrape of Service Discovery from BOSH.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
	)
}
