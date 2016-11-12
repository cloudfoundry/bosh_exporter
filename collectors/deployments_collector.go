package collectors

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/deployments"
)

type DeploymentsCollector struct {
	deploymentReleaseInfoDesc                *prometheus.Desc
	deploymentStemcellInfoDesc               *prometheus.Desc
	lastDeploymentsScrapeTimestampDesc       *prometheus.Desc
	lastDeploymentsScrapeDurationSecondsDesc *prometheus.Desc
}

func NewDeploymentsCollector(namespace string) *DeploymentsCollector {
	deploymentReleaseInfoDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "deployment", "release_info"),
		"Labeled BOSH Deployment Release Info with a constant '1' value.",
		[]string{"bosh_deployment", "bosh_release_name", "bosh_release_version"},
		nil,
	)

	deploymentStemcellInfoDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "deployment", "stemcell_info"),
		"Labeled BOSH Deployment Stemcell Info with a constant '1' value.",
		[]string{"bosh_deployment", "bosh_stemcell_name", "bosh_stemcell_version", "bosh_stemcell_os_name"},
		nil,
	)

	lastDeploymentsScrapeTimestampDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "last_deployments_scrape_timestamp"),
		"Number of seconds since 1970 since last scrape of Deployments metrics from BOSH.",
		[]string{},
		nil,
	)

	lastDeploymentsScrapeDurationSecondsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "last_deployments_scrape_duration_seconds"),
		"Duration of the last scrape of Deployments metrics from BOSH.",
		[]string{},
		nil,
	)

	collector := &DeploymentsCollector{
		deploymentReleaseInfoDesc:                deploymentReleaseInfoDesc,
		deploymentStemcellInfoDesc:               deploymentStemcellInfoDesc,
		lastDeploymentsScrapeTimestampDesc:       lastDeploymentsScrapeTimestampDesc,
		lastDeploymentsScrapeDurationSecondsDesc: lastDeploymentsScrapeDurationSecondsDesc,
	}
	return collector
}

func (c *DeploymentsCollector) Collect(deployments []deployments.DeploymentInfo, ch chan<- prometheus.Metric) error {
	var begun = time.Now()

	for _, deployment := range deployments {
		c.reportDeploymentReleaseInfoMetrics(deployment, ch)
		c.reportDeploymentStemcellInfoMetrics(deployment, ch)
	}

	ch <- prometheus.MustNewConstMetric(
		c.lastDeploymentsScrapeTimestampDesc,
		prometheus.GaugeValue,
		float64(time.Now().Unix()),
	)

	ch <- prometheus.MustNewConstMetric(
		c.lastDeploymentsScrapeDurationSecondsDesc,
		prometheus.GaugeValue,
		time.Since(begun).Seconds(),
	)

	return nil
}

func (c *DeploymentsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.deploymentReleaseInfoDesc
	ch <- c.deploymentStemcellInfoDesc
	ch <- c.lastDeploymentsScrapeTimestampDesc
	ch <- c.lastDeploymentsScrapeDurationSecondsDesc
}

func (c *DeploymentsCollector) reportDeploymentReleaseInfoMetrics(
	deployment deployments.DeploymentInfo,
	ch chan<- prometheus.Metric,
) {
	for _, release := range deployment.Releases {
		ch <- prometheus.MustNewConstMetric(
			c.deploymentReleaseInfoDesc,
			prometheus.GaugeValue,
			float64(1),
			deployment.Name,
			release.Name,
			release.Version,
		)
	}
}

func (c *DeploymentsCollector) reportDeploymentStemcellInfoMetrics(
	deployment deployments.DeploymentInfo,
	ch chan<- prometheus.Metric,
) {
	for _, stemcell := range deployment.Stemcells {
		ch <- prometheus.MustNewConstMetric(
			c.deploymentStemcellInfoDesc,
			prometheus.GaugeValue,
			float64(1),
			deployment.Name,
			stemcell.Name,
			stemcell.Version,
			stemcell.OSName,
		)
	}
}
