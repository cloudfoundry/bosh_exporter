package collectors_test

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/director/directorfakes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/bosh-prometheus/bosh_exporter/deployments"
	"github.com/bosh-prometheus/bosh_exporter/filters"

	. "github.com/bosh-prometheus/bosh_exporter/collectors"
	. "github.com/bosh-prometheus/bosh_exporter/utils/test_matchers"
)

func init() {
	log.Base().SetLevel("fatal")
}

var _ = Describe("BoshCollector", func() {
	var (
		err                      error
		namespace                string
		environment              string
		boshName                 string
		boshUUID                 string
		tmpfile                  *os.File
		serviceDiscoveryFilename string

		boshDeployments    []string
		boshClient         *directorfakes.FakeDirector
		deploymentsFilter  *filters.DeploymentsFilter
		deploymentsFetcher *deployments.Fetcher
		collectorsFilter   *filters.CollectorsFilter
		azsFilter          *filters.AZsFilter
		processesFilter    *filters.RegexpFilter
		cidrsFilter        *filters.CidrFilter
		boshCollector      *BoshCollector

		totalBoshScrapesMetric              prometheus.Counter
		totalBoshScrapeErrorsMetric         prometheus.Counter
		lastBoshScrapeErrorMetric           prometheus.Gauge
		lastBoshScrapeTimestampMetric       prometheus.Gauge
		lastBoshScrapeDurationSecondsMetric prometheus.Gauge
	)

	BeforeEach(func() {
		namespace = "test_exporter"
		environment = "test_environment"
		boshName = "test_bosh_name"
		boshUUID = "test_bosh_uuid"
		tmpfile, err = os.CreateTemp("", "service_discovery_collector_test_")
		Expect(err).ToNot(HaveOccurred())
		serviceDiscoveryFilename = tmpfile.Name()

		boshDeployments = []string{}
		boshClient = &directorfakes.FakeDirector{}
		deploymentsFilter = filters.NewDeploymentsFilter(boshDeployments, boshClient)
		deploymentsFetcher = deployments.NewFetcher(*deploymentsFilter)
		collectorsFilter, err = filters.NewCollectorsFilter([]string{})
		Expect(err).ToNot(HaveOccurred())
		azsFilter = filters.NewAZsFilter([]string{})
		cidrsFilter, err = filters.NewCidrFilter([]string{})
		Expect(err).ToNot(HaveOccurred())
		processesFilter, err = filters.NewRegexpFilter([]string{})
		Expect(err).ToNot(HaveOccurred())

		totalBoshScrapesMetric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "scrapes_total",
				Help:      "Total number of times BOSH was scraped for metrics.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
		)

		totalBoshScrapesMetric.Inc()

		totalBoshScrapeErrorsMetric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "scrape_errors_total",
				Help:      "Total number of times an error occured scraping BOSH.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
		)

		lastBoshScrapeErrorMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_scrape_error",
				Help:      "Whether the last scrape of metrics from BOSH resulted in an error (1 for error, 0 for success).",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
		)

		lastBoshScrapeErrorMetric.Set(float64(0))

		lastBoshScrapeTimestampMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_scrape_timestamp",
				Help:      "Number of seconds since 1970 since last scrape from BOSH.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
		)

		lastBoshScrapeDurationSecondsMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_scrape_duration_seconds",
				Help:      "Duration of the last scrape from BOSH.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
		)
	})

	AfterEach(func() {
		err = os.Remove(serviceDiscoveryFilename)
		Expect(err).ToNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		boshCollector = NewBoshCollector(
			namespace,
			environment,
			boshName,
			boshUUID,
			serviceDiscoveryFilename,
			deploymentsFetcher,
			collectorsFilter,
			azsFilter,
			processesFilter,
			cidrsFilter,
		)
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
			Eventually(metrics).Should(Receive(PrometheusMetric(totalBoshScrapesMetric)))
		})

		It("returns a scrape_errors_total metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(totalBoshScrapeErrorsMetric)))
		})

		It("returns a last_scrape_error metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(lastBoshScrapeErrorMetric)))
		})

		Context("when it fails to get the deployment", func() {
			BeforeEach(func() {
				boshClient.DeploymentsReturns([]director.Deployment{}, errors.New("no deployments"))

				totalBoshScrapeErrorsMetric.Inc()
				lastBoshScrapeErrorMetric.Set(float64(1))
			})

			It("returns a scrape_errors_total metric", func() {
				Eventually(metrics).Should(Receive(PrometheusMetric(totalBoshScrapeErrorsMetric)))
			})

			It("returns a last_scrape_error metric", func() {
				Eventually(metrics).Should(Receive(PrometheusMetric(lastBoshScrapeErrorMetric)))
			})
		})
	})
})
