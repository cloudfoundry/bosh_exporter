package filters_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh_exporter/filters"
)

var _ = ginkgo.Describe("AZsFilter", func() {
	var (
		filter    []string
		azsFilter *filters.AZsFilter
	)

	ginkgo.BeforeEach(func() {
		filter = []string{"fake-az-1", "fake-az-3"}
	})

	ginkgo.JustBeforeEach(func() {
		azsFilter = filters.NewAZsFilter(filter)
	})

	ginkgo.Describe("Enabled", func() {
		ginkgo.Context("when az is enabled", func() {
			ginkgo.It("returns true", func() {
				gomega.Expect(azsFilter.Enabled("fake-az-1")).To(gomega.BeTrue())
			})
		})

		ginkgo.Context("when az is not enabled", func() {
			ginkgo.It("returns false", func() {
				gomega.Expect(azsFilter.Enabled("fake-az-2")).To(gomega.BeFalse())
			})
		})

		ginkgo.Context("when there is no filter", func() {
			ginkgo.BeforeEach(func() {
				filter = []string{}
			})

			ginkgo.It("returns true", func() {
				gomega.Expect(azsFilter.Enabled("fake-az-2")).To(gomega.BeTrue())
			})
		})

		ginkgo.Context("when a filter has leading and/or trailing whitespaces", func() {
			ginkgo.BeforeEach(func() {
				filter = []string{"   fake-az-1  "}
			})

			ginkgo.It("returns true", func() {
				gomega.Expect(azsFilter.Enabled("fake-az-1")).To(gomega.BeTrue())
			})
		})
	})
})
