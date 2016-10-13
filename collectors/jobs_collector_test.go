package collectors_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director/fakes"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/collectors"
)

var _ = Describe("JobsCollector", func() {
	var (
		namespace       string
		boshDeployments []string
		boshClient      *fakes.FakeDirector
		jobsCollector   *collectors.JobsCollector
	)

	BeforeEach(func() {
		namespace = "test_exporter"
		boshDeployments = []string{}
		boshClient = &fakes.FakeDirector{}
		jobsCollector = collectors.NewJobsCollector(namespace, boshDeployments, boshClient)
	})

	Describe("Describe", func() {
		var (
			descriptions                      chan *prometheus.Desc
			jobHealthyDesc                    *prometheus.Desc
			jobLoadAvg01Desc                  *prometheus.Desc
			jobLoadAvg05Desc                  *prometheus.Desc
			jobLoadAvg15Desc                  *prometheus.Desc
			jobCPUSysDesc                     *prometheus.Desc
			jobCPUUserDesc                    *prometheus.Desc
			jobCPUWaitDesc                    *prometheus.Desc
			jobMemKBDesc                      *prometheus.Desc
			jobMemPercentDesc                 *prometheus.Desc
			jobSwapKBDesc                     *prometheus.Desc
			jobSwapPercentDesc                *prometheus.Desc
			jobSystemDiskInodePercentDesc     *prometheus.Desc
			jobSystemDiskPercentDesc          *prometheus.Desc
			jobEphemeralDiskInodePercentDesc  *prometheus.Desc
			jobEphemeralDiskPercentDesc       *prometheus.Desc
			jobPersistentDiskInodePercentDesc *prometheus.Desc
			jobPersistentDiskPercentDesc      *prometheus.Desc
		)

		BeforeEach(func() {
			descriptions = make(chan *prometheus.Desc)

			jobHealthyDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_healthy"),
				"BOSH Job Healthy.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobLoadAvg01Desc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_load_avg01"),
				"BOSH Job Load avg01.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobLoadAvg05Desc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_load_avg05"),
				"BOSH Job Load avg05.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobLoadAvg15Desc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_load_avg15"),
				"BOSH Job Load avg15.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobCPUSysDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_cpu_sys"),
				"BOSH Job CPU System.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobCPUUserDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_cpu_user"),
				"BOSH Job CPU User.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobCPUWaitDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_cpu_wait"),
				"BOSH Job CPU Wait.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobMemKBDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_mem_kb"),
				"BOSH Job Memory KB.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobMemPercentDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_mem_percent"),
				"BOSH Job Memory Percent.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobSwapKBDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_swap_kb"),
				"BOSH Job Swap KB.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobSwapPercentDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_swap_percent"),
				"BOSH Job Swap Percent.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobSystemDiskInodePercentDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_system_disk_inode_percent"),
				"BOSH Job System Disk Inode Percent.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobSystemDiskPercentDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_system_disk_percent"),
				"BOSH Job System Disk Percent.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobEphemeralDiskInodePercentDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_ephemeral_disk_inode_percent"),
				"BOSH Job Ephemeral Disk Inode Percent.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobEphemeralDiskPercentDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_ephemeral_disk_percent"),
				"BOSH Job Ephemeral Disk Percent.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobPersistentDiskInodePercentDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_persistent_disk_inode_percent"),
				"BOSH Job Persistent Disk Inode Percent.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			jobPersistentDiskPercentDesc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "bosh", "job_persistent_disk_percent"),
				"BOSH Job Persistent Disk Percent.",
				[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az"},
				nil,
			)

			go jobsCollector.Describe(descriptions)
		})

		It("returns a bosh_job_healthy metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobHealthyDesc)))
		})

		It("returns a bosh_job_load_avg01 metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobLoadAvg01Desc)))
		})

		It("returns a bosh_job_load_avg05 metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobLoadAvg05Desc)))
		})

		It("returns a bosh_job_load_avg15 metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobLoadAvg15Desc)))
		})

		It("returns a bosh_job_cpu_sys metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobCPUSysDesc)))
		})

		It("returns a bosh_job_cpu_user metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobCPUUserDesc)))
		})

		It("returns a bosh_job_cpu_wait metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobCPUWaitDesc)))
		})

		It("returns a bosh_job_mem_kb metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobMemKBDesc)))
		})

		It("returns a bosh_job_mem_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobMemPercentDesc)))
		})

		It("returns a bosh_job_swap_kb metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobSwapKBDesc)))
		})

		It("returns a bosh_job_swap_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobSwapPercentDesc)))
		})

		It("returns a bosh_job_system_disk_inode_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobSystemDiskInodePercentDesc)))
		})

		It("returns a bosh_job_system_disk_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobSystemDiskPercentDesc)))
		})

		It("returns a bosh_job_ephemeral_disk_inode_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobEphemeralDiskInodePercentDesc)))
		})

		It("returns a bosh_job_ephemeral_disk_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobEphemeralDiskPercentDesc)))
		})

		It("returns a bosh_job_persistent_disk_inode_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobPersistentDiskInodePercentDesc)))
		})

		It("returns a bosh_job_persistent_disk_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobPersistentDiskPercentDesc)))
		})
	})
})
