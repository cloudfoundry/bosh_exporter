package filters_test

import (
	"errors"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/director/directorfakes"
	"github.com/prometheus/common/log"

	"github.com/cloudfoundry/bosh_exporter/filters"
)

func init() {
	_ = log.Base().SetLevel("fatal")
}

var _ = ginkgo.Describe("DeploymentsFilter", func() {
	var (
		err               error
		filtersArray      []string
		boshClient        *directorfakes.FakeDirector
		deploymentsFilter *filters.DeploymentsFilter
	)

	ginkgo.Describe("GetDeployments", func() {
		var (
			deployment1    director.Deployment
			deployment2    director.Deployment
			allDeployments []director.Deployment

			deployments []director.Deployment
		)

		ginkgo.BeforeEach(func() {
			filtersArray = []string{}
			boshClient = &directorfakes.FakeDirector{}

			deployment1 = &directorfakes.FakeDeployment{
				NameStub: func() string { return "fake-deployment-name-1" },
			}
			deployment2 = &directorfakes.FakeDeployment{
				NameStub: func() string { return "fake-deployment-name-2" },
			}
			allDeployments = []director.Deployment{}
		})

		ginkgo.JustBeforeEach(func() {
			deploymentsFilter = filters.NewDeploymentsFilter(filtersArray, boshClient)
			deployments, err = deploymentsFilter.GetDeployments()
		})

		ginkgo.Context("when there are no filters", func() {
			ginkgo.BeforeEach(func() {
				boshClient.DeploymentsReturns(allDeployments, nil)
			})

			ginkgo.It("returns all deployments", func() {
				gomega.Expect(deployments).To(gomega.Equal(allDeployments))
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})

			ginkgo.Context("and there are no deployments", func() {
				ginkgo.BeforeEach(func() {
					boshClient.DeploymentsReturns([]director.Deployment{}, nil)
				})

				ginkgo.It("does not return any deployment", func() {
					gomega.Expect(deployments).To(gomega.BeEmpty())
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
				})
			})

			ginkgo.Context("and it fails to get the deployments", func() {
				ginkgo.BeforeEach(func() {
					boshClient.DeploymentsReturns(nil, errors.New("no deployments"))
				})

				ginkgo.It("does not return any deployment", func() {
					gomega.Expect(deployments).To(gomega.BeEmpty())
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when there are filters", func() {
			ginkgo.BeforeEach(func() {
				filtersArray = []string{"fake-deployment-name-1"}
				boshClient.FindDeploymentReturns(deployment1, nil)
			})

			ginkgo.It("returns the filtered deployments", func() {
				gomega.Expect(boshClient.FindDeploymentArgsForCall(0)).To(gomega.Equal("fake-deployment-name-1"))
				gomega.Expect(deployments).To(gomega.ContainElement(deployment1))
				gomega.Expect(deployments).ToNot(gomega.ContainElement(deployment2))
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})

			ginkgo.Context("and it fails to get the deployment", func() {
				ginkgo.BeforeEach(func() {
					boshClient.FindDeploymentReturns(nil, errors.New("deployment does not exists"))
				})

				ginkgo.It("does not return any deployment", func() {
					gomega.Expect(deployments).To(gomega.BeEmpty())
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})

			ginkgo.Context("and the deployment name has leading and/or trailing whitespaces", func() {
				ginkgo.BeforeEach(func() {
					filtersArray = []string{"   fake-deployment-name-1  "}
				})

				ginkgo.It("returns the filtered deployments", func() {
					gomega.Expect(boshClient.FindDeploymentArgsForCall(0)).To(gomega.Equal("fake-deployment-name-1"))
					gomega.Expect(deployments).To(gomega.ContainElement(deployment1))
					gomega.Expect(deployments).ToNot(gomega.ContainElement(deployment2))
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
				})
			})
		})
	})
})
