package collectors_test

import (
	"errors"
	"flag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/director/fakes"
	"github.com/cppforlife/go-semi-semantic/version"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/filters"

	. "github.com/cloudfoundry-community/bosh_exporter/collectors"
)

func init() {
	flag.Set("log.level", "fatal")
}

var _ = Describe("DeploymentsCollector", func() {
	var (
		namespace            string
		boshDeployments      []string
		deploymentsFilter    *filters.DeploymentsFilter
		boshClient           *fakes.FakeDirector
		deploymentsCollector *DeploymentsCollector

		deploymentReleaseInfoDesc                *prometheus.Desc
		deploymentStemcellInfoDesc               *prometheus.Desc
		lastDeploymentsScrapeTimestampDesc       *prometheus.Desc
		lastDeploymentsScrapeDurationSecondsDesc *prometheus.Desc
	)

	BeforeEach(func() {
		namespace = "test_exporter"
		boshDeployments = []string{}
		boshClient = &fakes.FakeDirector{}

		deploymentReleaseInfoDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "deployment", "release_info"),
			"BOSH Deployment Release Info.",
			[]string{"bosh_deployment", "bosh_release_name", "bosh_release_version"},
			nil,
		)

		deploymentStemcellInfoDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "deployment", "stemcell_info"),
			"BOSH Deployment Stemcell Info.",
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
		deploymentsFilter = filters.NewDeploymentsFilter(boshDeployments, boshClient)
		deploymentsCollector = NewDeploymentsCollector(namespace, *deploymentsFilter)
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

			release = &fakes.FakeRelease{
				NameStub:    func() string { return releaseName },
				VersionStub: func() version.Version { return version.MustNewVersionFromString(releaseVersion) },
			}
			releases = []director.Release{release}

			stemcell = &fakes.FakeStemcell{
				NameStub:    func() string { return stemcellName },
				VersionStub: func() version.Version { return version.MustNewVersionFromString(stemcellVersion) },
				OSNameStub:  func() string { return stemcellOSName },
			}
			stemcells = []director.Stemcell{stemcell}

			deployments []director.Deployment
			deployment  director.Deployment

			metrics                      chan prometheus.Metric
			deploymentReleaseInfoMetric  prometheus.Metric
			deploymentStemcellInfoMetric prometheus.Metric
		)

		BeforeEach(func() {
			deployment = &fakes.FakeDeployment{
				NameStub:      func() string { return deploymentName },
				ReleasesStub:  func() ([]director.Release, error) { return releases, nil },
				StemcellsStub: func() ([]director.Stemcell, error) { return stemcells, nil },
			}

			deployments = []director.Deployment{deployment}
			boshClient.DeploymentsReturns(deployments, nil)

			metrics = make(chan prometheus.Metric)

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
			go deploymentsCollector.Collect(metrics)
		})

		It("returns a deployment_release_info metric", func() {
			Eventually(metrics).Should(Receive(Equal(deploymentReleaseInfoMetric)))
		})

		It("returns a deployment_stemcell_info metric", func() {
			Eventually(metrics).Should(Receive(Equal(deploymentStemcellInfoMetric)))
		})

		Context("when there are no deployments", func() {
			BeforeEach(func() {
				boshClient.DeploymentsReturns([]director.Deployment{}, nil)
			})

			It("returns only a last_deployments_scrape_timestamp & last_deployments_scrape_duration_seconds metric", func() {
				Eventually(metrics).Should(Receive())
				Eventually(metrics).Should(Receive())
				Consistently(metrics).ShouldNot(Receive())
			})
		})

		Context("when it does not return any Release", func() {
			BeforeEach(func() {
				deployment = &fakes.FakeDeployment{
					NameStub:      func() string { return deploymentName },
					StemcellsStub: func() ([]director.Stemcell, error) { return stemcells, nil },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("should not return a deployment_release_info metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(deploymentReleaseInfoMetric)))
			})
		})

		Context("when it fails to get the Releases for a deployment", func() {
			BeforeEach(func() {
				deployment = &fakes.FakeDeployment{
					NameStub:      func() string { return deploymentName },
					ReleasesStub:  func() ([]director.Release, error) { return nil, errors.New("no Releases") },
					StemcellsStub: func() ([]director.Stemcell, error) { return stemcells, nil },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("should not return a deployment_release_info metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(deploymentReleaseInfoMetric)))
			})
		})

		Context("when it does not return any Stemcell", func() {
			BeforeEach(func() {
				deployment = &fakes.FakeDeployment{
					NameStub:     func() string { return deploymentName },
					ReleasesStub: func() ([]director.Release, error) { return releases, nil },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("should not return a deployment_stemcell_info metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(deploymentStemcellInfoMetric)))
			})
		})

		Context("when it fails to get the Stemcells for a deployment", func() {
			BeforeEach(func() {
				deployment = &fakes.FakeDeployment{
					NameStub:      func() string { return deploymentName },
					ReleasesStub:  func() ([]director.Release, error) { return releases, nil },
					StemcellsStub: func() ([]director.Stemcell, error) { return nil, errors.New("no Stemcells") },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("should not return a deployment_stemcell_info metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(deploymentStemcellInfoMetric)))
			})
		})
	})
})
