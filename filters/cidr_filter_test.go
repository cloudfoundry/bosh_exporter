package filters_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/bosh-prometheus/bosh_exporter/filters"
)

var _ = Describe("Cidr Filter", func() {
	var (
		err        error
		cidrs      []string
		cidrFilter *CidrFilter
	)

	JustBeforeEach(func() {
		cidrFilter, err = NewCidrFilter(cidrs)
	})

	Describe("New", func() {
		Context("when valid cidr", func() {
			BeforeEach(func() {
				cidrs = []string{"0.0.0.0/0", "10.250.0.0/16"}
			})

			It("does not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when invalid cidr", func() {
			BeforeEach(func() {
				cidrs = []string{"not.a.cidr"}
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid CIDR address: not.a.cidr"))
			})
		})
	})

	Describe("Select", func() {
		Describe("with default cidr", func() {
			BeforeEach(func() {
				cidrs = []string{"0.0.0.0/0"}
			})

			Context("when selecting single ip", func() {
				It("returns first ip/true", func() {
					ip, found := cidrFilter.Select([]string{"192.168.0.1"})
					Expect(found).To(BeTrue())
					Expect(ip).To(Equal("192.168.0.1"))
				})
			})

			Context("when selecting multiple ips", func() {
				It("returns first ip/true", func() {
					ip, found := cidrFilter.Select([]string{"192.168.0.1", "10.254.12.57"})
					Expect(found).To(BeTrue())
					Expect(ip).To(Equal("192.168.0.1"))
				})
			})

			Context("when selecting empty list", func() {
				It("returns empty/false", func() {
					ip, found := cidrFilter.Select([]string{})
					Expect(found).To(BeFalse())
					Expect(ip).To(Equal(""))
				})
			})
		})

		Describe("with multiple cidr", func() {
			BeforeEach(func() {
				cidrs = []string{"10.254.0.0/16", "0.0.0.0/0"}
			})

			Context("when selecting single ip", func() {
				It("returns first ip/true", func() {
					ip, found := cidrFilter.Select([]string{"192.168.0.1"})
					Expect(found).To(BeTrue())
					Expect(ip).To(Equal("192.168.0.1"))
				})
			})

			Context("when selecting multiple ips", func() {
				It("returns second ip/true", func() {
					ip, found := cidrFilter.Select([]string{"192.168.0.1", "10.254.12.57"})
					Expect(found).To(BeTrue())
					Expect(ip).To(Equal("10.254.12.57"))
				})
			})
		})

		Describe("with specific cidr", func() {
			BeforeEach(func() {
				cidrs = []string{"10.254.0.0/16"}
			})

			Context("with matching ip", func() {
				It("returns first ip/true", func() {
					ip, found := cidrFilter.Select([]string{"10.254.1.1"})
					Expect(found).To(BeTrue())
					Expect(ip).To(Equal("10.254.1.1"))
				})
			})

			Context("with unmatching ip", func() {
				It("returns empty/false", func() {
					ip, found := cidrFilter.Select([]string{"192.168.0.1"})
					Expect(found).To(BeFalse())
					Expect(ip).To(Equal(""))
				})
			})
		})
	})
})
