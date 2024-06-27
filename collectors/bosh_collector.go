package collectors

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/cloudfoundry/bosh_exporter/deployments"
	"github.com/cloudfoundry/bosh_exporter/filters"
)

type BoshCollector struct {
	enabledCollectors                   []Collector
	deploymentsFetcher                  *deployments.Fetcher
	totalBoshScrapesMetric              prometheus.Counter
	totalBoshScrapeErrorsMetric         prometheus.Counter
	lastBoshScrapeErrorMetric           prometheus.Gauge
	lastBoshScrapeTimestampMetric       prometheus.Gauge
	lastBoshScrapeDurationSecondsMetric prometheus.Gauge
}

func NewBoshCollector(
	namespace string,
	environment string,
	boshName string,
	boshUUID string,
	serviceDiscoveryFilename string,
	deploymentsFetcher *deployments.Fetcher,
	collectorsFilter *filters.CollectorsFilter,
	azsFilter *filters.AZsFilter,
	processesFilter *filters.RegexpFilter,
	cidrsFilter *filters.CidrFilter,
) *BoshCollector {
	var enabledCollectors []Collector

	if collectorsFilter.Enabled(filters.DeploymentsCollector) {
		deploymentsCollector := NewDeploymentsCollector(namespace, environment, boshName, boshUUID)
		enabledCollectors = append(enabledCollectors, deploymentsCollector)
	}

	if collectorsFilter.Enabled(filters.JobsCollector) {
		jobsCollector := NewJobsCollector(namespace, environment, boshName, boshUUID, azsFilter, cidrsFilter)
		enabledCollectors = append(enabledCollectors, jobsCollector)
	}

	if collectorsFilter.Enabled(filters.ServiceDiscoveryCollector) {
		serviceDiscoveryCollector := NewServiceDiscoveryCollector(
			namespace,
			environment,
			boshName,
			boshUUID,
			serviceDiscoveryFilename,
			azsFilter,
			processesFilter,
			cidrsFilter,
		)
		enabledCollectors = append(enabledCollectors, serviceDiscoveryCollector)
	}

	metrics := NewBoshCollectorMetrics(namespace, environment, boshName, boshUUID)
	return &BoshCollector{
		enabledCollectors:                   enabledCollectors,
		deploymentsFetcher:                  deploymentsFetcher,
		totalBoshScrapesMetric:              metrics.NewTotalBoshScrapesMetric(),
		totalBoshScrapeErrorsMetric:         metrics.NewTotalBoshScrapeErrorsMetric(),
		lastBoshScrapeErrorMetric:           metrics.NewLastBoshScrapeErrorMetric(),
		lastBoshScrapeTimestampMetric:       metrics.NewLastBoshScrapeTimestampMetric(),
		lastBoshScrapeDurationSecondsMetric: metrics.NewLastBoshScrapeDurationSecondsMetric(),
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

	c.totalBoshScrapesMetric.Describe(ch)
	c.totalBoshScrapeErrorsMetric.Describe(ch)
	c.lastBoshScrapeErrorMetric.Describe(ch)
	c.lastBoshScrapeTimestampMetric.Describe(ch)
	c.lastBoshScrapeDurationSecondsMetric.Describe(ch)
}

func (c *BoshCollector) Collect(ch chan<- prometheus.Metric) {
	var begun = time.Now()

	scrapeError := 0
	c.totalBoshScrapesMetric.Inc()
	ds, err := c.deploymentsFetcher.Deployments()
	if err != nil {
		log.Error(err)
		scrapeError = 1
		c.totalBoshScrapeErrorsMetric.Inc()
	} else {
		if err := c.executeCollectors(ds, ch); err != nil {
			log.Error(err)
			scrapeError = 1
			c.totalBoshScrapeErrorsMetric.Inc()
		}
	}

	c.totalBoshScrapesMetric.Collect(ch)

	c.totalBoshScrapeErrorsMetric.Collect(ch)

	c.lastBoshScrapeErrorMetric.Set(float64(scrapeError))
	c.lastBoshScrapeErrorMetric.Collect(ch)

	c.lastBoshScrapeTimestampMetric.Set(float64(time.Now().Unix()))
	c.lastBoshScrapeTimestampMetric.Collect(ch)

	c.lastBoshScrapeDurationSecondsMetric.Set(time.Since(begun).Seconds())
	c.lastBoshScrapeDurationSecondsMetric.Collect(ch)
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
