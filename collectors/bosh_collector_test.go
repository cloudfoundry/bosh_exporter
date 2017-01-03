package collectors_test

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/director/directorfakes"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/deployments"
	"github.com/cloudfoundry-community/bosh_exporter/filters"

	. "github.com/cloudfoundry-community/bosh_exporter/collectors"
)

func init() {
	flag.Set("log.level", "fatal")
}

var _ = Describe("BoshCollector", func() {
	var (
		err                      error
		namespace                string
		tmpfile                  *os.File
		serviceDiscoveryFilename string

		boshDeployments    []string
		boshClient         *directorfakes.FakeDirector
		deploymentsFilter  *filters.DeploymentsFilter
		deploymentsFetcher *deployments.Fetcher
		collectorsFilter   *filters.CollectorsFilter
		azsFilter          *filters.AZsFilter
		processesFilter    *filters.RegexpFilter
		boshCollector      *BoshCollector

		totalBoshScrapesMetric              prometheus.Counter
		totalBoshScrapeErrorsMetric         prometheus.Counter
		lastBoshScrapeErrorMetric           prometheus.Gauge
		lastBoshScrapeTimestampMetric       prometheus.Gauge
		lastBoshScrapeDurationSecondsMetric prometheus.Gauge
	)

	BeforeEach(func() {
		namespace = "test_exporter"
		tmpfile, err = ioutil.TempFile("", "service_discovery_collector_test_")
		Expect(err).ToNot(HaveOccurred())
		serviceDiscoveryFilename = tmpfile.Name()

		boshDeployments = []string{}
		boshClient = &directorfakes.FakeDirector{}
		deploymentsFilter = filters.NewDeploymentsFilter(boshDeployments, boshClient)
		deploymentsFetcher = deployments.NewFetcher(*deploymentsFilter)
		collectorsFilter, err = filters.NewCollectorsFilter([]string{})
		Expect(err).ToNot(HaveOccurred())
		azsFilter = filters.NewAZsFilter([]string{})
		processesFilter, err = filters.NewRegexpFilter([]string{})
		Expect(err).ToNot(HaveOccurred())

		totalBoshScrapesMetric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "scrapes_total",
				Help:      "Total number of times BOSH was scraped for metrics.",
			},
		)

		totalBoshScrapesMetric.Inc()

		totalBoshScrapeErrorsMetric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "scrape_errors_total",
				Help:      "Total number of times an error occured scraping BOSH.",
			},
		)

		lastBoshScrapeErrorMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_scrape_error",
				Help:      "Whether the last scrape of metrics from BOSH resulted in an error (1 for error, 0 for success).",
			},
		)

		lastBoshScrapeErrorMetric.Set(float64(0))

		lastBoshScrapeTimestampMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_scrape_timestamp",
				Help:      "Number of seconds since 1970 since last scrape from BOSH.",
			},
		)

		lastBoshScrapeDurationSecondsMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_scrape_duration_seconds",
				Help:      "Duration of the last scrape from BOSH.",
			},
		)
	})

	AfterEach(func() {
		err = os.Remove(serviceDiscoveryFilename)
		Expect(err).ToNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		boshCollector = NewBoshCollector(namespace, serviceDiscoveryFilename, deploymentsFetcher, collectorsFilter, azsFilter, processesFilter)
	})

	Describe("Describe", func() {
		var (
			descriptions chan *prometheus.Desc
		)

		BeforeEach(func() {
			descriptions = make(chan *prometheus.Desc)
		})

		JustBeforeEach(func() {
			go boshCollector.Describe(descriptions)
		})

		It("returns a scrapes_total description", func() {
			Eventually(descriptions).Should(Receive(Equal(totalBoshScrapesMetric.Desc())))
		})

		It("returns a scrape_errors_total description", func() {
			Eventually(descriptions).Should(Receive(Equal(totalBoshScrapeErrorsMetric.Desc())))
		})

		It("returns a last_scrape_error description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastBoshScrapeErrorMetric.Desc())))
		})

		It("returns a last_scrape_timestamp metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastBoshScrapeTimestampMetric.Desc())))
		})

		It("returns a last_scrape_duration_seconds metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastBoshScrapeDurationSecondsMetric.Desc())))
		})
	})

	Describe("Collect", func() {
		var (
			metrics chan prometheus.Metric
		)

		BeforeEach(func() {
			metrics = make(chan prometheus.Metric)
		})

		JustBeforeEach(func() {
			go boshCollector.Collect(metrics)
		})

		It("returns a scrapes_total metric", func() {
			Eventually(metrics).Should(Receive(Equal(totalBoshScrapesMetric)))
		})

		It("returns a scrape_errors_total metric", func() {
			Eventually(metrics).Should(Receive(Equal(totalBoshScrapeErrorsMetric)))
		})

		It("returns a last_scrape_error metric", func() {
			Eventually(metrics).Should(Receive(Equal(lastBoshScrapeErrorMetric)))
		})

		Context("when it fails to get the deployment", func() {
			BeforeEach(func() {
				boshClient.DeploymentsReturns([]director.Deployment{}, errors.New("no deployments"))

				totalBoshScrapeErrorsMetric.Inc()
				lastBoshScrapeErrorMetric.Set(float64(1))
			})

			It("returns a scrape_errors_total metric", func() {
				Eventually(metrics).Should(Receive(Equal(totalBoshScrapeErrorsMetric)))
			})

			It("returns a last_scrape_error metric", func() {
				Eventually(metrics).Should(Receive(Equal(lastBoshScrapeErrorMetric)))
			})
		})
	})
})
