package collectors_test

import (
	"io/ioutil"
	"os"

	. "github.com/benjamintf1/unmarshalledmatchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/bosh-prometheus/bosh_exporter/deployments"
	"github.com/bosh-prometheus/bosh_exporter/filters"

	. "github.com/bosh-prometheus/bosh_exporter/collectors"
)

func init() {
	log.Base().SetLevel("fatal")
}

var _ = Describe("ServiceDiscoveryCollector", func() {
	var (
		err                       error
		namespace                 string
		environment               string
		boshName                  string
		boshUUID                  string
		tmpfile                   *os.File
		serviceDiscoveryFilename  string
		azsFilter                 *filters.AZsFilter
		processesFilter           *filters.RegexpFilter
		cidrsFilter               *filters.CidrFilter
		serviceDiscoveryCollector *ServiceDiscoveryCollector

		lastServiceDiscoveryScrapeTimestampMetric       prometheus.Gauge
		lastServiceDiscoveryScrapeDurationSecondsMetric prometheus.Gauge
	)

	BeforeEach(func() {
		namespace = "test_exporter"
		environment = "test_environment"
		boshName = "test_bosh_name"
		boshUUID = "test_bosh_uuid"
		tmpfile, err = ioutil.TempFile("", "service_discovery_collector_test_")
		Expect(err).ToNot(HaveOccurred())
		serviceDiscoveryFilename = tmpfile.Name()
		azsFilter = filters.NewAZsFilter([]string{})
		cidrsFilter, err = filters.NewCidrFilter([]string{"0.0.0.0/0"})
		processesFilter, err = filters.NewRegexpFilter([]string{})

		lastServiceDiscoveryScrapeTimestampMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_service_discovery_scrape_timestamp",
				Help:      "Number of seconds since 1970 since last scrape of Service Discovery from BOSH.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
		)

		lastServiceDiscoveryScrapeDurationSecondsMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_service_discovery_scrape_duration_seconds",
				Help:      "Duration of the last scrape of Service Discovery from BOSH.",
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
		serviceDiscoveryCollector = NewServiceDiscoveryCollector(
			namespace,
			environment,
			boshName,
			boshUUID,
			serviceDiscoveryFilename,
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
			go serviceDiscoveryCollector.Describe(descriptions)
		})

		It("returns a last_service_discovery_scrape_duration_seconds metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastServiceDiscoveryScrapeTimestampMetric.Desc())))
		})

		It("returns a last_service_discovery_scrape_duration_seconds metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastServiceDiscoveryScrapeDurationSecondsMetric.Desc())))
		})
	})

	Describe("Collect", func() {
		var (
			deployment1Name     = "fake-deployment-1-name"
			deployment2Name     = "fake-deployment-2-name"
			job1Name            = "fake-job-1-name"
			job2Name            = "fake-job-2-name"
			job1AZ              = "fake-job-1-az"
			job2AZ              = "fake-job-2-az"
			job1IP              = "1.2.3.4"
			job2IP              = "5.6.7.8"
			jobProcess1Name     = "fake-process-1-name"
			jobProcess2Name     = "fake-process-2-name"
			targetGroupsContent = `[
				{"targets":["1.2.3.4"],"labels":{"__meta_bosh_deployment":"fake-deployment-1-name","__meta_bosh_job_process_name":"fake-process-1-name"}},
				{"targets":["1.2.3.4"],"labels":{"__meta_bosh_deployment":"fake-deployment-1-name","__meta_bosh_job_process_name":"fake-process-2-name"}},
				{"targets":["5.6.7.8"],"labels":{"__meta_bosh_deployment":"fake-deployment-2-name","__meta_bosh_job_process_name":"fake-process-2-name"}}
			]`

			deployment1Processes []deployments.Process
			deployment2Processes []deployments.Process
			deployment1Instances []deployments.Instance
			deployment2Instances []deployments.Instance
			deployment1Info      deployments.DeploymentInfo
			deployment2Info      deployments.DeploymentInfo
			deploymentsInfo      []deployments.DeploymentInfo

			metrics    chan prometheus.Metric
			errMetrics chan error
		)

		BeforeEach(func() {
			deployment1Processes = []deployments.Process{
				{
					Name: jobProcess1Name,
				},
				{
					Name: jobProcess2Name,
				},
			}

			deployment2Processes = []deployments.Process{
				{
					Name: jobProcess2Name,
				},
			}
			deployment1Instances = []deployments.Instance{
				{
					Name:      job1Name,
					IPs:       []string{job1IP},
					AZ:        job1AZ,
					Processes: deployment1Processes,
				},
			}

			deployment2Instances = []deployments.Instance{
				{
					Name:      job2Name,
					IPs:       []string{job2IP},
					AZ:        job2AZ,
					Processes: deployment2Processes,
				},
			}

			deployment1Info = deployments.DeploymentInfo{
				Name:      deployment1Name,
				Instances: deployment1Instances,
			}

			deployment2Info = deployments.DeploymentInfo{
				Name:      deployment2Name,
				Instances: deployment2Instances,
			}

			deploymentsInfo = []deployments.DeploymentInfo{deployment1Info, deployment2Info}

			metrics = make(chan prometheus.Metric)
			errMetrics = make(chan error, 1)
		})

		JustBeforeEach(func() {
			go func() {
				if err := serviceDiscoveryCollector.Collect(deploymentsInfo, metrics); err != nil {
					errMetrics <- err
				}
			}()
		})

		It("writes a target groups file", func() {
			Eventually(metrics).Should(Receive())
			targetGroups, err := ioutil.ReadFile(serviceDiscoveryFilename)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(targetGroups)).To(MatchUnorderedJSON(targetGroupsContent))
		})

		It("returns a last_service_discovery_scrape_timestamp & last_service_discovery_scrape_duration_seconds", func() {
			Eventually(metrics).Should(Receive())
			Eventually(metrics).Should(Receive())
			Consistently(metrics).ShouldNot(Receive())
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there are no deployments", func() {
			BeforeEach(func() {
				deploymentsInfo = []deployments.DeploymentInfo{}
			})

			It("writes an empty target groups file", func() {
				Eventually(metrics).Should(Receive())
				targetGroups, err := ioutil.ReadFile(serviceDiscoveryFilename)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(targetGroups)).To(Equal("[]"))
			})

			It("returns only last_service_discovery_scrape_timestamp & last_service_discovery_scrape_duration_seconds", func() {
				Eventually(metrics).Should(Receive())
				Eventually(metrics).Should(Receive())
				Consistently(metrics).ShouldNot(Receive())
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		Context("when there are no instances", func() {
			BeforeEach(func() {
				deployment1Info.Instances = []deployments.Instance{}
				deploymentsInfo = []deployments.DeploymentInfo{deployment1Info}
			})

			It("writes an empty target groups file", func() {
				Eventually(metrics).Should(Receive())
				targetGroups, err := ioutil.ReadFile(serviceDiscoveryFilename)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(targetGroups)).To(Equal("[]"))
			})

			It("returns only last_service_discovery_scrape_timestamp & last_service_discovery_scrape_duration_seconds", func() {
				Eventually(metrics).Should(Receive())
				Eventually(metrics).Should(Receive())
				Consistently(metrics).ShouldNot(Receive())
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		Context("when instance has no IP", func() {
			BeforeEach(func() {
				deployment1Info.Instances[0].IPs = []string{}
				deploymentsInfo = []deployments.DeploymentInfo{deployment1Info}
			})

			It("writes an empty target groups file", func() {
				Eventually(metrics).Should(Receive())
				targetGroups, err := ioutil.ReadFile(serviceDiscoveryFilename)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(targetGroups)).To(Equal("[]"))
			})

			It("returns only last_service_discovery_scrape_timestamp & last_service_discovery_scrape_duration_seconds", func() {
				Eventually(metrics).Should(Receive())
				Eventually(metrics).Should(Receive())
				Consistently(metrics).ShouldNot(Receive())
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		Context("when no IP is found for an instance", func() {
			BeforeEach(func() {
				cidrsFilter, err = filters.NewCidrFilter([]string{"10.254.0.0/16"})
			})

			It("writes an empty target groups file", func() {
				Eventually(metrics).Should(Receive())
				targetGroups, err := ioutil.ReadFile(serviceDiscoveryFilename)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(targetGroups)).To(Equal("[]"))
			})

			It("returns only last_service_discovery_scrape_timestamp & last_service_discovery_scrape_duration_seconds", func() {
				Eventually(metrics).Should(Receive())
				Eventually(metrics).Should(Receive())
				Consistently(metrics).ShouldNot(Receive())
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		Context("when there are no processes", func() {
			BeforeEach(func() {
				deployment1Info.Instances[0].Processes = []deployments.Process{}
				deploymentsInfo = []deployments.DeploymentInfo{deployment1Info}
			})

			It("writes an empty target groups file", func() {
				Eventually(metrics).Should(Receive())
				targetGroups, err := ioutil.ReadFile(serviceDiscoveryFilename)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(targetGroups)).To(Equal("[]"))
			})

			It("returns only last_service_discovery_scrape_timestamp & last_service_discovery_scrape_duration_seconds", func() {
				Eventually(metrics).Should(Receive())
				Eventually(metrics).Should(Receive())
				Consistently(metrics).ShouldNot(Receive())
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})
	})
})
