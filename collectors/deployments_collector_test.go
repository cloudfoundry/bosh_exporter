package collectors_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/deployments"

	. "github.com/cloudfoundry-community/bosh_exporter/collectors"
)

var _ = Describe("DeploymentsCollector", func() {
	var (
		namespace            string
		deploymentsCollector *DeploymentsCollector

		deploymentReleaseInfoDesc                *prometheus.Desc
		deploymentStemcellInfoDesc               *prometheus.Desc
		lastDeploymentsScrapeTimestampDesc       *prometheus.Desc
		lastDeploymentsScrapeDurationSecondsDesc *prometheus.Desc
	)

	BeforeEach(func() {
		namespace = "test_exporter"

		deploymentReleaseInfoDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "deployment", "release_info"),
			"Labeled BOSH Deployment Release Info with a constant '1' value.",
			[]string{"bosh_deployment", "bosh_release_name", "bosh_release_version"},
			nil,
		)

		deploymentStemcellInfoDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "deployment", "stemcell_info"),
			"Labeled BOSH Deployment Stemcell Info with a constant '1' value.",
			[]string{"bosh_deployment", "bosh_stemcell_name", "bosh_stemcell_version", "bosh_stemcell_os_name"},
			nil,
		)

		lastDeploymentsScrapeTimestampDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "last_deployments_scrape_timestamp"),
			"Number of seconds since 1970 since last scrape of Deployments metrics from BOSH.",
			[]string{},
			nil,
		)

		lastDeploymentsScrapeDurationSecondsDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "last_deployments_scrape_duration_seconds"),
			"Duration of the last scrape of Deployments metrics from BOSH.",
			[]string{},
			nil,
		)
	})

	JustBeforeEach(func() {
		deploymentsCollector = NewDeploymentsCollector(namespace)
	})

	Describe("Describe", func() {
		var (
			descriptions chan *prometheus.Desc
		)

		BeforeEach(func() {
			descriptions = make(chan *prometheus.Desc)
		})

		JustBeforeEach(func() {
			go deploymentsCollector.Describe(descriptions)
		})

		It("returns a deployment_release_info description", func() {
			Eventually(descriptions).Should(Receive(Equal(deploymentReleaseInfoDesc)))
		})

		It("returns a deployment_stemcell_info metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(deploymentStemcellInfoDesc)))
		})

		It("returns a last_deployments_scrape_timestamp metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastDeploymentsScrapeTimestampDesc)))
		})

		It("returns a last_deployments_scrape_duration_seconds metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastDeploymentsScrapeDurationSecondsDesc)))
		})
	})

	Describe("Collect", func() {
		var (
			deploymentName  = "fake-deployment-name"
			releaseName     = "fake-release-name"
			releaseVersion  = "1.2.3"
			stemcellName    = "fake-stemcell-name"
			stemcellVersion = "4.5.6"
			stemcellOSName  = "fake-stemcell-os-name"

			release = deployments.Release{
				Name:    releaseName,
				Version: releaseVersion,
			}
			releases = []deployments.Release{release}

			stemcell = deployments.Stemcell{
				Name:    stemcellName,
				Version: stemcellVersion,
				OSName:  stemcellOSName,
			}
			stemcells = []deployments.Stemcell{stemcell}

			deploymentInfo deployments.DeploymentInfo

			deploymentsInfo []deployments.DeploymentInfo

			metrics                      chan prometheus.Metric
			errMetrics                   chan error
			deploymentReleaseInfoMetric  prometheus.Metric
			deploymentStemcellInfoMetric prometheus.Metric
		)

		BeforeEach(func() {
			deploymentInfo = deployments.DeploymentInfo{
				Name:      deploymentName,
				Releases:  releases,
				Stemcells: stemcells,
			}
			deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}

			metrics = make(chan prometheus.Metric)
			errMetrics = make(chan error, 1)

			deploymentReleaseInfoMetric = prometheus.MustNewConstMetric(
				deploymentReleaseInfoDesc,
				prometheus.GaugeValue,
				float64(1),
				deploymentName,
				releaseName,
				releaseVersion,
			)

			deploymentStemcellInfoMetric = prometheus.MustNewConstMetric(
				deploymentStemcellInfoDesc,
				prometheus.GaugeValue,
				float64(1),
				deploymentName,
				stemcellName,
				stemcellVersion,
				stemcellOSName,
			)
		})

		JustBeforeEach(func() {
			go func() {
				if err := deploymentsCollector.Collect(deploymentsInfo, metrics); err != nil {
					errMetrics <- err
				}
			}()
		})

		It("returns a deployment_release_info metric", func() {
			Eventually(metrics).Should(Receive(Equal(deploymentReleaseInfoMetric)))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		It("returns a deployment_stemcell_info metric", func() {
			Eventually(metrics).Should(Receive(Equal(deploymentStemcellInfoMetric)))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there are no deployments", func() {
			BeforeEach(func() {
				deploymentsInfo = []deployments.DeploymentInfo{}
			})

			It("returns only a last_deployments_scrape_timestamp & last_deployments_scrape_duration_seconds metric", func() {
				Eventually(metrics).Should(Receive())
				Eventually(metrics).Should(Receive())
				Consistently(metrics).ShouldNot(Receive())
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		Context("when there are no releases", func() {
			BeforeEach(func() {
				deploymentInfo.Releases = []deployments.Release{}
				deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}
			})

			It("should not return a deployment_release_info metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(deploymentReleaseInfoMetric)))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		Context("when there are no stemcells", func() {
			BeforeEach(func() {
				deploymentInfo.Stemcells = []deployments.Stemcell{}
				deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}
			})

			It("should not return a deployment_stemcell_info metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(deploymentStemcellInfoMetric)))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})
	})
})
