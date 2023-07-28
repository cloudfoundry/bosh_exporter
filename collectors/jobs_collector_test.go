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

type BaseLabelValues struct {
	deploymentName string
	jobName        string
	jobID          string
	jobIndex       string
	jobAZ          string
	jobIP          string
}

func (b *BaseLabelValues) AddLabelValues(gaugeVec *prometheus.GaugeVec, lvs ...string) prometheus.Gauge {
	values := []string{b.deploymentName, b.jobName, b.jobID, b.jobIndex, b.jobAZ, b.jobIP}
	values = append(values, lvs...)
	return gaugeVec.WithLabelValues(values...)
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
		metrics       *JobsCollectorMetrics
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
		jobProcessInfoMetric                *prometheus.GaugeVec
		jobProcessHealthyMetric             *prometheus.GaugeVec
		jobProcessUptimeMetric              *prometheus.GaugeVec
		jobProcessCPUTotalMetric            *prometheus.GaugeVec
		jobProcessMemKBMetric               *prometheus.GaugeVec
		jobProcessMemPercentMetric          *prometheus.GaugeVec
		lastJobsScrapeTimestampMetric       prometheus.Gauge
		lastJobsScrapeDurationSecondsMetric prometheus.Gauge

		baseLabelValues = BaseLabelValues{
			deploymentName: "fake-deployment-name",
			jobName:        "fake-job-name",
			jobID:          "fake-job-id",
			jobIndex:       "0",
			jobIP:          "1.2.3.4",
			jobAZ:          "fake-job-az",
		}
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
		jobProcessReleaseName         = "fake-process-release-name"
		jobProcessReleaseVersion      = "fake-process-release-version"
	)

	BeforeEach(func() {
		namespace = testNamespace
		environment = testEnvironment
		boshName = testBoshName
		boshUUID = testBoshUUID
		metrics = NewJobsCollectorMetrics(testNamespace, testEnvironment, testBoshName, testBoshUUID)
		azsFilter = filters.NewAZsFilter([]string{})
		cidrsFilter, err = filters.NewCidrFilter([]string{"0.0.0.0/0"})
		Expect(err).ToNot(HaveOccurred())

		jobHealthyMetric = metrics.NewJobHealthyMetric()
		baseLabelValues.AddLabelValues(jobHealthyMetric).Set(float64(1))

		jobLoadAvg01Metric = metrics.NewJobLoadAvg01Metric()
		baseLabelValues.AddLabelValues(jobLoadAvg01Metric).Set(jobLoadAvg01)

		jobLoadAvg05Metric = metrics.NewJobLoadAvg05Metric()
		baseLabelValues.AddLabelValues(jobLoadAvg05Metric).Set(jobLoadAvg05)

		jobLoadAvg15Metric = metrics.NewJobLoadAvg15Metric()
		baseLabelValues.AddLabelValues(jobLoadAvg15Metric).Set(jobLoadAvg15)

		jobCPUSysMetric = metrics.NewJobCPUSysMetric()
		baseLabelValues.AddLabelValues(jobCPUSysMetric).Set(jobCPUSys)

		jobCPUUserMetric = metrics.NewJobCPUUserMetric()
		baseLabelValues.AddLabelValues(jobCPUUserMetric).Set(jobCPUUser)

		jobCPUWaitMetric = metrics.NewJobCPUWaitMetric()
		baseLabelValues.AddLabelValues(jobCPUWaitMetric).Set(jobCPUWait)

		jobMemKBMetric = metrics.NewJobMemKBMetric()
		baseLabelValues.AddLabelValues(jobMemKBMetric).Set(float64(jobMemKB))

		jobMemPercentMetric = metrics.NewJobMemPercentMetric()
		baseLabelValues.AddLabelValues(jobMemPercentMetric).Set(float64(jobMemPercent))

		jobSwapKBMetric = metrics.NewJobSwapKBMetric()
		baseLabelValues.AddLabelValues(jobSwapKBMetric).Set(float64(jobSwapKB))

		jobSwapPercentMetric = metrics.NewJobSwapPercentMetric()
		baseLabelValues.AddLabelValues(jobSwapPercentMetric).Set(float64(jobSwapPercent))

		jobSystemDiskInodePercentMetric = metrics.NewJobSystemDiskInodePercentMetric()
		baseLabelValues.AddLabelValues(jobSystemDiskInodePercentMetric).Set(float64(jobSystemDiskInodePercent))

		jobSystemDiskPercentMetric = metrics.NewJobSystemDiskPercentMetric()
		baseLabelValues.AddLabelValues(jobSystemDiskPercentMetric).Set(float64(jobSystemDiskPercent))

		jobEphemeralDiskInodePercentMetric = metrics.NewJobEphemeralDiskInodePercentMetric()
		baseLabelValues.AddLabelValues(jobEphemeralDiskInodePercentMetric).Set(float64(jobEphemeralDiskInodePercent))

		jobEphemeralDiskPercentMetric = metrics.NewJobEphemeralDiskPercentMetric()
		baseLabelValues.AddLabelValues(jobEphemeralDiskPercentMetric).Set(float64(jobEphemeralDiskPercent))

		jobPersistentDiskInodePercentMetric = metrics.NewJobPersistentDiskInodePercentMetric()
		baseLabelValues.AddLabelValues(jobPersistentDiskInodePercentMetric).Set(float64(jobPersistentDiskInodePercent))

		jobPersistentDiskPercentMetric = metrics.NewJobPersistentDiskPercentMetric()
		baseLabelValues.AddLabelValues(jobPersistentDiskPercentMetric).Set(float64(jobPersistentDiskPercent))

		jobProcessInfoMetric = metrics.NewJobProcessInfoMetric()
		baseLabelValues.AddLabelValues(jobProcessInfoMetric, jobProcessName, jobProcessReleaseName, jobProcessReleaseVersion).Set(float64(1))

		jobProcessHealthyMetric = metrics.NewJobProcessHealthyMetric()
		baseLabelValues.AddLabelValues(jobProcessHealthyMetric, jobProcessName).Set(float64(1))

		jobProcessUptimeMetric = metrics.NewJobProcessUptimeMetric()
		baseLabelValues.AddLabelValues(jobProcessUptimeMetric, jobProcessName).Set(float64(jobProcessUptime))

		jobProcessCPUTotalMetric = metrics.NewJobProcessCPUTotalMetric()
		baseLabelValues.AddLabelValues(jobProcessCPUTotalMetric, jobProcessName).Set(jobProcessCPUTotal)

		jobProcessMemKBMetric = metrics.NewJobProcessMemKBMetric()
		baseLabelValues.AddLabelValues(jobProcessMemKBMetric, jobProcessName).Set(float64(jobProcessMemKB))

		jobProcessMemPercentMetric = metrics.NewJobProcessMemPercentMetric()
		baseLabelValues.AddLabelValues(jobProcessMemPercentMetric, jobProcessName).Set(jobProcessMemPercent)

		lastJobsScrapeTimestampMetric = metrics.NewLastJobsScrapeTimestampMetric()
		lastJobsScrapeDurationSecondsMetric = metrics.NewLastJobsScrapeDurationSecondsMetric()
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
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobHealthyMetric).Desc())))
		})

		It("returns a job_load_avg01 metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobLoadAvg01Metric).Desc())))
		})

		It("returns a job_load_avg05 metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobLoadAvg05Metric).Desc())))
		})

		It("returns a job_load_avg15 metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobLoadAvg15Metric).Desc())))
		})

		It("returns a job_cpu_sys metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobCPUSysMetric).Desc())))
		})

		It("returns a job_cpu_user metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobCPUUserMetric).Desc())))
		})

		It("returns a job_cpu_wait metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobCPUWaitMetric).Desc())))
		})

		It("returns a job_mem_kb metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobMemKBMetric).Desc())))
		})

		It("returns a job_mem_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobMemPercentMetric).Desc())))
		})

		It("returns a job_swap_kb metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobSwapKBMetric).Desc())))
		})

		It("returns a job_swap_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobSwapPercentMetric).Desc())))
		})

		It("returns a job_system_disk_inode_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobSystemDiskInodePercentMetric).Desc())))
		})

		It("returns a job_system_disk_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobSystemDiskPercentMetric).Desc())))
		})

		It("returns a job_ephemeral_disk_inode_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobEphemeralDiskInodePercentMetric).Desc())))
		})

		It("returns a job_ephemeral_disk_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobEphemeralDiskPercentMetric).Desc())))
		})

		It("returns a job_persistent_disk_inode_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobPersistentDiskInodePercentMetric).Desc())))
		})

		It("returns a job_persistent_disk_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobPersistentDiskPercentMetric).Desc())))
		})

		It("returns a job_process_info metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobProcessInfoMetric, jobProcessName, jobProcessReleaseName, jobProcessReleaseVersion).Desc())))
		})

		It("returns a job_process_healthy metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobProcessHealthyMetric, jobProcessName).Desc())))
		})

		It("returns a job_process_uptime_seconds metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobProcessUptimeMetric, jobProcessName).Desc())))
		})

		It("returns a job_process_cpu_total metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobProcessCPUTotalMetric, jobProcessName).Desc())))
		})

		It("returns a job_process_mem_kb metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobProcessMemKBMetric, jobProcessName).Desc())))
		})

		It("returns a job_process_mem_percent metric description", func() {
			Eventually(descriptions).Should(Receive(Equal(baseLabelValues.AddLabelValues(jobProcessMemPercentMetric, jobProcessName).Desc())))
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
					Name:      baseLabelValues.jobName,
					ID:        baseLabelValues.jobID,
					Index:     baseLabelValues.jobIndex,
					IPs:       []string{baseLabelValues.jobIP},
					AZ:        baseLabelValues.jobAZ,
					Healthy:   jobHealthy,
					Vitals:    vitals,
					Processes: processes,
				},
			}

			deploymentInfo = deployments.DeploymentInfo{
				Name:      baseLabelValues.deploymentName,
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
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobHealthyMetric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when the process is not running", func() {
			BeforeEach(func() {
				instances[0].Healthy = false
				baseLabelValues.AddLabelValues(jobHealthyMetric).Set(float64(0))
			})

			It("returns a job_process_healthy metric", func() {
				Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobHealthyMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_load_avg01 metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg01Metric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		It("returns a job_load_avg05 metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg05Metric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		It("returns a job_load_avg15 metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg15Metric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no load avg values", func() {
			BeforeEach(func() {
				instances[0].Vitals.Load = []string{}
			})

			It("does not return any job_load_avg metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg01Metric))))
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg05Metric))))
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg15Metric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_cpu_sys metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUSysMetric))))
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
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUSysMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_cpu_user metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUUserMetric))))
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
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUUserMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_cpu_wait metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUWaitMetric))))
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
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUWaitMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_mem_kb metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobMemKBMetric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no mem kb value", func() {
			BeforeEach(func() {
				instances[0].Vitals.Mem = deployments.Mem{
					Percent: strconv.Itoa(jobMemPercent),
				}
			})

			It("does not return a job_mem_kb metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobMemKBMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_mem_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobMemPercentMetric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no mem percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.Mem = deployments.Mem{
					KB: strconv.Itoa(jobMemKB),
				}
			})

			It("does not return a job_mem_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobMemPercentMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_swap_kb metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobSwapKBMetric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no swap kb value", func() {
			BeforeEach(func() {
				instances[0].Vitals.Swap = deployments.Mem{
					Percent: strconv.Itoa(jobSwapPercent),
				}
			})

			It("does not return a job_swap_kb metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobSwapKBMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_swap_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobSwapPercentMetric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no swap percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.Swap = deployments.Mem{
					KB: strconv.Itoa(jobSwapKB),
				}
			})

			It("does not return a job_swap_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobSwapPercentMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_system_disk_inode_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobSystemDiskInodePercentMetric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no system disk inode percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.SystemDisk = deployments.Disk{
					Percent: strconv.Itoa(jobSystemDiskPercent),
				}
			})

			It("does not return a job_system_disk_inode_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobSystemDiskInodePercentMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_system_disk_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobSystemDiskPercentMetric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no system disk percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.SystemDisk = deployments.Disk{
					InodePercent: strconv.Itoa(jobSystemDiskInodePercent),
				}
			})

			It("does not return a job_system_disk_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobSystemDiskPercentMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_ephemeral_disk_inode_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobEphemeralDiskInodePercentMetric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no ephemeral disk inode percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.EphemeralDisk = deployments.Disk{
					Percent: strconv.Itoa(jobEphemeralDiskPercent),
				}
			})

			It("does not return a job_ephemeral_disk_inode_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobEphemeralDiskInodePercentMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_ephemeral_disk_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobEphemeralDiskPercentMetric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no ephemeral disk percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.EphemeralDisk = deployments.Disk{
					InodePercent: strconv.Itoa(jobEphemeralDiskInodePercent),
				}
			})

			It("does not return a job_Ephemeral_disk_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobEphemeralDiskPercentMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_persistent_disk_inode_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobPersistentDiskInodePercentMetric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no persistent disk inode percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.PersistentDisk = deployments.Disk{
					Percent: strconv.Itoa(jobPersistentDiskPercent),
				}
			})

			It("does not return a job_persistent_disk_inode_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobPersistentDiskInodePercentMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_persistent_disk_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobPersistentDiskPercentMetric))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no persistent disk percent value", func() {
			BeforeEach(func() {
				instances[0].Vitals.PersistentDisk = deployments.Disk{
					InodePercent: strconv.Itoa(jobPersistentDiskInodePercent),
				}
			})

			It("does not return a job_persistent_disk_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobPersistentDiskPercentMetric))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a healthy job_process_healthy metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessHealthyMetric, jobProcessName))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when a process is not running", func() {
			BeforeEach(func() {
				instances[0].Processes[0].Healthy = false
				baseLabelValues.AddLabelValues(jobProcessHealthyMetric, jobProcessName).Set(float64(0))
			})

			It("returns an unhealthy job_process_healthy metric", func() {
				Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessHealthyMetric, jobProcessName))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_process_uptime_seconds metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessUptimeMetric, jobProcessName))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no process uptime value", func() {
			BeforeEach(func() {
				instances[0].Processes[0].Uptime = nil
			})

			It("does not return a job_process_uptime_seconds metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessUptimeMetric, jobProcessName))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_process_cpu_total metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessCPUTotalMetric,
				jobProcessName,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no process cpu total value", func() {
			BeforeEach(func() {
				instances[0].Processes[0].CPU = deployments.CPU{}
			})

			It("does not return a job_process_cpu_total metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessCPUTotalMetric,
					jobProcessName,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_process_mem_kb metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessMemKBMetric,
				jobProcessName,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no process mem kb value", func() {
			BeforeEach(func() {
				instances[0].Processes[0].Mem = deployments.MemInt{Percent: &jobProcessMemPercent}
			})

			It("does not return a job_process_mem_kb metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessMemKBMetric,
					jobProcessName,
				))))
				Consistently(errMetrics).ShouldNot(Receive())
			})
		})

		It("returns a job_process_mem_percent metric", func() {
			Eventually(metrics).Should(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessMemPercentMetric,
				jobProcessName,
			))))
			Consistently(errMetrics).ShouldNot(Receive())
		})

		Context("when there is no process mem percent value", func() {
			BeforeEach(func() {
				instances[0].Processes[0].Mem = deployments.MemInt{KB: &jobProcessMemKB}
			})

			It("does not return a job_process_mem_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessMemPercentMetric,
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
