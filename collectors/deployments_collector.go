package collectors

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry/bosh_exporter/deployments"
)

type DeploymentsCollector struct {
	deploymentReleaseInfoMetric                *prometheus.GaugeVec
	deploymentReleaseJobInfoMetric             *prometheus.GaugeVec
	deploymentReleasePackageInfoMetric         *prometheus.GaugeVec
	deploymentStemcellInfoMetric               *prometheus.GaugeVec
	deploymentInstancesMetric                  *prometheus.GaugeVec
	lastDeploymentsScrapeTimestampMetric       prometheus.Gauge
	lastDeploymentsScrapeDurationSecondsMetric prometheus.Gauge
}

func NewDeploymentsCollector(
	namespace string,
	environment string,
	boshName string,
	boshUUID string,
) *DeploymentsCollector {
	metrics := NewDeploymentsCollectorMetrics(namespace, environment, boshName, boshUUID)
	collector := &DeploymentsCollector{
		deploymentReleaseInfoMetric:                metrics.NewDeploymentReleaseInfoMetric(),
		deploymentReleaseJobInfoMetric:             metrics.NewDeploymentReleaseJobInfoMetric(),
		deploymentReleasePackageInfoMetric:         metrics.NewDeploymentReleasePackageInfoMetric(),
		deploymentStemcellInfoMetric:               metrics.NewDeploymentStemcellInfoMetric(),
		deploymentInstancesMetric:                  metrics.NewDeploymentInstancesMetric(),
		lastDeploymentsScrapeTimestampMetric:       metrics.NewLastDeploymentsScrapeTimestampMetric(),
		lastDeploymentsScrapeDurationSecondsMetric: metrics.NewLastDeploymentsScrapeDurationSecondsMetric(),
	}
	return collector
}

func (c *DeploymentsCollector) Collect(deployments []deployments.DeploymentInfo, ch chan<- prometheus.Metric) error {
	var begun = time.Now()

	c.deploymentReleaseInfoMetric.Reset()
	c.deploymentReleaseJobInfoMetric.Reset()
	c.deploymentReleasePackageInfoMetric.Reset()
	c.deploymentStemcellInfoMetric.Reset()
	c.deploymentInstancesMetric.Reset()

	for _, deployment := range deployments {
		c.reportDeploymentReleaseInfoMetrics(deployment)
		c.reportDeploymentStemcellInfoMetrics(deployment)
		c.reportDeploymentInstancesMetrics(deployment)
	}

	c.deploymentReleaseInfoMetric.Collect(ch)
	c.deploymentReleaseJobInfoMetric.Collect(ch)
	c.deploymentReleasePackageInfoMetric.Collect(ch)
	c.deploymentStemcellInfoMetric.Collect(ch)
	c.deploymentInstancesMetric.Collect(ch)

	c.lastDeploymentsScrapeTimestampMetric.Set(float64(time.Now().Unix()))
	c.lastDeploymentsScrapeTimestampMetric.Collect(ch)

	c.lastDeploymentsScrapeDurationSecondsMetric.Set(time.Since(begun).Seconds())
	c.lastDeploymentsScrapeDurationSecondsMetric.Collect(ch)

	return nil
}

func (c *DeploymentsCollector) Describe(ch chan<- *prometheus.Desc) {
	c.deploymentReleaseInfoMetric.Describe(ch)
	c.deploymentReleaseJobInfoMetric.Describe(ch)
	c.deploymentReleasePackageInfoMetric.Describe(ch)
	c.deploymentStemcellInfoMetric.Describe(ch)
	c.deploymentInstancesMetric.Describe(ch)
	c.lastDeploymentsScrapeTimestampMetric.Describe(ch)
	c.lastDeploymentsScrapeDurationSecondsMetric.Describe(ch)
}

func (c *DeploymentsCollector) reportDeploymentReleaseInfoMetrics(
	deployment deployments.DeploymentInfo,
) {
	for _, release := range deployment.Releases {
		c.deploymentReleaseInfoMetric.WithLabelValues(
			deployment.Name,
			release.Name,
			release.Version,
		).Set(float64(1))
		for _, jobName := range release.JobNames {
			c.deploymentReleaseJobInfoMetric.WithLabelValues(
				deployment.Name,
				release.Name,
				release.Version,
				jobName,
			).Set(float64(1))
		}
		for _, packageName := range release.PackageNames {
			c.deploymentReleasePackageInfoMetric.WithLabelValues(
				deployment.Name,
				release.Name,
				release.Version,
				packageName,
			).Set(float64(1))
		}
	}
}

func (c *DeploymentsCollector) reportDeploymentStemcellInfoMetrics(
	deployment deployments.DeploymentInfo,
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

func (c *DeploymentsCollector) reportDeploymentInstancesMetrics(
	deployment deployments.DeploymentInfo,
) {
	for _, instance := range deployment.Instances {
		c.deploymentInstancesMetric.WithLabelValues(
			deployment.Name,
			instance.VMType,
		).Add(float64(1))
	}
}
