package collectors

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/deployments"
)

type DeploymentsCollector struct {
	deploymentReleaseInfoMetric                *prometheus.GaugeVec
	deploymentStemcellInfoMetric               *prometheus.GaugeVec
	lastDeploymentsScrapeTimestampMetric       prometheus.Gauge
	lastDeploymentsScrapeDurationSecondsMetric prometheus.Gauge
}

func NewDeploymentsCollector(namespace string) *DeploymentsCollector {
	deploymentReleaseInfoMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "deployment",
			Name:      "release_info",
			Help:      "Labeled BOSH Deployment Release Info with a constant '1' value.",
		},
		[]string{"bosh_deployment", "bosh_release_name", "bosh_release_version"},
	)

	deploymentStemcellInfoMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "deployment",
			Name:      "stemcell_info",
			Help:      "Labeled BOSH Deployment Stemcell Info with a constant '1' value.",
		},
		[]string{"bosh_deployment", "bosh_stemcell_name", "bosh_stemcell_version", "bosh_stemcell_os_name"},
	)

	lastDeploymentsScrapeTimestampMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "last_deployments_scrape_timestamp",
			Help:      "Number of seconds since 1970 since last scrape of Deployments metrics from BOSH.",
		},
	)

	lastDeploymentsScrapeDurationSecondsMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "last_deployments_scrape_duration_seconds",
			Help:      "Duration of the last scrape of Deployments metrics from BOSH.",
		},
	)

	collector := &DeploymentsCollector{
		deploymentReleaseInfoMetric:                deploymentReleaseInfoMetric,
		deploymentStemcellInfoMetric:               deploymentStemcellInfoMetric,
		lastDeploymentsScrapeTimestampMetric:       lastDeploymentsScrapeTimestampMetric,
		lastDeploymentsScrapeDurationSecondsMetric: lastDeploymentsScrapeDurationSecondsMetric,
	}
	return collector
}

func (c *DeploymentsCollector) Collect(deployments []deployments.DeploymentInfo, ch chan<- prometheus.Metric) error {
	var begun = time.Now()

	c.deploymentReleaseInfoMetric.Reset()
	c.deploymentStemcellInfoMetric.Reset()

	for _, deployment := range deployments {
		c.reportDeploymentReleaseInfoMetrics(deployment, ch)
		c.reportDeploymentStemcellInfoMetrics(deployment, ch)
	}

	c.deploymentReleaseInfoMetric.Collect(ch)
	c.deploymentStemcellInfoMetric.Collect(ch)

	c.lastDeploymentsScrapeTimestampMetric.Set(float64(time.Now().Unix()))
	c.lastDeploymentsScrapeTimestampMetric.Collect(ch)

	c.lastDeploymentsScrapeDurationSecondsMetric.Set(time.Since(begun).Seconds())
	c.lastDeploymentsScrapeDurationSecondsMetric.Collect(ch)

	return nil
}

func (c *DeploymentsCollector) Describe(ch chan<- *prometheus.Desc) {
	c.deploymentReleaseInfoMetric.Describe(ch)
	c.deploymentStemcellInfoMetric.Describe(ch)
	c.lastDeploymentsScrapeTimestampMetric.Describe(ch)
	c.lastDeploymentsScrapeDurationSecondsMetric.Describe(ch)
}

func (c *DeploymentsCollector) reportDeploymentReleaseInfoMetrics(
	deployment deployments.DeploymentInfo,
	ch chan<- prometheus.Metric,
) {
	for _, release := range deployment.Releases {
		c.deploymentReleaseInfoMetric.WithLabelValues(
			deployment.Name,
			release.Name,
			release.Version,
		).Set(float64(1))
	}
}

func (c *DeploymentsCollector) reportDeploymentStemcellInfoMetrics(
	deployment deployments.DeploymentInfo,
	ch chan<- prometheus.Metric,
) {
	for _, stemcell := range deployment.Stemcells {
		c.deploymentStemcellInfoMetric.WithLabelValues(
			deployment.Name,
			stemcell.Name,
			stemcell.Version,
			stemcell.OSName,
		).Set(float64(1))
	}
}
