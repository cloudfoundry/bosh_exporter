package collectors

import (
	"sync"
	"time"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/cloudfoundry-community/bosh_exporter/filters"
)

type DeploymentsCollector struct {
	namespace                                string
	deploymentsFilter                        filters.DeploymentsFilter
	deploymentReleaseInfoDesc                *prometheus.Desc
	deploymentStemcellInfoDesc               *prometheus.Desc
	lastDeploymentsScrapeTimestampDesc       *prometheus.Desc
	lastDeploymentsScrapeDurationSecondsDesc *prometheus.Desc
}

func NewDeploymentsCollector(
	namespace string,
	deploymentsFilter filters.DeploymentsFilter,
) *DeploymentsCollector {
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
		namespace:                                namespace,
		deploymentsFilter:                        deploymentsFilter,
		deploymentReleaseInfoDesc:                deploymentReleaseInfoDesc,
		deploymentStemcellInfoDesc:               deploymentStemcellInfoDesc,
		lastDeploymentsScrapeTimestampDesc:       lastDeploymentsScrapeTimestampDesc,
		lastDeploymentsScrapeDurationSecondsDesc: lastDeploymentsScrapeDurationSecondsDesc,
	}
	return collector
}

func (c DeploymentsCollector) Collect(ch chan<- prometheus.Metric) {
	var begun = time.Now()

	deployments := c.deploymentsFilter.GetDeployments()

	var wg sync.WaitGroup
	for _, deployment := range deployments {
		wg.Add(1)
		go func(deployment director.Deployment, ch chan<- prometheus.Metric) {
			defer wg.Done()
			c.reportDeploymentMetrics(deployment, ch)
		}(deployment, ch)
	}
	wg.Wait()

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
}

func (c DeploymentsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.deploymentReleaseInfoDesc
	ch <- c.deploymentStemcellInfoDesc
	ch <- c.lastDeploymentsScrapeTimestampDesc
	ch <- c.lastDeploymentsScrapeDurationSecondsDesc
}

func (c DeploymentsCollector) reportDeploymentMetrics(
	deployment director.Deployment,
	ch chan<- prometheus.Metric,
) {
	c.reportDeploymentReleaseInfoMetrics(deployment, ch)
	c.reportDeploymentStemcellInfoMetrics(deployment, ch)
}

func (c DeploymentsCollector) reportDeploymentReleaseInfoMetrics(
	deployment director.Deployment,
	ch chan<- prometheus.Metric,
) {
	log.Debugf("Reading Releases info for deployment `%s`:", deployment.Name())
	releases, err := deployment.Releases()
	if err != nil {
		log.Errorf("Error while reading Release info for deployment `%s`: %v", deployment.Name(), err)
		return
	}

	for _, release := range releases {
		ch <- prometheus.MustNewConstMetric(
			c.deploymentReleaseInfoDesc,
			prometheus.GaugeValue,
			float64(1),
			deployment.Name(),
			release.Name(),
			release.Version().AsString(),
		)
	}
}

func (c DeploymentsCollector) reportDeploymentStemcellInfoMetrics(
	deployment director.Deployment,
	ch chan<- prometheus.Metric,
) {
	log.Debugf("Reading Stemcells info for deployment `%s`:", deployment.Name())
	stemcells, err := deployment.Stemcells()
	if err != nil {
		log.Errorf("Error while reading Stemcells info for deployment `%s`: %v", deployment.Name(), err)
		return
	}

	for _, stemcell := range stemcells {
		ch <- prometheus.MustNewConstMetric(
			c.deploymentStemcellInfoDesc,
			prometheus.GaugeValue,
			float64(1),
			deployment.Name(),
			stemcell.Name(),
			stemcell.Version().AsString(),
			stemcell.OSName(),
		)
	}
}
