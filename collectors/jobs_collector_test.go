package collectors_test

import (
	"strconv"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/cloudfoundry/bosh_exporter/deployments"
	"github.com/cloudfoundry/bosh_exporter/filters"

	"github.com/cloudfoundry/bosh_exporter/collectors"
	"github.com/cloudfoundry/bosh_exporter/utils/matchers"
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

var _ = ginkgo.Describe("JobsCollector", func() {
	var (
		err           error
		namespace     string
		environment   string
		boshName      string
		boshUUID      string
		azsFilter     *filters.AZsFilter
		cidrsFilter   *filters.CidrFilter
		metrics       *collectors.JobsCollectorMetrics
		jobsCollector *collectors.JobsCollector

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

	ginkgo.BeforeEach(func() {
		namespace = testNamespace
		environment = testEnvironment
		boshName = testBoshName
		boshUUID = testBoshUUID
		metrics = collectors.NewJobsCollectorMetrics(testNamespace, testEnvironment, testBoshName, testBoshUUID)
		azsFilter = filters.NewAZsFilter([]string{})
		cidrsFilter, err = filters.NewCidrFilter([]string{"0.0.0.0/0"})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

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

	ginkgo.JustBeforeEach(func() {
		jobsCollector = collectors.NewJobsCollector(namespace, environment, boshName, boshUUID, azsFilter, cidrsFilter)
	})

	ginkgo.Describe("ginkgo.Describe", func() {
		var (
			descriptions chan *prometheus.Desc
		)

		ginkgo.BeforeEach(func() {
			descriptions = make(chan *prometheus.Desc)
		})

		ginkgo.JustBeforeEach(func() {
			go jobsCollector.Describe(descriptions)
		})

		ginkgo.It("returns a job_healthy metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobHealthyMetric).Desc())))
		})

		ginkgo.It("returns a job_load_avg01 metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobLoadAvg01Metric).Desc())))
		})

		ginkgo.It("returns a job_load_avg05 metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobLoadAvg05Metric).Desc())))
		})

		ginkgo.It("returns a job_load_avg15 metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobLoadAvg15Metric).Desc())))
		})

		ginkgo.It("returns a job_cpu_sys metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobCPUSysMetric).Desc())))
		})

		ginkgo.It("returns a job_cpu_user metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobCPUUserMetric).Desc())))
		})

		ginkgo.It("returns a job_cpu_wait metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobCPUWaitMetric).Desc())))
		})

		ginkgo.It("returns a job_mem_kb metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobMemKBMetric).Desc())))
		})

		ginkgo.It("returns a job_mem_percent metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobMemPercentMetric).Desc())))
		})

		ginkgo.It("returns a job_swap_kb metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobSwapKBMetric).Desc())))
		})

		ginkgo.It("returns a job_swap_percent metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobSwapPercentMetric).Desc())))
		})

		ginkgo.It("returns a job_system_disk_inode_percent metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobSystemDiskInodePercentMetric).Desc())))
		})

		ginkgo.It("returns a job_system_disk_percent metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobSystemDiskPercentMetric).Desc())))
		})

		ginkgo.It("returns a job_ephemeral_disk_inode_percent metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobEphemeralDiskInodePercentMetric).Desc())))
		})

		ginkgo.It("returns a job_ephemeral_disk_percent metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobEphemeralDiskPercentMetric).Desc())))
		})

		ginkgo.It("returns a job_persistent_disk_inode_percent metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobPersistentDiskInodePercentMetric).Desc())))
		})

		ginkgo.It("returns a job_persistent_disk_percent metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobPersistentDiskPercentMetric).Desc())))
		})

		ginkgo.It("returns a job_process_info metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobProcessInfoMetric, jobProcessName, jobProcessReleaseName, jobProcessReleaseVersion).Desc())))
		})

		ginkgo.It("returns a job_process_healthy metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobProcessHealthyMetric, jobProcessName).Desc())))
		})

		ginkgo.It("returns a job_process_uptime_seconds metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobProcessUptimeMetric, jobProcessName).Desc())))
		})

		ginkgo.It("returns a job_process_cpu_total metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobProcessCPUTotalMetric, jobProcessName).Desc())))
		})

		ginkgo.It("returns a job_process_mem_kb metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobProcessMemKBMetric, jobProcessName).Desc())))
		})

		ginkgo.It("returns a job_process_mem_percent metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(baseLabelValues.AddLabelValues(jobProcessMemPercentMetric, jobProcessName).Desc())))
		})

		ginkgo.It("returns a last_jobs_scrape_timestamp metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(lastJobsScrapeTimestampMetric.Desc())))
		})

		ginkgo.It("returns a last_jobs_scrape_duration_seconds metric description", func() {
			gomega.Eventually(descriptions).Should(gomega.Receive(gomega.Equal(lastJobsScrapeDurationSecondsMetric.Desc())))
		})
	})

	ginkgo.Describe("Collect", func() {
		var (
			processes       []deployments.Process
			vitals          deployments.Vitals
			instances       []deployments.Instance
			deploymentInfo  deployments.DeploymentInfo
			deploymentsInfo []deployments.DeploymentInfo

			metrics    chan prometheus.Metric
			errMetrics chan error
		)

		ginkgo.BeforeEach(func() {
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

		ginkgo.JustBeforeEach(func() {
			go func() {
				if err := jobsCollector.Collect(deploymentsInfo, metrics); err != nil {
					errMetrics <- err
				}
			}()
		})

		ginkgo.It("returns a job_process_healthy metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobHealthyMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when the process is not running", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Healthy = false
				baseLabelValues.AddLabelValues(jobHealthyMetric).Set(float64(0))
			})

			ginkgo.It("returns a job_process_healthy metric", func() {
				gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobHealthyMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_load_avg01 metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg01Metric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.It("returns a job_load_avg05 metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg05Metric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.It("returns a job_load_avg15 metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg15Metric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no load avg values", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.Load = []string{}
			})

			ginkgo.It("does not return any job_load_avg metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg01Metric))))
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg05Metric))))
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobLoadAvg15Metric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_cpu_sys metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUSysMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no cpu sys value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.CPU = deployments.CPU{
					User: strconv.FormatFloat(jobCPUUser, 'E', -1, 64),
					Wait: strconv.FormatFloat(jobCPUWait, 'E', -1, 64),
				}
			})

			ginkgo.It("does not return a job_cpu_sys metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUSysMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_cpu_user metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUUserMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no cpu user value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.CPU = deployments.CPU{
					Sys:  strconv.FormatFloat(jobCPUSys, 'E', -1, 64),
					Wait: strconv.FormatFloat(jobCPUWait, 'E', -1, 64),
				}
			})

			ginkgo.It("does not return a job_cpu_user metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUUserMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_cpu_wait metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUWaitMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no cpu wait value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.CPU = deployments.CPU{
					Sys:  strconv.FormatFloat(jobCPUSys, 'E', -1, 64),
					User: strconv.FormatFloat(jobCPUUser, 'E', -1, 64),
				}
			})

			ginkgo.It("does not return a job_cpu_wait metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobCPUWaitMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_mem_kb metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobMemKBMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no mem kb value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.Mem = deployments.Mem{
					Percent: strconv.Itoa(jobMemPercent),
				}
			})

			ginkgo.It("does not return a job_mem_kb metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobMemKBMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_mem_percent metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobMemPercentMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no mem percent value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.Mem = deployments.Mem{
					KB: strconv.Itoa(jobMemKB),
				}
			})

			ginkgo.It("does not return a job_mem_percent metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobMemPercentMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_swap_kb metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobSwapKBMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no swap kb value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.Swap = deployments.Mem{
					Percent: strconv.Itoa(jobSwapPercent),
				}
			})

			ginkgo.It("does not return a job_swap_kb metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobSwapKBMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_swap_percent metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobSwapPercentMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no swap percent value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.Swap = deployments.Mem{
					KB: strconv.Itoa(jobSwapKB),
				}
			})

			ginkgo.It("does not return a job_swap_percent metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobSwapPercentMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_system_disk_inode_percent metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobSystemDiskInodePercentMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no system disk inode percent value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.SystemDisk = deployments.Disk{
					Percent: strconv.Itoa(jobSystemDiskPercent),
				}
			})

			ginkgo.It("does not return a job_system_disk_inode_percent metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobSystemDiskInodePercentMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_system_disk_percent metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobSystemDiskPercentMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no system disk percent value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.SystemDisk = deployments.Disk{
					InodePercent: strconv.Itoa(jobSystemDiskInodePercent),
				}
			})

			ginkgo.It("does not return a job_system_disk_percent metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobSystemDiskPercentMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_ephemeral_disk_inode_percent metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobEphemeralDiskInodePercentMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no ephemeral disk inode percent value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.EphemeralDisk = deployments.Disk{
					Percent: strconv.Itoa(jobEphemeralDiskPercent),
				}
			})

			ginkgo.It("does not return a job_ephemeral_disk_inode_percent metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobEphemeralDiskInodePercentMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_ephemeral_disk_percent metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobEphemeralDiskPercentMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no ephemeral disk percent value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.EphemeralDisk = deployments.Disk{
					InodePercent: strconv.Itoa(jobEphemeralDiskInodePercent),
				}
			})

			ginkgo.It("does not return a job_Ephemeral_disk_percent metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobEphemeralDiskPercentMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_persistent_disk_inode_percent metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobPersistentDiskInodePercentMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no persistent disk inode percent value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.PersistentDisk = deployments.Disk{
					Percent: strconv.Itoa(jobPersistentDiskPercent),
				}
			})

			ginkgo.It("does not return a job_persistent_disk_inode_percent metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobPersistentDiskInodePercentMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_persistent_disk_percent metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobPersistentDiskPercentMetric))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no persistent disk percent value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Vitals.PersistentDisk = deployments.Disk{
					InodePercent: strconv.Itoa(jobPersistentDiskInodePercent),
				}
			})

			ginkgo.It("does not return a job_persistent_disk_percent metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobPersistentDiskPercentMetric))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a healthy job_process_healthy metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessHealthyMetric, jobProcessName))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when a process is not running", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Processes[0].Healthy = false
				baseLabelValues.AddLabelValues(jobProcessHealthyMetric, jobProcessName).Set(float64(0))
			})

			ginkgo.It("returns an unhealthy job_process_healthy metric", func() {
				gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessHealthyMetric, jobProcessName))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_process_uptime_seconds metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessUptimeMetric, jobProcessName))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no process uptime value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Processes[0].Uptime = nil
			})

			ginkgo.It("does not return a job_process_uptime_seconds metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessUptimeMetric, jobProcessName))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_process_cpu_total metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessCPUTotalMetric,
				jobProcessName,
			))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no process cpu total value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Processes[0].CPU = deployments.CPU{}
			})

			ginkgo.It("does not return a job_process_cpu_total metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessCPUTotalMetric,
					jobProcessName,
				))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_process_mem_kb metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessMemKBMetric,
				jobProcessName,
			))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no process mem kb value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Processes[0].Mem = deployments.MemInt{Percent: &jobProcessMemPercent}
			})

			ginkgo.It("does not return a job_process_mem_kb metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessMemKBMetric,
					jobProcessName,
				))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.It("returns a job_process_mem_percent metric", func() {
			gomega.Eventually(metrics).Should(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessMemPercentMetric,
				jobProcessName,
			))))
			gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
		})

		ginkgo.Context("when there is no process mem percent value", func() {
			ginkgo.BeforeEach(func() {
				instances[0].Processes[0].Mem = deployments.MemInt{KB: &jobProcessMemKB}
			})

			ginkgo.It("does not return a job_process_mem_percent metric", func() {
				gomega.Consistently(metrics).ShouldNot(gomega.Receive(matchers.PrometheusMetric(baseLabelValues.AddLabelValues(jobProcessMemPercentMetric,
					jobProcessName,
				))))
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.Context("when there are no deployments", func() {
			ginkgo.BeforeEach(func() {
				deploymentsInfo = []deployments.DeploymentInfo{}
			})

			ginkgo.It("returns only a last_jobs_scrape_timestamp & last_jobs_scrape_duration_seconds metric", func() {
				gomega.Eventually(metrics).Should(gomega.Receive())
				gomega.Eventually(metrics).Should(gomega.Receive())
				gomega.Consistently(metrics).ShouldNot(gomega.Receive())
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})

		ginkgo.Context("when there are no instances", func() {
			ginkgo.BeforeEach(func() {
				deploymentInfo.Instances = []deployments.Instance{}
				deploymentsInfo = []deployments.DeploymentInfo{deploymentInfo}
			})

			ginkgo.It("returns only a last_jobs_scrape_timestamp & last_jobs_scrape_duration_seconds metric", func() {
				gomega.Eventually(metrics).Should(gomega.Receive())
				gomega.Eventually(metrics).Should(gomega.Receive())
				gomega.Consistently(metrics).ShouldNot(gomega.Receive())
				gomega.Consistently(errMetrics).ShouldNot(gomega.Receive())
			})
		})
	})
})
