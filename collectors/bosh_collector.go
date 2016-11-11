package collectors

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/deployments"
	"github.com/cloudfoundry-community/bosh_exporter/filters"
)

type BoshCollector struct {
	enabledCollectors  []Collector
	deploymentsFetcher *deployments.Fetcher
}

func NewBoshCollector(
	metricsNamespace string,
	serviceDiscoveryFilename string,
	deploymentsFetcher *deployments.Fetcher,
	collectorsFilter *filters.CollectorsFilter,
	processesFilter *filters.RegexpFilter,
) *BoshCollector {
	enabledCollectors := []Collector{}

	if collectorsFilter.Enabled(filters.DeploymentsCollector) {
		deploymentsCollector := NewDeploymentsCollector(metricsNamespace)
		enabledCollectors = append(enabledCollectors, deploymentsCollector)
	}

	if collectorsFilter.Enabled(filters.JobsCollector) {
		jobsCollector := NewJobsCollector(metricsNamespace)
		enabledCollectors = append(enabledCollectors, jobsCollector)
	}

	if collectorsFilter.Enabled(filters.ServiceDiscoveryCollector) {
		serviceDiscoveryCollector := NewServiceDiscoveryCollector(
			metricsNamespace,
			serviceDiscoveryFilename,
			*processesFilter,
		)
		enabledCollectors = append(enabledCollectors, serviceDiscoveryCollector)
	}

	return &BoshCollector{
		enabledCollectors:  enabledCollectors,
		deploymentsFetcher: deploymentsFetcher,
	}
}

func (c *BoshCollector) Describe(ch chan<- *prometheus.Desc) {
	var wg = &sync.WaitGroup{}

	for _, collector := range c.enabledCollectors {
		wg.Add(1)
		go func(collector Collector, ch chan<- *prometheus.Desc) {
			defer wg.Done()
			collector.Describe(ch)
		}(collector, ch)
	}
	wg.Wait()
}

func (c *BoshCollector) Collect(ch chan<- prometheus.Metric) {
	var wg = &sync.WaitGroup{}

	deployments := c.deploymentsFetcher.Deployments()
	for _, collector := range c.enabledCollectors {
		wg.Add(1)
		go func(collector Collector, ch chan<- prometheus.Metric) {
			defer wg.Done()
			collector.Collect(deployments, ch)
		}(collector, ch)
	}
	wg.Wait()
}
