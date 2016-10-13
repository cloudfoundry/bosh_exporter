package collectors_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director/fakes"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/collectors"
)

var _ = Describe("ProcessesCollector", func() {
	var (
		namespace          string
		boshDeployments    []string
		boshClient         *fakes.FakeDirector
		processesCollector *collectors.ProcessesCollector
	)

	BeforeEach(func() {
		namespace = "test_exporter"
		boshDeployments = []string{}
		boshClient = &fakes.FakeDirector{}
		processesCollector = collectors.NewProcessesCollector(namespace, boshDeployments, boshClient)
	})

	Describe("Describe", func() {
		var (
			descriptions          chan *prometheus.Desc
			processHealthyDesc    *prometheus.Desc
			processUptimeDesc     *prometheus.Desc
			processCPUTotalDesc   *prometheus.Desc
			processMemKBDesc      *prometheus.Desc
			processMemPercentDesc *prometheus.Desc
		)

		BeforeEach(func() {
			descriptions = make(chan *prometheus.Desc)

			processHealthyDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_process_healthy"),
				"BOSH Job Process Healthy.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_process"},
				nil,
			)

			processUptimeDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_process_uptime_seconds"),
				"BOSH Job Process Uptime in seconds.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_process"},
				nil,
			)

			processCPUTotalDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_process_cpu_total"),
				"BOSH Job Process CPU Total.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_process"},
				nil,
			)

			processMemKBDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_process_mem_kb"),
				"BOSH Job Process Memory KB.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_process"},
				nil,
			)

			processMemPercentDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_process_mem_percent"),
				"BOSH Job Process Memory Percent.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_process"},
				nil,
			)

			go processesCollector.Describe(descriptions)
		})

		It("returns a bosh_job_process_healthy metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(processHealthyDesc)))
		})

		It("returns a bosh_job_process_uptime_seconds metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(processUptimeDesc)))
		})

		It("returns a bosh_job_process_cpu_total metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(processCPUTotalDesc)))
		})

		It("returns a bosh_job_process_mem_kb metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(processMemKBDesc)))
		})

		It("returns a bosh_job_process_mem_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(processMemPercentDesc)))
		})
	})
})
