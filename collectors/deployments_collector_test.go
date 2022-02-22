package collectors_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/bosh-prometheus/bosh_exporter/deployments"

	. "github.com/bosh-prometheus/bosh_exporter/collectors"
	. "github.com/bosh-prometheus/bosh_exporter/utils/test_matchers"
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
		deploymentsCollector *DeploymentsCollector

		deploymentReleaseInfoMetric                *prometheus.GaugeVec
		deploymentStemcellInfoMetric               *prometheus.GaugeVec
		deploymentInstancesMetric                  *prometheus.GaugeVec
		lastDeploymentsScrapeTimestampMetric       prometheus.Gauge
		lastDeploymentsScrapeDurationSecondsMetric prometheus.Gauge

		deploymentName  = "fake-deployment-name"
		releaseName     = "fake-release-name"
		releaseVersion  = "1.2.3"
		stemcellName    = "fake-stemcell-name"
		stemcellVersion = "4.5.6"
		stemcellOSName  = "fake-stemcell-os-name"
		vmTypeSmall     = "fake-vm-type-small"
		vmTypeMedium    = "fake-vm-type-medium"
		vmTypeLarge     = "fake-vm-type-large"
	)

	BeforeEach(func() {
		namespace = "test_exporter"
		environment = "test_environment"
		boshName = "test_bosh_name"
		boshUUID = "test_bosh_uuid"

		deploymentReleaseInfoMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "deployment",
				Name:      "release_info",
				Help:      "Labeled BOSH Deployment Release Info with a constant '1' value.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_release_name", "bosh_release_version"},
		)

		deploymentReleaseInfoMetric.WithLabelValues(
			deploymentName,
			releaseName,
			releaseVersion,
		).Set(float64(1))

		deploymentStemcellInfoMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "deployment",
				Name:      "stemcell_info",
				Help:      "Labeled BOSH Deployment Stemcell Info with a constant '1' value.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_stemcell_name", "bosh_stemcell_version", "bosh_stemcell_os_name"},
		)

		deploymentStemcellInfoMetric.WithLabelValues(
			deploymentName,
			stemcellName,
			stemcellVersion,
			stemcellOSName,
		).Set(float64(1))

		deploymentInstancesMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "deployment",
				Name:      "instances",
				Help:      "Number of instances in this deployment",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_vm_type"},
		)

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

		lastDeploymentsScrapeTimestampMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_deployments_scrape_timestamp",
				Help:      "Number of seconds since 1970 since last scrape of Deployments metrics from BOSH.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
		)

		lastDeploymentsScrapeDurationSecondsMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_deployments_scrape_duration_seconds",
				Help:      "Duration of the last scrape of Deployments metrics from BOSH.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
		)
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
