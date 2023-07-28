package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

type DeploymentsCollectorMetrics struct {
	namespace   string
	environment string
	boshName    string
	boshUUID    string
}

func NewDeploymentsCollectorMetrics(
	namespace string,
	environment string,
	boshName string,
	boshUUID string,
) *DeploymentsCollectorMetrics {
	return &DeploymentsCollectorMetrics{
		namespace:   namespace,
		environment: environment,
		boshName:    boshName,
		boshUUID:    boshUUID,
	}
}

func (m *DeploymentsCollectorMetrics) NewLastDeploymentsScrapeDurationSecondsMetric() prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "",
			Name:      "last_deployments_scrape_duration_seconds",
			Help:      "Duration of the last scrape of Deployments metrics from BOSH.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
	)
}

func (m *DeploymentsCollectorMetrics) NewLastDeploymentsScrapeTimestampMetric() prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "",
			Name:      "last_deployments_scrape_timestamp",
			Help:      "Number of seconds since 1970 since last scrape of Deployments metrics from BOSH.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
	)
}

func (m *DeploymentsCollectorMetrics) NewDeploymentInstancesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "deployment",
			Name:      "instances",
			Help:      "Number of instances in this deployment",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_vm_type"},
	)
}

func (m *DeploymentsCollectorMetrics) NewDeploymentStemcellInfoMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "deployment",
			Name:      "stemcell_info",
			Help:      "Labeled BOSH Deployment Stemcell Info with a constant '1' value.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_stemcell_name", "bosh_stemcell_version", "bosh_stemcell_os_name"},
	)
}

func (m *DeploymentsCollectorMetrics) NewDeploymentReleasePackageInfoMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "deployment",
			Name:      "release_package_info",
			Help:      "Labeled BOSH Deployment Release Package Info with a constant '1' value.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_release_name", "bosh_release_version", "bosh_release_package_name"},
	)
}

func (m *DeploymentsCollectorMetrics) NewDeploymentReleaseJobInfoMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "deployment",
			Name:      "release_job_info",
			Help:      "Labeled BOSH Deployment Release Job Info with a constant '1' value.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_release_name", "bosh_release_version", "bosh_release_job_name"},
	)
}

func (m *DeploymentsCollectorMetrics) NewDeploymentReleaseInfoMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "deployment",
			Name:      "release_info",
			Help:      "Labeled BOSH Deployment Release Info with a constant '1' value.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_release_name", "bosh_release_version"},
	)
}
