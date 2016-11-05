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

	"github.com/cloudfoundry-community/bosh_exporter/filters"

	. "github.com/cloudfoundry-community/bosh_exporter/collectors"
)

func init() {
	flag.Set("log.level", "fatal")
}

var _ = Describe("ServiceDiscoveryCollector", func() {
	var (
		err                       error
		namespace                 string
		boshDeployments           []string
		deploymentsFilter         *filters.DeploymentsFilter
		tmpfile                   *os.File
		serviceDiscoveryFilename  string
		processesFilter           *filters.RegexpFilter
		boshClient                *fakes.FakeDirector
		serviceDiscoveryCollector *ServiceDiscoveryCollector

		lastServiceDiscoveryScrapeTimestampDesc       *prometheus.Desc
		lastServiceDiscoveryScrapeDurationSecondsDesc *prometheus.Desc
	)

	BeforeEach(func() {
		namespace = "test_exporter"
		boshDeployments = []string{}
		tmpfile, err = ioutil.TempFile("", "service_discovery_collector_test_")
		Expect(err).ToNot(HaveOccurred())
		serviceDiscoveryFilename = tmpfile.Name()
		boshClient = &fakes.FakeDirector{}

		lastServiceDiscoveryScrapeTimestampDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "last_service_discovery_scrape_timestamp"),
			"Number of seconds since 1970 since last scrape of Service Discovery from BOSH.",
			[]string{},
			nil,
		)

		lastServiceDiscoveryScrapeDurationSecondsDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "last_service_discovery_scrape_duration_seconds"),
			"Duration of the last scrape of Service Discovery from BOSH.",
			[]string{},
			nil,
		)
	})

	AfterEach(func() {
		os.Remove(serviceDiscoveryFilename)
	})

	JustBeforeEach(func() {
		deploymentsFilter = filters.NewDeploymentsFilter(boshDeployments, boshClient)
		processesFilter, err = filters.NewRegexpFilter([]string{})
		serviceDiscoveryCollector = NewServiceDiscoveryCollector(namespace, *deploymentsFilter, serviceDiscoveryFilename, *processesFilter)
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
			Eventually(descriptions).Should(Receive(Equal(lastServiceDiscoveryScrapeTimestampDesc)))
		})

		It("returns a last_service_discovery_scrape_duration_seconds metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastServiceDiscoveryScrapeDurationSecondsDesc)))
		})
	})

	Describe("Collect", func() {
		var (
			deploymentName      = "fake-deployment-name"
			jobName             = "fake-job-name"
			jobID               = "fake-job-id"
			jobIndex            = 0
			jobAZ               = "fake-job-az"
			jobIP               = "1.2.3.4"
			processState        = "running"
			jobProcessName      = "fake-process-name"
			jobProcessState     = "running"
			targetGroupsContent = "[{\"targets\":[\"1.2.3.4\"],\"labels\":{\"__meta_bosh_process\":\"fake-process-name\"}}]"

			vmProcesses   []director.VMInfoProcess
			instanceInfos []director.VMInfo
			deployments   []director.Deployment
			deployment    director.Deployment

			metrics chan prometheus.Metric
		)

		BeforeEach(func() {
			vmProcesses = []director.VMInfoProcess{
				{
					Name:  jobProcessName,
					State: jobProcessState,
				},
			}

			instanceInfos = []director.VMInfo{
				{
					JobName:      jobName,
					ID:           jobID,
					Index:        &jobIndex,
					ProcessState: processState,
					IPs:          []string{jobIP},
					AZ:           jobAZ,
					Processes:    vmProcesses,
				},
			}

			deployment = &fakes.FakeDeployment{
				NameStub:          func() string { return deploymentName },
				InstanceInfosStub: func() ([]director.VMInfo, error) { return instanceInfos, nil },
			}

			deployments = []director.Deployment{deployment}
			boshClient.DeploymentsReturns(deployments, nil)

			metrics = make(chan prometheus.Metric)
		})

		JustBeforeEach(func() {
			go serviceDiscoveryCollector.Collect(metrics)
		})

		It("writes a target groups file", func() {
			Eventually(metrics).Should(Receive())
			targetGroups, err := ioutil.ReadFile(serviceDiscoveryFilename)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(targetGroups)).To(Equal(targetGroupsContent))
		})

		It("returns a last_service_discovery_scrape_timestamp & last_service_discovery_scrape_duration_seconds", func() {
			Eventually(metrics).Should(Receive())
			Eventually(metrics).Should(Receive())
			Consistently(metrics).ShouldNot(Receive())
		})

		Context("when there are no deployments", func() {
			BeforeEach(func() {
				boshClient.DeploymentsReturns([]director.Deployment{}, nil)
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
			})
		})

		Context("when it does not return any InstanceInfos", func() {
			BeforeEach(func() {
				deployment = &fakes.FakeDeployment{
					NameStub:          func() string { return deploymentName },
					InstanceInfosStub: func() ([]director.VMInfo, error) { return nil, nil },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
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
			})
		})

		Context("when it fails to get the InstanceInfos for a deployment", func() {
			BeforeEach(func() {
				deployment = &fakes.FakeDeployment{
					NameStub:          func() string { return deploymentName },
					InstanceInfosStub: func() ([]director.VMInfo, error) { return nil, errors.New("no InstanceInfo") },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
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
			})
		})

	})
})
