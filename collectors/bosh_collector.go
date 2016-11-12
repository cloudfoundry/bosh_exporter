package collectors

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/deployments"
	"github.com/cloudfoundry-community/bosh_exporter/filters"
)

type BoshCollector struct {
	enabledCollectors                 []Collector
	deploymentsFetcher                *deployments.Fetcher
	totalScrapes                      uint64
	totalScrapesDesc                  *prometheus.Desc
	lastBoshScrapeTimestampDesc       *prometheus.Desc
	lastBoshScrapeDurationSecondsDesc *prometheus.Desc
}

func NewBoshCollector(
	namespace string,
	serviceDiscoveryFilename string,
	deploymentsFetcher *deployments.Fetcher,
	collectorsFilter *filters.CollectorsFilter,
	processesFilter *filters.RegexpFilter,
) *BoshCollector {
	enabledCollectors := []Collector{}

	if collectorsFilter.Enabled(filters.DeploymentsCollector) {
		deploymentsCollector := NewDeploymentsCollector(namespace)
		enabledCollectors = append(enabledCollectors, deploymentsCollector)
	}

	if collectorsFilter.Enabled(filters.JobsCollector) {
		jobsCollector := NewJobsCollector(namespace)
		enabledCollectors = append(enabledCollectors, jobsCollector)
	}

	if collectorsFilter.Enabled(filters.ServiceDiscoveryCollector) {
		serviceDiscoveryCollector := NewServiceDiscoveryCollector(
			namespace,
			serviceDiscoveryFilename,
			*processesFilter,
		)
		enabledCollectors = append(enabledCollectors, serviceDiscoveryCollector)
	}

	totalScrapesDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scrapes_total"),
		"Total number of times BOSH was scraped for metrics.",
		[]string{},
		nil,
	)

	lastBoshScrapeTimestampDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "last_scrape_timestamp"),
		"Number of seconds since 1970 since last scrape from BOSH.",
		[]string{},
		nil,
	)

	lastBoshScrapeDurationSecondsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "last_scrape_duration_seconds"),
		"Duration of the last scrape from BOSH.",
		[]string{},
		nil,
	)

	return &BoshCollector{
		enabledCollectors:                 enabledCollectors,
		deploymentsFetcher:                deploymentsFetcher,
		totalScrapes:                      0,
		totalScrapesDesc:                  totalScrapesDesc,
		lastBoshScrapeTimestampDesc:       lastBoshScrapeTimestampDesc,
		lastBoshScrapeDurationSecondsDesc: lastBoshScrapeDurationSecondsDesc,
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

	ch <- c.totalScrapesDesc
	ch <- c.lastBoshScrapeTimestampDesc
	ch <- c.lastBoshScrapeDurationSecondsDesc
}

func (c *BoshCollector) Collect(ch chan<- prometheus.Metric) {
	var begun = time.Now()
	var wg = &sync.WaitGroup{}

	c.totalScrapes++
	deployments := c.deploymentsFetcher.Deployments()
	for _, collector := range c.enabledCollectors {
		wg.Add(1)
		go func(collector Collector, ch chan<- prometheus.Metric) {
			defer wg.Done()
			collector.Collect(deployments, ch)
		}(collector, ch)
	}
	wg.Wait()

	ch <- prometheus.MustNewConstMetric(
		c.totalScrapesDesc,
		prometheus.CounterValue,
		float64(c.totalScrapes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.lastBoshScrapeTimestampDesc,
		prometheus.GaugeValue,
		float64(time.Now().Unix()),
	)

	ch <- prometheus.MustNewConstMetric(
		c.lastBoshScrapeDurationSecondsDesc,
		prometheus.GaugeValue,
		time.Since(begun).Seconds(),
	)
}
