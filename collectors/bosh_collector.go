package collectors

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/cloudfoundry-community/bosh_exporter/deployments"
	"github.com/cloudfoundry-community/bosh_exporter/filters"
)

type BoshCollector struct {
	enabledCollectors                 []Collector
	deploymentsFetcher                *deployments.Fetcher
	totalBoshScrapes                  uint64
	totalBoshScrapesDesc              *prometheus.Desc
	totalBoshScrapeErrors             uint64
	totalBoshScrapeErrorsDesc         *prometheus.Desc
	lastBoshScrapeErrorDesc           *prometheus.Desc
	lastBoshScrapeTimestampDesc       *prometheus.Desc
	lastBoshScrapeDurationSecondsDesc *prometheus.Desc
}

func NewBoshCollector(
	namespace string,
	serviceDiscoveryFilename string,
	deploymentsFetcher *deployments.Fetcher,
	collectorsFilter *filters.CollectorsFilter,
	azsFilter *filters.AZsFilter,
	processesFilter *filters.RegexpFilter,
) *BoshCollector {
	enabledCollectors := []Collector{}

	if collectorsFilter.Enabled(filters.DeploymentsCollector) {
		deploymentsCollector := NewDeploymentsCollector(namespace)
		enabledCollectors = append(enabledCollectors, deploymentsCollector)
	}

	if collectorsFilter.Enabled(filters.JobsCollector) {
		jobsCollector := NewJobsCollector(namespace, azsFilter)
		enabledCollectors = append(enabledCollectors, jobsCollector)
	}

	if collectorsFilter.Enabled(filters.ServiceDiscoveryCollector) {
		serviceDiscoveryCollector := NewServiceDiscoveryCollector(
			namespace,
			serviceDiscoveryFilename,
			azsFilter,
			processesFilter,
		)
		enabledCollectors = append(enabledCollectors, serviceDiscoveryCollector)
	}

	totalBoshScrapesDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scrapes_total"),
		"Total number of times BOSH was scraped for metrics.",
		[]string{},
		nil,
	)

	totalBoshScrapeErrorsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scrape_errors_total"),
		"Total number of times an error occured scraping BOSH.",
		[]string{},
		nil,
	)

	lastBoshScrapeErrorDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "last_scrape_error"),
		"Whether the last scrape of metrics from BOSH resulted in an error (1 for error, 0 for success).",
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
		totalBoshScrapes:                  0,
		totalBoshScrapesDesc:              totalBoshScrapesDesc,
		totalBoshScrapeErrors:             0,
		totalBoshScrapeErrorsDesc:         totalBoshScrapeErrorsDesc,
		lastBoshScrapeErrorDesc:           lastBoshScrapeErrorDesc,
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

	ch <- c.totalBoshScrapesDesc
	ch <- c.totalBoshScrapeErrorsDesc
	ch <- c.lastBoshScrapeErrorDesc
	ch <- c.lastBoshScrapeTimestampDesc
	ch <- c.lastBoshScrapeDurationSecondsDesc
}

func (c *BoshCollector) Collect(ch chan<- prometheus.Metric) {
	var begun = time.Now()

	scrapeError := 0
	c.totalBoshScrapes++
	deployments, err := c.deploymentsFetcher.Deployments()
	if err != nil {
		log.Error(err)
		scrapeError = 1
		c.totalBoshScrapeErrors++
	} else {
		if err := c.executeCollectors(deployments, ch); err != nil {
			log.Error(err)
			scrapeError = 1
			c.totalBoshScrapeErrors++
		}
	}

	ch <- prometheus.MustNewConstMetric(
		c.totalBoshScrapesDesc,
		prometheus.CounterValue,
		float64(c.totalBoshScrapes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.totalBoshScrapeErrorsDesc,
		prometheus.CounterValue,
		float64(c.totalBoshScrapeErrors),
	)

	ch <- prometheus.MustNewConstMetric(
		c.lastBoshScrapeErrorDesc,
		prometheus.GaugeValue,
		float64(scrapeError),
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

func (c *BoshCollector) executeCollectors(deployments []deployments.DeploymentInfo, ch chan<- prometheus.Metric) error {
	var wg = &sync.WaitGroup{}

	doneChannel := make(chan bool, 1)
	errChannel := make(chan error, 1)

	for _, collector := range c.enabledCollectors {
		wg.Add(1)
		go func(collector Collector) {
			defer wg.Done()
			if err := collector.Collect(deployments, ch); err != nil {
				errChannel <- err
			}
		}(collector)
	}

	go func() {
		wg.Wait()
		close(doneChannel)
	}()

	select {
	case <-doneChannel:
	case err := <-errChannel:
		return err
	}

	return nil
}
