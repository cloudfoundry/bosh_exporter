package collectors_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry/bosh_exporter/deployments"

	"github.com/cloudfoundry/bosh_exporter/collectors"
	"github.com/cloudfoundry/bosh_exporter/utils/matchers"
)

var _ = ginkgo.Describe("DeploymentsCollector", func() {
	var (
		namespace            string
		environment          string
		boshName             string
		boshUUID             string
		metrics              *collectors.DeploymentsCollectorMetrics
		deploymentsCollector *collectors.DeploymentsCollector

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

	ginkgo.BeforeEach(func() {
		namespace = testNamespace
		environment = testEnvironment
		boshName = testBoshName
		boshUUID = testBoshUUID
		metrics = collectors.NewDeploymentsCollectorMetrics(testNamespace, testEnvironment, testBoshName, testBoshUUID)

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

	ginkgo.JustBeforeEach(func() {
		deploymentsCollector = collectors.NewDeploymentsCollector(
			namespace,
			environment,
			boshName,
			boshUUID,
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
			go deploymentsCollector.Describe(descriptions)
		})

		ginkgo.It("returns a deployment_release_info description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(deploymentReleaseInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
			).Desc())))
		})

		ginkgo.It("returns a deployment_release_job_info description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(deploymentReleaseJobInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
				releaseJobName,
			).Desc())))
		})

		ginkgo.It("returns a deployment_release_package_info description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(deploymentReleasePackageInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
				releasePackageName,
			).Desc())))
		})

		ginkgo.It("returns a deployment_stemcell_info metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(deploymentStemcellInfoMetric.WithLabelValues(
				deploymentName,
				stemcellName,
				stemcellVersion,
				stemcellOSName,
			).Desc())))
		})

		ginkgo.It("returns a deployment_instances metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(deploymentInstancesMetric.WithLabelValues(
				deploymentName,
				vmTypeSmall,
			).Desc())))
		})

		ginkgo.It("returns a last_deployments_scrape_timestamp metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(lastDeploymentsScrapeTimestampMetric.Desc())))
		})

		ginkgo.It("returns a last_deployments_scrape_duration_seconds metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(lastDeploymentsScrapeDurationSecondsMetric.Desc())))
		})
	})

	ginkgo.Describe("Collect", func() {
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

		ginkgo.BeforeEach(func() {
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

		ginkgo.JustBeforeEach(func() {
			go func() {
				if err := deploymentsCollector.Collect(deploymentsInfo, metrics); err != nil {
					errMetrics <- err
				}
			}()
		})

		ginkgo.It("returns a deployment_release_info metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(deploymentReleaseInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
			))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.It("returns a deployment_release_job_info metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(deploymentReleaseJobInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
				releaseJobName,
			))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.It("returns a deployment_release_package_info metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(deploymentReleasePackageInfoMetric.WithLabelValues(
				deploymentName,
				releaseName,
				releaseVersion,
				releasePackageName,
			))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.It("returns a deployment_stemcell_info metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(deploymentStemcellInfoMetric.WithLabelValues(
				deploymentName,
				stemcellName,
				stemcellVersion,
				stemcellOSName,
			))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.It("returns a deployment_instances for small vmType instance", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(deploymentInstancesMetric.WithLabelValues(
				deploymentName,
				vmTypeSmall,
			))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.It("returns a deployment_instances for medium vmType instance", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(deploymentInstancesMetric.WithLabelValues(
				deploymentName,
				vmTypeMedium,
			))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.It("returns a deployment_instances for large vmType instance", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(deploymentInstancesMetric.WithLabelValues(
				deploymentName,
				vmTypeLarge,
			))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there are no deployments", func() {
			ginkgo.BeforeEach(func() {
				deploymentsInfo = []deployments.DeploymentInfo{}
			})

			ginkgo.It("returns only a last_deployments_scrape_timestamp & last_deployments_scrape_duration_seconds metric", func() {
				gomega.Eventually(metrics).Should(gomega.Receive())
				gomega.Eventually(metrics).Should(gomega.Receive())
				gomega.Consistently(metrics).ShouldNot(gomega.Receive())
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.Context("when there are no releases", func() {
			ginkgo.BeforeEach(func() {
				deploymentInfo.Releases = []deployments.Release{}
				deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}
			})

			ginkgo.It("should not return a deployment_release_info metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(deploymentReleaseInfoMetric.WithLabelValues(
					deploymentName,
					releaseName,
					releaseVersion,
				))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.Context("when there are no stemcells", func() {
			ginkgo.BeforeEach(func() {
				deploymentInfo.Stemcells = []deployments.Stemcell{}
				deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}
			})

			ginkgo.It("should not return a deployment_stemcell_info metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(deploymentStemcellInfoMetric.WithLabelValues(
					deploymentName,
					stemcellName,
					stemcellVersion,
					stemcellOSName,
				))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.Context("when there are no instances", func() {
			ginkgo.BeforeEach(func() {
				deploymentInfo.Instances = []deployments.Instance{}
				deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}
			})

			ginkgo.It("should not return a deployment_instances metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(deploymentInstancesMetric.WithLabelValues(
					deploymentName,
					vmTypeSmall,
				))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})
	})
})
