package filters_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bosh-prometheus/bosh_exporter/filters"
)

var _ = Describe("RegexpFilter", func() {
	var (
		err     error
		filters []string

		regexpFilter *RegexpFilter
	)

	JustBeforeEach(func() {
		regexpFilter, err = NewRegexpFilter(filters)
	})

	Describe("New", func() {
		Context("when filters compile", func() {
			BeforeEach(func() {
				filters = []string{"bosh_exporter", "[a-z]+_collector"}
			})

			It("does not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when filters does not compile", func() {
			BeforeEach(func() {
				filters = []string{"[a-(z]+_exporter"}
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error parsing regexp: invalid character class range: `a-(`"))
			})
		})
	})

	Describe("Enabled", func() {
		BeforeEach(func() {
			filters = []string{"bosh_exporter", "[a-z]+_collector"}
		})

		Context("when there is a match", func() {
			It("returns true", func() {
				Expect(regexpFilter.Enabled("deployments_collector")).To(BeTrue())
			})
		})

		Context("when there is not a match", func() {
			It("returns false", func() {
				Expect(regexpFilter.Enabled("deployments_exporter")).To(BeFalse())
			})
		})

		Context("when there are no filters", func() {
			BeforeEach(func() {
				filters = []string{}
			})

			It("returns true", func() {
				Expect(regexpFilter.Enabled("deployments_exporter")).To(BeTrue())
			})
		})
	})
})
