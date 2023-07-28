package collectors_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/bosh-prometheus/bosh_exporter/deployments"

	. "github.com/bosh-prometheus/bosh_exporter/collectors"
	. "github.com/bosh-prometheus/bosh_exporter/utils/matchers"
)

func init() {
	_ = log.Base().SetLevel("fatal")
}

var _ = Describe("DeploymentsCollector", func() {
	var (
		namespace            string
		environment          string
		boshName             string
		boshUUID             string
		metrics              *DeploymentsCollectorMetrics
		deploymentsCollector *DeploymentsCollector

		deploymentReleaseInfoMetric                *prometheus.GaugeVec
		deploymentReleaseJobInfoMetric             *prometheus.GaugeVec
		deploymentReleasePackageInfoMetric         *prometheus.GaugeVec
		deploymentStemcellInfoMetric               *prometheus.GaugeVec
		deploymentInstancesMetric                  *prometheus.GaugeVec
		lastDeploymentsScrapeTimestampMetric       prometheus.Gauge
		lastDeploymentsScrapeDurationSecondsMetric prometheus.Gauge

		deploymentName     = "fake-deployment-name"
		releaseName        = "fake-release-name"
		releaseVersion     = "1.2.3"
		releaseJobName     = "fake-release-job-name"
		releasePackageName = "fake-release-package-name"
		stemcellName       = "fake-stemcell-name"
		stemcellVersion    = "4.5.6"
		stemcellOSName     = "fake-stemcell-os-name"
		vmTypeSmall        = "fake-vm-type-small"
		vmTypeMedium       = "fake-vm-type-medium"
		vmTypeLarge        = "fake-vm-type-large"
	)

	BeforeEach(func() {
		namespace = testNamespace
		environment = testEnvironment
		boshName = testBoshName
		boshUUID = testBoshUUID
		metrics = NewDeploymentsCollectorMetrics(testNamespace, testEnvironment, testBoshName, testBoshUUID)

		deploymentReleaseInfoMetric = metrics.NewDeploymentReleaseInfoMetric()
		deploymentReleaseInfoMetric.WithLabelValues(
			deploymentName,
			releaseName,
			releaseVersion,
		).Set(float64(1))

		deploymentReleaseJobInfoMetric = metrics.NewDeploymentReleaseJobInfoMetric()
		deploymentReleaseJobInfoMetric.WithLabelValues(
			deploymentName,
			releaseName,
			releaseVersion,
			releaseJobName,
		).Set(float64(1))

		deploymentReleasePackageInfoMetric = metrics.NewDeploymentReleasePackageInfoMetric()
		deploymentReleasePackageInfoMetric.WithLabelValues(
			deploymentName,
			releaseName,
			releaseVersion,
			releasePackageName,
		).Set(float64(1))

		deploymentStemcellInfoMetric = metrics.NewDeploymentStemcellInfoMetric()
		deploymentStemcellInfoMetric.WithLabelValues(
			deploymentName,
			stemcellName,
			stemcellVersion,
			stemcellOSName,
		).Set(float64(1))

		deploymentInstancesMetric = metrics.NewDeploymentInstancesMetric()
		deploymentInstancesMetric.WithLabelValues(
			deploymentName,
			vmTypeSmall,
		).Set(float64(1))
		deploymentInstancesMetric.WithLabelValues(
			deploymentName,
			vmTypeMedium,
		).Set(float64(2))
		deploymentInstancesMetric.WithLabelValues(
			deploymentName,
			vmTypeLarge,
		).Set(float64(3))

		lastDeploymentsScrapeTimestampMetric = metrics.NewLastDeploymentsScrapeTimestampMetric()

		lastDeploymentsScrapeDurationSecondsMetric = metrics.NewLastDeploymentsScrapeDurationSecondsMetric()
	})

	JustBeforeEach(func() {
		deploymentsCollector = NewDeploymentsCollector(
			namespace,
			environment,
			boshName,
			boshUUID,
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
			go deploymentsCollector.Describe(descriptions)
		})

		It("returns a deployment_release_info description", func() {
			Eventually(descriptions).Should(Receive(Equal(deploymentReleaseInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
			).Desc())))
		})

		It("returns a deployment_release_job_info description", func() {
			Eventually(descriptions).Should(Receive(Equal(deploymentReleaseJobInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
				releaseJobName,
			).Desc())))
		})

		It("returns a deployment_release_package_info description", func() {
			Eventually(descriptions).Should(Receive(Equal(deploymentReleasePackageInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
				releasePackageName,
			).Desc())))
		})

		It("returns a deployment_stemcell_info metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(deploymentStemcellInfoMetric.WithLabelValues(
				deploymentName,
				stemcellName,
				stemcellVersion,
				stemcellOSName,
			).Desc())))
		})

		It("returns a deployment_instances metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(deploymentInstancesMetric.WithLabelValues(
				deploymentName,
				vmTypeSmall,
			).Desc())))
		})

		It("returns a last_deployments_scrape_timestamp metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastDeploymentsScrapeTimestampMetric.Desc())))
		})

		It("returns a last_deployments_scrape_duration_seconds metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastDeploymentsScrapeDurationSecondsMetric.Desc())))
		})
	})

	Describe("Collect", func() {
		var (
			release = deployments.Release{
				Name:         releaseName,
				Version:      releaseVersion,
				JobNames:     []string{releaseJobName},
				PackageNames: []string{releasePackageName},
			}
			releases = []deployments.Release{release}

			stemcell = deployments.Stemcell{
				Name:    stemcellName,
				Version: stemcellVersion,
				OSName:  stemcellOSName,
			}
			stemcells = []deployments.Stemcell{stemcell}

			instances = []deployments.Instance{
				{VMType: vmTypeSmall},
				{VMType: vmTypeMedium},
				{VMType: vmTypeMedium},
				{VMType: vmTypeLarge},
				{VMType: vmTypeLarge},
				{VMType: vmTypeLarge},
			}

			deploymentInfo deployments.DeploymentInfo

			deploymentsInfo []deployments.DeploymentInfo

			metrics    chan prometheus.Metric
			errMetrics chan error
		)

		BeforeEach(func() {
			deploymentInfo = deployments.DeploymentInfo{
				Name:      deploymentName,
				Releases:  releases,
				Stemcells: stemcells,
				Instances: instances,
			}
			deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}

			metrics = make(chan prometheus.Metric)
			errMetrics = make(chan error, 1)
		})

		JustBeforeEach(func() {
			go func() {
				if err := deploymentsCollector.Collect(deploymentsInfo, metrics); err != nil {
					errMetrics <- err
				}
			}()
		})

		It("returns a deployment_release_info metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(deploymentReleaseInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		It("returns a deployment_release_job_info metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(deploymentReleaseJobInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
				releaseJobName,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		It("returns a deployment_release_package_info metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(deploymentReleasePackageInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
				releasePackageName,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		It("returns a deployment_stemcell_info metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(deploymentStemcellInfoMetric.WithLabelValues(
				deploymentName,
				stemcellName,
				stemcellVersion,
				stemcellOSName,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		It("returns a deployment_instances for small vmType instance", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(deploymentInstancesMetric.WithLabelValues(
				deploymentName,
				vmTypeSmall,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		It("returns a deployment_instances for medium vmType instance", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(deploymentInstancesMetric.WithLabelValues(
				deploymentName,
				vmTypeMedium,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		It("returns a deployment_instances for large vmType instance", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(deploymentInstancesMetric.WithLabelValues(
				deploymentName,
				vmTypeLarge,
			))))
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
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(deploymentReleaseInfoMetric.WithLabelValues(
					deploymentName,
					releaseName,
					releaseVersion,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		Context("when there are no stemcells", func() {
			BeforeEach(func() {
				deploymentInfo.Stemcells = []deployments.Stemcell{}
				deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}
			})

			It("should not return a deployment_stemcell_info metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(deploymentStemcellInfoMetric.WithLabelValues(
					deploymentName,
					stemcellName,
					stemcellVersion,
					stemcellOSName,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		Context("when there are no instances", func() {
			BeforeEach(func() {
				deploymentInfo.Instances = []deployments.Instance{}
				deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}
			})

			It("should not return a deployment_instances metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(deploymentInstancesMetric.WithLabelValues(
					deploymentName,
					vmTypeSmall,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})
	})
})
