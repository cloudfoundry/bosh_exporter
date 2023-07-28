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
	. "github.com/bosh-prometheus/bosh_exporter/utils/matchers"
)

func init() {
	_ = log.Base().SetLevel("fatal")
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
		metrics            *BoshCollectorMetrics
		boshCollector      *BoshCollector

		totalBoshScrapesMetric              prometheus.Counter
		totalBoshScrapeErrorsMetric         prometheus.Counter
		lastBoshScrapeErrorMetric           prometheus.Gauge
		lastBoshScrapeTimestampMetric       prometheus.Gauge
		lastBoshScrapeDurationSecondsMetric prometheus.Gauge
	)

	BeforeEach(func() {
		namespace = testNamespace
		environment = testEnvironment
		boshName = testBoshName
		boshUUID = testBoshUUID
		metrics = NewBoshCollectorMetrics(testNamespace, testEnvironment, testBoshName, testBoshUUID)
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

		totalBoshScrapesMetric = metrics.NewTotalBoshScrapesMetric()
		totalBoshScrapesMetric.Inc()
		totalBoshScrapeErrorsMetric = metrics.NewTotalBoshScrapeErrorsMetric()
		lastBoshScrapeErrorMetric = metrics.NewLastBoshScrapeErrorMetric()
		lastBoshScrapeErrorMetric.Set(float64(0))
		lastBoshScrapeTimestampMetric = metrics.NewLastBoshScrapeTimestampMetric()
		lastBoshScrapeDurationSecondsMetric = metrics.NewLastBoshScrapeDurationSecondsMetric()
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
