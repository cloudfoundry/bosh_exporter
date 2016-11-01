package filters_test

import (
	"errors"
	"flag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/director/fakes"

	. "github.com/cloudfoundry-community/bosh_exporter/filters"
)

func init() {
	flag.Set("log.level", "fatal")
}

var _ = Describe("DeploymentsFilter", func() {
	var (
		filter            []string
		boshClient        *fakes.FakeDirector
		deploymentsFilter *DeploymentsFilter
	)

	Describe("GetDeployments", func() {
		var (
			deployment1    director.Deployment
			deployment2    director.Deployment
			allDeployments []director.Deployment

			deployments []director.Deployment
		)

		BeforeEach(func() {
			filter = []string{}
			boshClient = &fakes.FakeDirector{}

			deployment1 = &fakes.FakeDeployment{
				NameStub: func() string { return "fake-deployment-name-1" },
			}
			deployment2 = &fakes.FakeDeployment{
				NameStub: func() string { return "fake-deployment-name-2" },
			}
			allDeployments = []director.Deployment{}
		})

		JustBeforeEach(func() {
			deploymentsFilter = NewDeploymentsFilter(filter, boshClient)
			deployments = deploymentsFilter.GetDeployments()
		})

		Context("when there are no filters", func() {
			BeforeEach(func() {
				boshClient.DeploymentsReturns(allDeployments, nil)
			})

			It("returns all deployments", func() {
				Expect(deployments).To(Equal(allDeployments))
			})

			Context("and there are no deployments", func() {
				BeforeEach(func() {
					boshClient.DeploymentsReturns([]director.Deployment{}, nil)
				})

				It("does not return any deployment", func() {
					Expect(deployments).To(BeEmpty())
				})
			})

			Context("and it fails to get the deployments", func() {
				BeforeEach(func() {
					boshClient.DeploymentsReturns(nil, errors.New("no deployments"))
				})

				It("does not return any deployment", func() {
					Expect(deployments).To(BeEmpty())
				})
			})
		})

		Context("when there are filters", func() {
			BeforeEach(func() {
				filter = []string{"fake-deployment-name-1"}
				boshClient.FindDeploymentReturns(deployment1, nil)
			})

			It("returns the filtered deployments", func() {
				Expect(boshClient.FindDeploymentArgsForCall(0)).To(Equal("fake-deployment-name-1"))
				Expect(deployments).To(ContainElement(deployment1))
				Expect(deployments).ToNot(ContainElement(deployment2))
			})

			Context("and it fails to get the deployment", func() {
				BeforeEach(func() {
					boshClient.FindDeploymentReturns(nil, errors.New("deployment does not exists"))
				})

				It("does not return any deployment", func() {
					Expect(deployments).To(BeEmpty())
				})
			})
		})
	})
})
