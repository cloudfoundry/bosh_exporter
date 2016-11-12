package collectors_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director/fakes"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/deployments"
	"github.com/cloudfoundry-community/bosh_exporter/filters"

	. "github.com/cloudfoundry-community/bosh_exporter/collectors"
)

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

		totalScrapesDesc                  *prometheus.Desc
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

		totalScrapesDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "scrapes_total"),
			"Total number of times BOSH was scraped for metrics.",
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
		os.Remove(serviceDiscoveryFilename)
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
			Eventually(descriptions).Should(Receive(Equal(totalScrapesDesc)))
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
			metrics            chan prometheus.Metric
			totalScrapesMetric prometheus.Metric
		)

		BeforeEach(func() {
			metrics = make(chan prometheus.Metric)

			totalScrapesMetric = prometheus.MustNewConstMetric(
				totalScrapesDesc,
				prometheus.CounterValue,
				float64(1),
			)
		})

		JustBeforeEach(func() {
			go boshCollector.Collect(metrics)
		})

		It("returns a scrapes_total metric", func() {
			Eventually(metrics).Should(Receive(Equal(totalScrapesMetric)))
		})
	})
})
