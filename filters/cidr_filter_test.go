package filters_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh_exporter/filters"
)

var _ = ginkgo.Describe("Cidr Filter", func() {
	var (
		err        error
		cidrs      []string
		cidrFilter *filters.CidrFilter
	)

	ginkgo.JustBeforeEach(func() {
		cidrFilter, err = filters.NewCidrFilter(cidrs)
	})

	ginkgo.Describe("New", func() {
		ginkgo.Context("when valid cidr", func() {
			ginkgo.BeforeEach(func() {
				cidrs = []string{"0.0.0.0/0", "10.250.0.0/16"}
			})

			ginkgo.It("does not return an error", func() {
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when invalid cidr", func() {
			ginkgo.BeforeEach(func() {
				cidrs = []string{"not.a.cidr"}
			})

			ginkgo.It("returns an error", func() {
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.Equal("invalid CIDR address: not.a.cidr"))
			})
		})
	})

	ginkgo.Describe("Select", func() {
		ginkgo.Describe("with default cidr", func() {
			ginkgo.BeforeEach(func() {
				cidrs = []string{"0.0.0.0/0"}
			})

			ginkgo.Context("when selecting single ip", func() {
				ginkgo.It("returns first ip/true", func() {
					ip, found := cidrFilter.Select([]string{"192.168.0.1"})
					gomega.Expect(found).To(gomega.BeTrue())
					gomega.Expect(ip).To(gomega.Equal("192.168.0.1"))
				})
			})

			ginkgo.Context("when selecting multiple ips", func() {
				ginkgo.It("returns first ip/true", func() {
					ip, found := cidrFilter.Select([]string{"192.168.0.1", "10.254.12.57"})
					gomega.Expect(found).To(gomega.BeTrue())
					gomega.Expect(ip).To(gomega.Equal("192.168.0.1"))
				})
			})

			ginkgo.Context("when selecting empty list", func() {
				ginkgo.It("returns empty/false", func() {
					ip, found := cidrFilter.Select([]string{})
					gomega.Expect(found).To(gomega.BeFalse())
					gomega.Expect(ip).To(gomega.Equal(""))
				})
			})
		})

		ginkgo.Describe("with multiple cidr", func() {
			ginkgo.BeforeEach(func() {
				cidrs = []string{"10.254.0.0/16", "0.0.0.0/0"}
			})

			ginkgo.Context("when selecting single ip", func() {
				ginkgo.It("returns first ip/true", func() {
					ip, found := cidrFilter.Select([]string{"192.168.0.1"})
					gomega.Expect(found).To(gomega.BeTrue())
					gomega.Expect(ip).To(gomega.Equal("192.168.0.1"))
				})
			})

			ginkgo.Context("when selecting multiple ips", func() {
				ginkgo.It("returns second ip/true", func() {
					ip, found := cidrFilter.Select([]string{"192.168.0.1", "10.254.12.57"})
					gomega.Expect(found).To(gomega.BeTrue())
					gomega.Expect(ip).To(gomega.Equal("10.254.12.57"))
				})
			})
		})

		ginkgo.Describe("with specific cidr", func() {
			ginkgo.BeforeEach(func() {
				cidrs = []string{"10.254.0.0/16"}
			})

			ginkgo.Context("with matching ip", func() {
				ginkgo.It("returns first ip/true", func() {
					ip, found := cidrFilter.Select([]string{"10.254.1.1"})
					gomega.Expect(found).To(gomega.BeTrue())
					gomega.Expect(ip).To(gomega.Equal("10.254.1.1"))
				})
			})

			ginkgo.Context("with unmatching ip", func() {
				ginkgo.It("returns empty/false", func() {
					ip, found := cidrFilter.Select([]string{"192.168.0.1"})
					gomega.Expect(found).To(gomega.BeFalse())
					gomega.Expect(ip).To(gomega.Equal(""))
				})
			})
		})
	})
})
