package filters_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bosh-prometheus/bosh_exporter/filters"
)

var _ = Describe("AZsFilter", func() {
	var (
		filter    []string
		azsFilter *AZsFilter
	)

	BeforeEach(func() {
		filter = []string{"fake-az-1", "fake-az-3"}
	})

	JustBeforeEach(func() {
		azsFilter = NewAZsFilter(filter)
	})

	Describe("Enabled", func() {
		Context("when az is enabled", func() {
			It("returns true", func() {
				Expect(azsFilter.Enabled("fake-az-1")).To(BeTrue())
			})
		})

		Context("when az is not enabled", func() {
			It("returns false", func() {
				Expect(azsFilter.Enabled("fake-az-2")).To(BeFalse())
			})
		})

		Context("when there is no filter", func() {
			BeforeEach(func() {
				filter = []string{}
			})

			It("returns true", func() {
				Expect(azsFilter.Enabled("fake-az-2")).To(BeTrue())
			})
		})

		Context("when a filter has leading and/or trailing whitespaces", func() {
			BeforeEach(func() {
				filter = []string{"   fake-az-1  "}
			})

			It("returns true", func() {
				Expect(azsFilter.Enabled("fake-az-1")).To(BeTrue())
			})
		})
	})
})
