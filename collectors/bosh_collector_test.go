package collectors_test

import (
	"errors"
	"os"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/director/directorfakes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/bosh-prometheus/bosh_exporter/deployments"
	"github.com/bosh-prometheus/bosh_exporter/filters"

	"github.com/bosh-prometheus/bosh_exporter/collectors"
	"github.com/bosh-prometheus/bosh_exporter/utils/matchers"
)

func init() {
	_ = log.Base().SetLevel("fatal")
}

var _ = ginkgo.Describe("BoshCollector", func() {
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
		metrics            *collectors.BoshCollectorMetrics
		boshCollector      *collectors.BoshCollector

		totalBoshScrapesMetric              prometheus.Counter
		totalBoshScrapeErrorsMetric         prometheus.Counter
		lastBoshScrapeErrorMetric           prometheus.Gauge
		lastBoshScrapeTimestampMetric       prometheus.Gauge
		lastBoshScrapeDurationSecondsMetric prometheus.Gauge
	)

	ginkgo.BeforeEach(func() {
		namespace = testNamespace
		environment = testEnvironment
		boshName = testBoshName
		boshUUID = testBoshUUID
		metrics = collectors.NewBoshCollectorMetrics(testNamespace, testEnvironment, testBoshName, testBoshUUID)
		tmpfile, err = os.CreateTemp("", "service_discovery_collector_test_")
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		serviceDiscoveryFilename = tmpfile.Name()

		boshDeployments = []string{}
		boshClient = &directorfakes.FakeDirector{}
		deploymentsFilter = filters.NewDeploymentsFilter(boshDeployments, boshClient)
		deploymentsFetcher = deployments.NewFetcher(*deploymentsFilter)
		collectorsFilter, err = filters.NewCollectorsFilter([]string{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		azsFilter = filters.NewAZsFilter([]string{})
		cidrsFilter, err = filters.NewCidrFilter([]string{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		processesFilter, err = filters.NewRegexpFilter([]string{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		totalBoshScrapesMetric = metrics.NewTotalBoshScrapesMetric()
		totalBoshScrapesMetric.Inc()
		totalBoshScrapeErrorsMetric = metrics.NewTotalBoshScrapeErrorsMetric()
		lastBoshScrapeErrorMetric = metrics.NewLastBoshScrapeErrorMetric()
		lastBoshScrapeErrorMetric.Set(float64(0))
		lastBoshScrapeTimestampMetric = metrics.NewLastBoshScrapeTimestampMetric()
		lastBoshScrapeDurationSecondsMetric = metrics.NewLastBoshScrapeDurationSecondsMetric()
	})

	ginkgo.AfterEach(func() {
		err = os.Remove(serviceDiscoveryFilename)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.JustBeforeEach(func() {
		boshCollector = collectors.NewBoshCollector(
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

	ginkgo.Describe("Describe", func() {
		var (
			descriptions chan *prometheus.Desc
		)

		ginkgo.BeforeEach(func() {
			descriptions = make(chan *prometheus.Desc)
		})

		ginkgo.JustBeforeEach(func() {
			go boshCollector.Describe(descriptions)
		})

		ginkgo.It("returns a scrapes_total description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(totalBoshScrapesMetric.Desc())))
		})

		ginkgo.It("returns a scrape_errors_total description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(totalBoshScrapeErrorsMetric.Desc())))
		})

		ginkgo.It("returns a last_scrape_error description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(lastBoshScrapeErrorMetric.Desc())))
		})

		ginkgo.It("returns a last_scrape_timestamp metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(lastBoshScrapeTimestampMetric.Desc())))
		})

		ginkgo.It("returns a last_scrape_duration_seconds metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(lastBoshScrapeDurationSecondsMetric.Desc())))
		})
	})

	ginkgo.Describe("Collect", func() {
		var (
			metrics chan prometheus.Metric
		)

		ginkgo.BeforeEach(func() {
			metrics = make(chan prometheus.Metric)
		})

		ginkgo.JustBeforeEach(func() {
			go boshCollector.Collect(metrics)
		})

		ginkgo.It("returns a scrapes_total metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(totalBoshScrapesMetric)))
		})

		ginkgo.It("returns a scrape_errors_total metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(totalBoshScrapeErrorsMetric)))
		})

		ginkgo.It("returns a last_scrape_error metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(lastBoshScrapeErrorMetric)))
		})

		ginkgo.Context("when it fails to get the deployment", func() {
			ginkgo.BeforeEach(func() {
				boshClient.DeploymentsReturns([]director.Deployment{}, errors.New("no deployments"))

				totalBoshScrapeErrorsMetric.Inc()
				lastBoshScrapeErrorMetric.Set(float64(1))
			})

			ginkgo.It("returns a scrape_errors_total metric", func() {
				gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(totalBoshScrapeErrorsMetric)))
			})

			ginkgo.It("returns a last_scrape_error metric", func() {
				gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(lastBoshScrapeErrorMetric)))
			})
		})
	})
})
