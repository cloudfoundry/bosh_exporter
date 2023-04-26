package collectors_test

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/bosh-prometheus/bosh_exporter/deployments"
	"github.com/bosh-prometheus/bosh_exporter/filters"

	. "github.com/bosh-prometheus/bosh_exporter/collectors"
	. "github.com/bosh-prometheus/bosh_exporter/utils/matchers"
)

func init() {
	_ = log.Base().SetLevel("fatal")
}

var _ = Describe("JobsCollector", func() {
	var (
		err           error
		namespace     string
		environment   string
		boshName      string
		boshUUID      string
		azsFilter     *filters.AZsFilter
		cidrsFilter   *filters.CidrFilter
		jobsCollector *JobsCollector

		jobHealthyMetric                    *prometheus.GaugeVec
		jobLoadAvg01Metric                  *prometheus.GaugeVec
		jobLoadAvg05Metric                  *prometheus.GaugeVec
		jobLoadAvg15Metric                  *prometheus.GaugeVec
		jobCPUSysMetric                     *prometheus.GaugeVec
		jobCPUUserMetric                    *prometheus.GaugeVec
		jobCPUWaitMetric                    *prometheus.GaugeVec
		jobMemKBMetric                      *prometheus.GaugeVec
		jobMemPercentMetric                 *prometheus.GaugeVec
		jobSwapKBMetric                     *prometheus.GaugeVec
		jobSwapPercentMetric                *prometheus.GaugeVec
		jobSystemDiskInodePercentMetric     *prometheus.GaugeVec
		jobSystemDiskPercentMetric          *prometheus.GaugeVec
		jobEphemeralDiskInodePercentMetric  *prometheus.GaugeVec
		jobEphemeralDiskPercentMetric       *prometheus.GaugeVec
		jobPersistentDiskInodePercentMetric *prometheus.GaugeVec
		jobPersistentDiskPercentMetric      *prometheus.GaugeVec
		jobProcessHealthyMetric             *prometheus.GaugeVec
		jobProcessUptimeMetric              *prometheus.GaugeVec
		jobProcessCPUTotalMetric            *prometheus.GaugeVec
		jobProcessMemKBMetric               *prometheus.GaugeVec
		jobProcessMemPercentMetric          *prometheus.GaugeVec
		lastJobsScrapeTimestampMetric       prometheus.Gauge
		lastJobsScrapeDurationSecondsMetric prometheus.Gauge

		deploymentName                = "fake-deployment-name"
		jobName                       = "fake-job-name"
		jobID                         = "fake-job-id"
		jobIndex                      = "0"
		jobIP                         = "1.2.3.4"
		jobAZ                         = "fake-job-az"
		jobHealthy                    = true
		jobCPUSys                     = float64(0.5)
		jobCPUUser                    = float64(1.0)
		jobCPUWait                    = float64(1.5)
		jobMemKB                      = 1000
		jobMemPercent                 = 10
		jobSwapKB                     = 2000
		jobSwapPercent                = 20
		jobLoadAvg01                  = float64(0.01)
		jobLoadAvg05                  = float64(0.05)
		jobLoadAvg15                  = float64(0.15)
		jobSystemDiskInodePercent     = 10
		jobSystemDiskPercent          = 20
		jobEphemeralDiskInodePercent  = 30
		jobEphemeralDiskPercent       = 40
		jobPersistentDiskInodePercent = 50
		jobPersistentDiskPercent      = 60
		jobProcessName                = "fake-process-name"
		jobProcessUptime              = uint64(3600)
		jobProcessHealthy             = true
		jobProcessCPUTotal            = float64(0.5)
		jobProcessMemKB               = uint64(2000)
		jobProcessMemPercent          = float64(20)
	)

	BeforeEach(func() {
		namespace = testNamespace
		environment = testEnvironment
		boshName = testBoshName
		boshUUID = testBoshUUID
		azsFilter = filters.NewAZsFilter([]string{})
		cidrsFilter, err = filters.NewCidrFilter([]string{"0.0.0.0/0"})
		Expect(err).ToNot(HaveOccurred())

		jobHealthyMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "healthy",
				Help:      "BOSH Job Healthy (1 for healthy, 0 for unhealthy).",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobHealthyMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(float64(1))

		jobLoadAvg01Metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "load_avg01",
				Help:      "BOSH Job Load avg01.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobLoadAvg01Metric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(jobLoadAvg01)

		jobLoadAvg05Metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "load_avg05",
				Help:      "BOSH Job Load avg05.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobLoadAvg05Metric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(jobLoadAvg05)

		jobLoadAvg15Metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "load_avg15",
				Help:      "BOSH Job Load avg15.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobLoadAvg15Metric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(jobLoadAvg15)

		jobCPUSysMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "cpu_sys",
				Help:      "BOSH Job CPU System.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobCPUSysMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(jobCPUSys)

		jobCPUUserMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "cpu_user",
				Help:      "BOSH Job CPU User.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobCPUUserMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(jobCPUUser)

		jobCPUWaitMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "cpu_wait",
				Help:      "BOSH Job CPU Wait.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobCPUWaitMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(jobCPUWait)

		jobMemKBMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "mem_kb",
				Help:      "BOSH Job Memory KB.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobMemKBMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(float64(jobMemKB))

		jobMemPercentMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "mem_percent",
				Help:      "BOSH Job Memory Percent.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobMemPercentMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(float64(jobMemPercent))

		jobSwapKBMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "swap_kb",
				Help:      "BOSH Job Swap KB.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobSwapKBMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(float64(jobSwapKB))

		jobSwapPercentMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "swap_percent",
				Help:      "BOSH Job Swap Percent.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobSwapPercentMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(float64(jobSwapPercent))

		jobSystemDiskInodePercentMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "system_disk_inode_percent",
				Help:      "BOSH Job System Disk Inode Percent.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobSystemDiskInodePercentMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(float64(jobSystemDiskInodePercent))

		jobSystemDiskPercentMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "system_disk_percent",
				Help:      "BOSH Job System Disk Percent.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobSystemDiskPercentMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(float64(jobSystemDiskPercent))

		jobEphemeralDiskInodePercentMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "ephemeral_disk_inode_percent",
				Help:      "BOSH Job Ephemeral Disk Inode Percent.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobEphemeralDiskInodePercentMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(float64(jobEphemeralDiskInodePercent))

		jobEphemeralDiskPercentMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "ephemeral_disk_percent",
				Help:      "BOSH Job Ephemeral Disk Percent.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobEphemeralDiskPercentMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(float64(jobEphemeralDiskPercent))

		jobPersistentDiskInodePercentMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "persistent_disk_inode_percent",
				Help:      "BOSH Job Persistent Disk Inode Percent.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobPersistentDiskInodePercentMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(float64(jobPersistentDiskInodePercent))

		jobPersistentDiskPercentMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job",
				Name:      "persistent_disk_percent",
				Help:      "BOSH Job Persistent Disk Percent.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		)

		jobPersistentDiskPercentMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
		).Set(float64(jobPersistentDiskPercent))

		jobProcessHealthyMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job_process",
				Name:      "healthy",
				Help:      "BOSH Job Process Healthy (1 for healthy, 0 for unhealthy).",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
		)

		jobProcessHealthyMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		).Set(float64(1))

		jobProcessUptimeMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job_process",
				Name:      "uptime_seconds",
				Help:      "BOSH Job Process Uptime in seconds.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
		)

		jobProcessUptimeMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		).Set(float64(jobProcessUptime))

		jobProcessCPUTotalMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job_process",
				Name:      "cpu_total",
				Help:      "BOSH Job Process CPU Total.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
		)

		jobProcessCPUTotalMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		).Set(jobProcessCPUTotal)

		jobProcessMemKBMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job_process",
				Name:      "mem_kb",
				Help:      "BOSH Job Process Memory KB.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
		)

		jobProcessMemKBMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		).Set(float64(jobProcessMemKB))

		jobProcessMemPercentMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "job_process",
				Name:      "mem_percent",
				Help:      "BOSH Job Process Memory Percent.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
			[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
		)

		jobProcessMemPercentMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		).Set(jobProcessMemPercent)

		lastJobsScrapeTimestampMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_jobs_scrape_timestamp",
				Help:      "Number of seconds since 1970 since last scrape of Job metrics from BOSH.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
		)

		lastJobsScrapeDurationSecondsMetric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "last_jobs_scrape_duration_seconds",
				Help:      "Duration of the last scrape of Job metrics from BOSH.",
				ConstLabels: prometheus.Labels{
					"environment": environment,
					"bosh_name":   boshName,
					"bosh_uuid":   boshUUID,
				},
			},
		)
	})

	JustBeforeEach(func() {
		jobsCollector = NewJobsCollector(namespace, environment, boshName, boshUUID, azsFilter, cidrsFilter)
	})

	Describe("Describe", func() {
		var (
			descriptions chan *prometheus.Desc
		)

		BeforeEach(func() {
			descriptions = make(chan *prometheus.Desc)
		})

		JustBeforeEach(func() {
			go jobsCollector.Describe(descriptions)
		})

		It("returns a job_healthy metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobHealthyMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_load_avg01 metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobLoadAvg01Metric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_load_avg05 metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobLoadAvg05Metric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_load_avg15 metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobLoadAvg15Metric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_cpu_sys metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobCPUSysMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_cpu_user metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobCPUUserMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_cpu_wait metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobCPUWaitMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_mem_kb metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobMemKBMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_mem_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobMemPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_swap_kb metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobSwapKBMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_swap_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobSwapPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_system_disk_inode_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobSystemDiskInodePercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_system_disk_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobSystemDiskPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_ephemeral_disk_inode_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobEphemeralDiskInodePercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_ephemeral_disk_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobEphemeralDiskPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_persistent_disk_inode_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobPersistentDiskInodePercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_persistent_disk_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobPersistentDiskPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Desc())))
		})

		It("returns a job_process_healthy metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobProcessHealthyMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
				jobProcessName,
			).Desc())))
		})

		It("returns a job_process_uptime_seconds metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobProcessUptimeMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
				jobProcessName,
			).Desc())))
		})

		It("returns a job_process_cpu_total metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobProcessCPUTotalMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
				jobProcessName,
			).Desc())))
		})

		It("returns a job_process_mem_kb metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobProcessMemKBMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
				jobProcessName,
			).Desc())))
		})

		It("returns a job_process_mem_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(jobProcessMemPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
				jobProcessName,
			).Desc())))
		})

		It("returns a last_jobs_scrape_timestamp metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastJobsScrapeTimestampMetric.Desc())))
		})

		It("returns a last_jobs_scrape_duration_seconds metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(lastJobsScrapeDurationSecondsMetric.Desc())))
		})
	})

	Describe("Collect", func() {
		var (
			processes       []deployments.Process
			vitals          deployments.Vitals
			instances       []deployments.Instance
			deploymentInfo  deployments.DeploymentInfo
			deploymentsInfo []deployments.DeploymentInfo

			metrics    chan prometheus.Metric
			errMetrics chan error
		)

		BeforeEach(func() {
			processes = []deployments.Process{
				{
					Name:    jobProcessName,
					Uptime:  &jobProcessUptime,
					Healthy: jobProcessHealthy,
					CPU:     deployments.CPU{Total: &jobProcessCPUTotal},
					Mem:     deployments.MemInt{KB: &jobProcessMemKB, Percent: &jobProcessMemPercent},
				},
			}

			vitals = deployments.Vitals{
				CPU: deployments.CPU{
					Sys:  strconv.FormatFloat(jobCPUSys, 'E', -1, 64),
					User: strconv.FormatFloat(jobCPUUser, 'E', -1, 64),
					Wait: strconv.FormatFloat(jobCPUWait, 'E', -1, 64),
				},
				Mem: deployments.Mem{
					KB:      strconv.Itoa(jobMemKB),
					Percent: strconv.Itoa(jobMemPercent),
				},
				Swap: deployments.Mem{
					KB:      strconv.Itoa(jobSwapKB),
					Percent: strconv.Itoa(jobSwapPercent),
				},
				Load: []string{
					strconv.FormatFloat(jobLoadAvg01, 'E', -1, 64),
					strconv.FormatFloat(jobLoadAvg05, 'E', -1, 64),
					strconv.FormatFloat(jobLoadAvg15, 'E', -1, 64),
				},
				SystemDisk: deployments.Disk{
					InodePercent: strconv.Itoa(jobSystemDiskInodePercent),
					Percent:      strconv.Itoa(jobSystemDiskPercent),
				},
				EphemeralDisk: deployments.Disk{
					InodePercent: strconv.Itoa(jobEphemeralDiskInodePercent),
					Percent:      strconv.Itoa(jobEphemeralDiskPercent),
				},
				PersistentDisk: deployments.Disk{
					InodePercent: strconv.Itoa(jobPersistentDiskInodePercent),
					Percent:      strconv.Itoa(jobPersistentDiskPercent),
				},
			}

			instances = []deployments.Instance{
				{
					Name:      jobName,
					ID:        jobID,
					Index:     jobIndex,
					IPs:       []string{jobIP},
					AZ:        jobAZ,
					Healthy:   jobHealthy,
					Vitals:    vitals,
					Processes: processes,
				},
			}

			deploymentInfo = deployments.DeploymentInfo{
				Name:      deploymentName,
				Instances: instances,
			}

			deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}

			metrics = make(chan prometheus.Metric)
			errMetrics = make(chan error, 1)
		})

		JustBeforeEach(func() {
			go func() {
				if err := jobsCollector.Collect(deploymentsInfo, metrics); err != nil {
					errMetrics <- err
				}
			}()
		})

		It("returns a job_process_healthy metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobHealthyMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when the process is not running", func() {
			BeforeEach(func() {
				instances[0].Healthy = false

				jobHealthyMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				).Set(float64(0))
			})

			It("returns a job_process_healthy metric", func() {
				Eventually(metrics).Should(Receive(PrometheusMetric(jobHealthyMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_load_avg01 metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobLoadAvg01Metric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		It("returns a job_load_avg05 metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobLoadAvg05Metric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		It("returns a job_load_avg15 metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobLoadAvg15Metric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no load avg values", func() {
			BeforeEach(func() {
				instances[0].Vitals.Load = []string{}
			})

			It("does not return any job_load_avg metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobLoadAvg01Metric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobLoadAvg05Metric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobLoadAvg15Metric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_cpu_sys metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobCPUSysMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no cpu sys value", func() {
			BeforeEach(func() {
				instances[0].Vitals.CPU = deployments.CPU{
					User: strconv.FormatFloat(jobCPUUser, 'E', -1, 64),
					Wait: strconv.FormatFloat(jobCPUWait, 'E', -1, 64),
				}
			})

			It("does not return a job_cpu_sys metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobCPUSysMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_cpu_user metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobCPUUserMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no cpu user value", func() {
			BeforeEach(func() {
				instances[0].Vitals.CPU = deployments.CPU{
					Sys:  strconv.FormatFloat(jobCPUSys, 'E', -1, 64),
					Wait: strconv.FormatFloat(jobCPUWait, 'E', -1, 64),
				}
			})

			It("does not return a job_cpu_user metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobCPUUserMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_cpu_wait metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobCPUWaitMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no cpu wait value", func() {
			BeforeEach(func() {
				instances[0].Vitals.CPU = deployments.CPU{
					Sys:  strconv.FormatFloat(jobCPUSys, 'E', -1, 64),
					User: strconv.FormatFloat(jobCPUUser, 'E', -1, 64),
				}
			})

			It("does not return a job_cpu_wait metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobCPUWaitMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_mem_kb metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobMemKBMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no mem kb value", func() {
			BeforeEach(func() {
				instances[0].Vitals.Mem = deployments.Mem{
					Percent: strconv.Itoa(jobMemPercent),
				}
			})

			It("does not return a job_mem_kb metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobMemKBMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_mem_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobMemPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no mem percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.Mem = deployments.Mem{
					KB: strconv.Itoa(jobMemKB),
				}
			})

			It("does not return a job_mem_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobMemPercentMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_swap_kb metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobSwapKBMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no swap kb value", func() {
			BeforeEach(func() {
				instances[0].Vitals.Swap = deployments.Mem{
					Percent: strconv.Itoa(jobSwapPercent),
				}
			})

			It("does not return a job_swap_kb metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobSwapKBMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_swap_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobSwapPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no swap percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.Swap = deployments.Mem{
					KB: strconv.Itoa(jobSwapKB),
				}
			})

			It("does not return a job_swap_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobSwapPercentMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_system_disk_inode_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobSystemDiskInodePercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no system disk inode percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.SystemDisk = deployments.Disk{
					Percent: strconv.Itoa(jobSystemDiskPercent),
				}
			})

			It("does not return a job_system_disk_inode_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobSystemDiskInodePercentMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_system_disk_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobSystemDiskPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no system disk percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.SystemDisk = deployments.Disk{
					InodePercent: strconv.Itoa(jobSystemDiskInodePercent),
				}
			})

			It("does not return a job_system_disk_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobSystemDiskPercentMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_ephemeral_disk_inode_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobEphemeralDiskInodePercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no ephemeral disk inode percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.EphemeralDisk = deployments.Disk{
					Percent: strconv.Itoa(jobEphemeralDiskPercent),
				}
			})

			It("does not return a job_ephemeral_disk_inode_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobEphemeralDiskInodePercentMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_ephemeral_disk_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobEphemeralDiskPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no ephemeral disk percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.EphemeralDisk = deployments.Disk{
					InodePercent: strconv.Itoa(jobEphemeralDiskInodePercent),
				}
			})

			It("does not return a job_Ephemeral_disk_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobEphemeralDiskPercentMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_persistent_disk_inode_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobPersistentDiskInodePercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no persistent disk inode percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.PersistentDisk = deployments.Disk{
					Percent: strconv.Itoa(jobPersistentDiskPercent),
				}
			})

			It("does not return a job_persistent_disk_inode_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobPersistentDiskInodePercentMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_persistent_disk_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobPersistentDiskPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no persistent disk percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.PersistentDisk = deployments.Disk{
					InodePercent: strconv.Itoa(jobPersistentDiskInodePercent),
				}
			})

			It("does not return a job_persistent_disk_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobPersistentDiskPercentMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a healthy job_process_healthy metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobProcessHealthyMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
				jobProcessName,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when a process is not running", func() {
			BeforeEach(func() {
				instances[0].Processes[0].Healthy = false

				jobProcessHealthyMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
					jobProcessName,
				).Set(float64(0))
			})

			It("returns an unhealthy job_process_healthy metric", func() {
				Eventually(metrics).Should(Receive(PrometheusMetric(jobProcessHealthyMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
					jobProcessName,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_process_uptime_seconds metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobProcessUptimeMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
				jobProcessName,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no process uptime value", func() {
			BeforeEach(func() {
				instances[0].Processes[0].Uptime = nil
			})

			It("does not return a job_process_uptime_seconds metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobProcessUptimeMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
					jobProcessName,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_process_cpu_total metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobProcessCPUTotalMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
				jobProcessName,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no process cpu total value", func() {
			BeforeEach(func() {
				instances[0].Processes[0].CPU = deployments.CPU{}
			})

			It("does not return a job_process_cpu_total metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobProcessCPUTotalMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
					jobProcessName,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_process_mem_kb metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobProcessMemKBMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
				jobProcessName,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no process mem kb value", func() {
			BeforeEach(func() {
				instances[0].Processes[0].Mem = deployments.MemInt{Percent: &jobProcessMemPercent}
			})

			It("does not return a job_process_mem_kb metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobProcessMemKBMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
					jobProcessName,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_process_mem_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(jobProcessMemPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
				jobProcessName,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no process mem percent value", func() {
			BeforeEach(func() {
				instances[0].Processes[0].Mem = deployments.MemInt{KB: &jobProcessMemKB}
			})

			It("does not return a job_process_mem_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(jobProcessMemPercentMetric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
					jobProcessName,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		Context("when there are no deployments", func() {
			BeforeEach(func() {
				deploymentsInfo = []deployments.DeploymentInfo{}
			})

			It("returns only a last_jobs_scrape_timestamp & last_jobs_scrape_duration_seconds metric", func() {
				Eventually(metrics).Should(Receive())
				Eventually(metrics).Should(Receive())
				Consistently(metrics).ShouldNot(Receive())
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		Context("when there are no instances", func() {
			BeforeEach(func() {
				deploymentInfo.Instances = []deployments.Instance{}
				deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}
			})

			It("returns only a last_jobs_scrape_timestamp & last_jobs_scrape_duration_seconds metric", func() {
				Eventually(metrics).Should(Receive())
				Eventually(metrics).Should(Receive())
				Consistently(metrics).ShouldNot(Receive())
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})
	})
})
