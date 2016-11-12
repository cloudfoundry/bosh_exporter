package collectors_test

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/director/fakes"
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
		boshClient         *fakes.FakeDirector
		deploymentsFilter  *filters.DeploymentsFilter
		deploymentsFetcher *deployments.Fetcher
		collectorsFilter   *filters.CollectorsFilter
		processesFilter    *filters.RegexpFilter
		boshCollector      *BoshCollector

		totalBoshScrapesDesc              *prometheus.Desc
		totalBoshScrapeErrorsDesc         *prometheus.Desc
		lastBoshScrapeErrorDesc           *prometheus.Desc
		lastBoshScrapeTimestampDesc       *prometheus.Desc
		lastBoshScrapeDurationSecondsDesc *prometheus.Desc
	)

	BeforeEach(func() {
		namespace = "test_exporter"
		tmpfile, err = ioutil.TempFile("", "service_discovery_collector_test_")
		Expect(err).ToNot(HaveOccurred())
		serviceDiscoveryFilename = tmpfile.Name()

		boshDeployments = []string{}
		boshClient = &fakes.FakeDirector{}
		deploymentsFilter = filters.NewDeploymentsFilter(boshDeployments, boshClient)
		deploymentsFetcher = deployments.NewFetcher(*deploymentsFilter)
		collectorsFilter, err = filters.NewCollectorsFilter([]string{})
		Expect(err).ToNot(HaveOccurred())
		processesFilter, err = filters.NewRegexpFilter([]string{})
		Expect(err).ToNot(HaveOccurred())

		totalBoshScrapesDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "scrapes_total"),
			"Total number of times BOSH was scraped for metrics.",
			[]string{},
			nil,
		)

		totalBoshScrapeErrorsDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "scrape_errors_total"),
			"Total number of times an error occured scraping BOSH.",
			[]string{},
			nil,
		)

		lastBoshScrapeErrorDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "last_scrape_error"),
			"Whether the last scrape of metrics from BOSH resulted in an error (1 for error, 0 for success).",
			[]string{},
			nil,
		)

		lastBoshScrapeTimestampDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "last_scrape_timestamp"),
			"Number of seconds since 1970 since last scrape from BOSH.",
			[]string{},
			nil,
		)

		lastBoshScrapeDurationSecondsDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "last_scrape_duration_seconds"),
			"Duration of the last scrape from BOSH.",
			[]string{},
			nil,
		)
	})

	AfterEach(func() {
		err = os.Remove(serviceDiscoveryFilename)
		Expect(err).ToNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		boshCollector = NewBoshCollector(namespace, serviceDiscoveryFilename, deploymentsFetcher, collectorsFilter, processesFilter)
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
			Eventually(descriptions).Should(Receive(Equal(totalBoshScrapesDesc)))
		})

		It("returns a scrape_errors_total description", func() {
			Eventually(descriptions).Should(Receive(Equal(totalBoshScrapeErrorsDesc)))
		})

		It("returns a last_scrape_error description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastBoshScrapeErrorDesc)))
		})

		It("returns a last_scrape_timestamp metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastBoshScrapeTimestampDesc)))
		})

		It("returns a last_scrape_duration_seconds metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastBoshScrapeDurationSecondsDesc)))
		})
	})

	Describe("Collect", func() {
		var (
			metrics                     chan prometheus.Metric
			totalBoshScrapesMetric      prometheus.Metric
			totalBoshScrapeErrorsMetric prometheus.Metric
			lastBoshScrapeErrorMetric   prometheus.Metric
		)

		BeforeEach(func() {
			metrics = make(chan prometheus.Metric)

			totalBoshScrapesMetric = prometheus.MustNewConstMetric(
				totalBoshScrapesDesc,
				prometheus.CounterValue,
				float64(1),
			)

			totalBoshScrapeErrorsMetric = prometheus.MustNewConstMetric(
				totalBoshScrapeErrorsDesc,
				prometheus.CounterValue,
				float64(0),
			)

			lastBoshScrapeErrorMetric = prometheus.MustNewConstMetric(
				lastBoshScrapeErrorDesc,
				prometheus.GaugeValue,
				float64(0),
			)
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

				totalBoshScrapeErrorsMetric = prometheus.MustNewConstMetric(
					totalBoshScrapeErrorsDesc,
					prometheus.CounterValue,
					float64(1),
				)

				lastBoshScrapeErrorMetric = prometheus.MustNewConstMetric(
					lastBoshScrapeErrorDesc,
					prometheus.GaugeValue,
					float64(1),
				)
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
