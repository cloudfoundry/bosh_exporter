package collectors_test

import (
	"errors"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
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
		namespace = "test_exporter"
		boshDeployments = []string{}
		boshClient = &fakes.FakeDirector{}

		jobHealthyDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_healthy"),
			"BOSH Job Healthy.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobLoadAvg01Desc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_load_avg01"),
			"BOSH Job Load avg01.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobLoadAvg05Desc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_load_avg05"),
			"BOSH Job Load avg05.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobLoadAvg15Desc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_load_avg15"),
			"BOSH Job Load avg15.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobCPUSysDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_cpu_sys"),
			"BOSH Job CPU System.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobCPUUserDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_cpu_user"),
			"BOSH Job CPU User.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobCPUWaitDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_cpu_wait"),
			"BOSH Job CPU Wait.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobMemKBDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_mem_kb"),
			"BOSH Job Memory KB.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobMemPercentDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_mem_percent"),
			"BOSH Job Memory Percent.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobSwapKBDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_swap_kb"),
			"BOSH Job Swap KB.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobSwapPercentDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_swap_percent"),
			"BOSH Job Swap Percent.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobSystemDiskInodePercentDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_system_disk_inode_percent"),
			"BOSH Job System Disk Inode Percent.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobSystemDiskPercentDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_system_disk_percent"),
			"BOSH Job System Disk Percent.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobEphemeralDiskInodePercentDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_ephemeral_disk_inode_percent"),
			"BOSH Job Ephemeral Disk Inode Percent.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobEphemeralDiskPercentDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_ephemeral_disk_percent"),
			"BOSH Job Ephemeral Disk Percent.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobPersistentDiskInodePercentDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_persistent_disk_inode_percent"),
			"BOSH Job Persistent Disk Inode Percent.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)

		jobPersistentDiskPercentDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bosh", "job_persistent_disk_percent"),
			"BOSH Job Persistent Disk Percent.",
			[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
			nil,
		)
	})

	JustBeforeEach(func() {
		jobsCollector = collectors.NewJobsCollector(namespace, boshDeployments, boshClient)
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

	Describe("Collect", func() {
		var (
			deploymentName                = "fake-deployment-name"
			jobName                       = "fake-job-name"
			jobIndex                      = 0
			jobAZ                         = "fake-job-az"
			jobIP                         = "1.2.3.4"
			processState                  = "running"
			jobLoadAvg01                  = float64(0.01)
			jobLoadAvg05                  = float64(0.05)
			jobLoadAvg15                  = float64(0.15)
			jobCPUSys                     = float64(0.5)
			jobCPUUser                    = float64(1.0)
			jobCPUWait                    = float64(1.5)
			jobMemKB                      = 1000
			jobMemPercent                 = 10
			jobSwapKB                     = 2000
			jobSwapPercent                = 20
			jobSystemDiskInodePercent     = 10
			jobSystemDiskPercent          = 20
			jobEphemeralDiskInodePercent  = 30
			jobEphemeralDiskPercent       = 40
			jobPersistentDiskInodePercent = 50
			jobPersistentDiskPercent      = 60

			vmVitals    director.VMInfoVitals
			vmInfos     []director.VMInfo
			deployments []director.Deployment
			deployment  director.Deployment

			metrics                             chan prometheus.Metric
			jobHealthyMetric                    prometheus.Metric
			jobUnHealthyMetric                  prometheus.Metric
			jobLoadAvg01Metric                  prometheus.Metric
			jobLoadAvg05Metric                  prometheus.Metric
			jobLoadAvg15Metric                  prometheus.Metric
			jobCPUSysMetric                     prometheus.Metric
			jobCPUUserMetric                    prometheus.Metric
			jobCPUWaitMetric                    prometheus.Metric
			jobMemKBMetric                      prometheus.Metric
			jobMemPercentMetric                 prometheus.Metric
			jobSwapKBMetric                     prometheus.Metric
			jobSwapPercentMetric                prometheus.Metric
			jobSystemDiskInodePercentMetric     prometheus.Metric
			jobSystemDiskPercentMetric          prometheus.Metric
			jobEphemeralDiskInodePercentMetric  prometheus.Metric
			jobEphemeralDiskPercentMetric       prometheus.Metric
			jobPersistentDiskInodePercentMetric prometheus.Metric
			jobPersistentDiskPercentMetric      prometheus.Metric
		)

		BeforeEach(func() {
			vmVitals = director.VMInfoVitals{
				CPU: director.VMInfoVitalsCPU{
					Sys:  strconv.FormatFloat(jobCPUSys, 'E', -1, 64),
					User: strconv.FormatFloat(jobCPUUser, 'E', -1, 64),
					Wait: strconv.FormatFloat(jobCPUWait, 'E', -1, 64),
				},
				Mem: director.VMInfoVitalsMemSize{
					KB:      strconv.Itoa(jobMemKB),
					Percent: strconv.Itoa(jobMemPercent),
				},
				Swap: director.VMInfoVitalsMemSize{
					KB:      strconv.Itoa(jobSwapKB),
					Percent: strconv.Itoa(jobSwapPercent),
				},
				Load: []string{
					strconv.FormatFloat(jobLoadAvg01, 'E', -1, 64),
					strconv.FormatFloat(jobLoadAvg05, 'E', -1, 64),
					strconv.FormatFloat(jobLoadAvg15, 'E', -1, 64),
				},
				Disk: map[string]director.VMInfoVitalsDiskSize{
					"system": director.VMInfoVitalsDiskSize{
						InodePercent: strconv.Itoa(int(jobSystemDiskInodePercent)),
						Percent:      strconv.Itoa(int(jobSystemDiskPercent)),
					},
					"ephemeral": director.VMInfoVitalsDiskSize{
						InodePercent: strconv.Itoa(int(jobEphemeralDiskInodePercent)),
						Percent:      strconv.Itoa(int(jobEphemeralDiskPercent)),
					},
					"persistent": director.VMInfoVitalsDiskSize{
						InodePercent: strconv.Itoa(int(jobPersistentDiskInodePercent)),
						Percent:      strconv.Itoa(int(jobPersistentDiskPercent)),
					},
				},
			}

			vmInfos = []director.VMInfo{
				{
					JobName:      jobName,
					Index:        &jobIndex,
					ProcessState: processState,
					IPs:          []string{jobIP},
					AZ:           jobAZ,
					Vitals:       vmVitals,
				},
			}

			deployment = &fakes.FakeDeployment{
				NameStub:    func() string { return deploymentName },
				VMInfosStub: func() ([]director.VMInfo, error) { return vmInfos, nil },
			}

			deployments = []director.Deployment{deployment}
			boshClient.DeploymentsReturns(deployments, nil)

			metrics = make(chan prometheus.Metric)

			jobHealthyMetric = prometheus.MustNewConstMetric(
				jobHealthyDesc,
				prometheus.GaugeValue,
				float64(1),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobUnHealthyMetric = prometheus.MustNewConstMetric(
				jobHealthyDesc,
				prometheus.GaugeValue,
				float64(0),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobLoadAvg01Metric = prometheus.MustNewConstMetric(
				jobLoadAvg01Desc,
				prometheus.GaugeValue,
				jobLoadAvg01,
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobLoadAvg05Metric = prometheus.MustNewConstMetric(
				jobLoadAvg05Desc,
				prometheus.GaugeValue,
				jobLoadAvg05,
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobLoadAvg15Metric = prometheus.MustNewConstMetric(
				jobLoadAvg15Desc,
				prometheus.GaugeValue,
				jobLoadAvg15,
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobCPUSysMetric = prometheus.MustNewConstMetric(
				jobCPUSysDesc,
				prometheus.GaugeValue,
				jobCPUSys,
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobCPUUserMetric = prometheus.MustNewConstMetric(
				jobCPUUserDesc,
				prometheus.GaugeValue,
				jobCPUUser,
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobCPUWaitMetric = prometheus.MustNewConstMetric(
				jobCPUWaitDesc,
				prometheus.GaugeValue,
				jobCPUWait,
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobMemKBMetric = prometheus.MustNewConstMetric(
				jobMemKBDesc,
				prometheus.GaugeValue,
				float64(jobMemKB),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobMemPercentMetric = prometheus.MustNewConstMetric(
				jobMemPercentDesc,
				prometheus.GaugeValue,
				float64(jobMemPercent),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobSwapKBMetric = prometheus.MustNewConstMetric(
				jobSwapKBDesc,
				prometheus.GaugeValue,
				float64(jobSwapKB),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobSwapPercentMetric = prometheus.MustNewConstMetric(
				jobSwapPercentDesc,
				prometheus.GaugeValue,
				float64(jobSwapPercent),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobSystemDiskInodePercentMetric = prometheus.MustNewConstMetric(
				jobSystemDiskInodePercentDesc,
				prometheus.GaugeValue,
				float64(jobSystemDiskInodePercent),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobSystemDiskPercentMetric = prometheus.MustNewConstMetric(
				jobSystemDiskPercentDesc,
				prometheus.GaugeValue,
				float64(jobSystemDiskPercent),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobEphemeralDiskInodePercentMetric = prometheus.MustNewConstMetric(
				jobEphemeralDiskInodePercentDesc,
				prometheus.GaugeValue,
				float64(jobEphemeralDiskInodePercent),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobEphemeralDiskPercentMetric = prometheus.MustNewConstMetric(
				jobEphemeralDiskPercentDesc,
				prometheus.GaugeValue,
				float64(jobEphemeralDiskPercent),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobPersistentDiskInodePercentMetric = prometheus.MustNewConstMetric(
				jobPersistentDiskInodePercentDesc,
				prometheus.GaugeValue,
				float64(jobPersistentDiskInodePercent),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)

			jobPersistentDiskPercentMetric = prometheus.MustNewConstMetric(
				jobPersistentDiskPercentDesc,
				prometheus.GaugeValue,
				float64(jobPersistentDiskPercent),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				jobIP,
			)
		})

		JustBeforeEach(func() {
			go jobsCollector.Collect(metrics)
		})

		It("returns a bosh_job_process_healthy metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobHealthyMetric)))
		})

		Context("when the process is not running", func() {
			BeforeEach(func() {
				vmInfos[0].ProcessState = "failing"
			})

			It("returns a bosh_job_process_healthy metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobUnHealthyMetric)))
			})
		})

		It("returns a bosh_job_load_avg01 metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobLoadAvg01Metric)))
		})

		It("returns a bosh_job_load_avg05 metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobLoadAvg05Metric)))
		})

		It("returns a bosh_job_load_avg15 metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobLoadAvg15Metric)))
		})

		Context("when there is no load avg values", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.Load = []string{}
			})

			It("does not return any bosh_job_load_avg metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobLoadAvg01Metric)))
				Consistently(metrics).ShouldNot(Receive(Equal(jobLoadAvg05Metric)))
				Consistently(metrics).ShouldNot(Receive(Equal(jobLoadAvg15Metric)))
			})
		})

		It("returns a bosh_job_cpu_sys metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobCPUSysMetric)))
		})

		Context("when there is no cpu sys value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.CPU = director.VMInfoVitalsCPU{
					User: strconv.FormatFloat(jobCPUUser, 'E', -1, 64),
					Wait: strconv.FormatFloat(jobCPUWait, 'E', -1, 64),
				}
			})

			It("does not return a bosh_job_cpu_sys metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobCPUSysMetric)))
			})
		})

		It("returns a bosh_job_cpu_user metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobCPUUserMetric)))
		})

		Context("when there is no cpu user value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.CPU = director.VMInfoVitalsCPU{
					Sys:  strconv.FormatFloat(jobCPUSys, 'E', -1, 64),
					Wait: strconv.FormatFloat(jobCPUWait, 'E', -1, 64),
				}
			})

			It("does not return a bosh_job_cpu_user metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobCPUUserMetric)))
			})
		})

		It("returns a bosh_job_cpu_wait metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobCPUWaitMetric)))
		})

		Context("when there is no cpu wait value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.CPU = director.VMInfoVitalsCPU{
					Sys:  strconv.FormatFloat(jobCPUSys, 'E', -1, 64),
					User: strconv.FormatFloat(jobCPUUser, 'E', -1, 64),
				}
			})

			It("does not return a bosh_job_cpu_wait metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobCPUWaitMetric)))
			})
		})

		It("returns a bosh_job_mem_kb metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobMemKBMetric)))
		})

		Context("when there is no mem kb value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.Mem = director.VMInfoVitalsMemSize{
					Percent: strconv.Itoa(jobMemPercent),
				}
			})

			It("does not return a bosh_job_mem_kb metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobMemKBMetric)))
			})
		})

		It("returns a bosh_job_mem_percent metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobMemPercentMetric)))
		})

		Context("when there is no mem percent value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.Mem = director.VMInfoVitalsMemSize{
					KB: strconv.Itoa(jobMemKB),
				}
			})

			It("does not return a bosh_job_mem_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobMemPercentMetric)))
			})
		})

		It("returns a bosh_job_swap_kb metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobSwapKBMetric)))
		})

		Context("when there is no swap kb value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.Swap = director.VMInfoVitalsMemSize{
					Percent: strconv.Itoa(jobSwapPercent),
				}
			})

			It("does not return a bosh_job_swap_kb metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobSwapKBMetric)))
			})
		})

		It("returns a bosh_job_swap_percent metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobSwapPercentMetric)))
		})

		Context("when there is no swap percent value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.Swap = director.VMInfoVitalsMemSize{
					KB: strconv.Itoa(jobSwapKB),
				}
			})

			It("does not return a bosh_job_swap_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobSwapPercentMetric)))
			})
		})

		It("returns a bosh_job_system_disk_inode_percent metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobSystemDiskInodePercentMetric)))
		})

		Context("when there is no system disk inode percent value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.Disk["system"] = director.VMInfoVitalsDiskSize{
					Percent: strconv.Itoa(int(jobSystemDiskPercent)),
				}
			})

			It("does not return a bosh_job_system_disk_inode_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobSystemDiskInodePercentMetric)))
			})
		})

		It("returns a bosh_job_system_disk_percent metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobSystemDiskPercentMetric)))
		})

		Context("when there is no system disk percent value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.Disk["system"] = director.VMInfoVitalsDiskSize{
					InodePercent: strconv.Itoa(int(jobSystemDiskInodePercent)),
				}
			})

			It("does not return a bosh_job_system_disk_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobSystemDiskPercentMetric)))
			})
		})

		It("returns a bosh_job_ephemeral_disk_inode_percent metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobEphemeralDiskInodePercentMetric)))
		})

		Context("when there is no ephemeral disk inode percent value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.Disk["ephemeral"] = director.VMInfoVitalsDiskSize{
					Percent: strconv.Itoa(int(jobEphemeralDiskPercent)),
				}
			})

			It("does not return a bosh_job_ephemeral_disk_inode_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobEphemeralDiskInodePercentMetric)))
			})
		})

		It("returns a bosh_job_ephemeral_disk_percent metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobEphemeralDiskPercentMetric)))
		})

		Context("when there is no ephemeral disk percent value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.Disk["ephemeral"] = director.VMInfoVitalsDiskSize{
					InodePercent: strconv.Itoa(int(jobEphemeralDiskInodePercent)),
				}
			})

			It("does not return a bosh_job_Ephemeral_disk_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobEphemeralDiskPercentMetric)))
			})
		})

		It("returns a bosh_job_persistent_disk_inode_percent metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobPersistentDiskInodePercentMetric)))
		})

		Context("when there is no persistent disk inode percent value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.Disk["persistent"] = director.VMInfoVitalsDiskSize{
					Percent: strconv.Itoa(int(jobPersistentDiskPercent)),
				}
			})

			It("does not return a bosh_job_persistent_disk_inode_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobPersistentDiskInodePercentMetric)))
			})
		})

		It("returns a bosh_job_persistent_disk_percent metric", func() {
			Eventually(metrics).Should(Receive(Equal(jobPersistentDiskPercentMetric)))
		})

		Context("when there is no persistent disk percent value", func() {
			BeforeEach(func() {
				vmInfos[0].Vitals.Disk["persistent"] = director.VMInfoVitalsDiskSize{
					InodePercent: strconv.Itoa(int(jobPersistentDiskInodePercent)),
				}
			})

			It("does not return a bosh_job_persistent_disk_percent metric", func() {
				Consistently(metrics).ShouldNot(Receive(Equal(jobPersistentDiskPercentMetric)))
			})
		})

		Context("when there are no deployments", func() {
			BeforeEach(func() {
				boshClient.DeploymentsReturns([]director.Deployment{}, nil)
			})

			It("does not return any metric", func() {
				Consistently(metrics).ShouldNot(Receive())
			})
		})

		Context("when it fails to get the deployments", func() {
			BeforeEach(func() {
				boshClient.DeploymentsReturns(nil, errors.New("no deployments"))
			})

			It("does not return any metric", func() {
				Consistently(metrics).ShouldNot(Receive())
			})
		})

		Context("when it dos not return any VMInfos", func() {
			BeforeEach(func() {
				deployment = &fakes.FakeDeployment{
					NameStub:    func() string { return deploymentName },
					VMInfosStub: func() ([]director.VMInfo, error) { return nil, nil },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("does not return any metric", func() {
				Consistently(metrics).ShouldNot(Receive())
			})
		})

		Context("when it fails to get the VMInfos for a deployment", func() {
			BeforeEach(func() {
				deployment = &fakes.FakeDeployment{
					NameStub:    func() string { return deploymentName },
					VMInfosStub: func() ([]director.VMInfo, error) { return nil, errors.New("no VMInfo") },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("does not return any metric", func() {
				Consistently(metrics).ShouldNot(Receive())
			})
		})

		Context("when there is a bosh deployment filter", func() {
			BeforeEach(func() {
				boshDeployments = []string{"fake-deployment-name"}
				boshClient.FindDeploymentReturns(deployment, nil)
			})

			It("returns a bosh_job_process_healthy metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobHealthyMetric)))
			})

			It("returns a bosh_job_load_avg01 metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobLoadAvg01Metric)))
			})

			It("returns a bosh_job_load_avg05 metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobLoadAvg05Metric)))
			})

			It("returns a bosh_job_load_avg15 metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobLoadAvg15Metric)))
			})

			It("returns a bosh_job_cpu_sys metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobCPUSysMetric)))
			})

			It("returns a bosh_job_cpu_user metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobCPUUserMetric)))
			})

			It("returns a bosh_job_cpu_wait metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobCPUWaitMetric)))
			})

			It("returns a bosh_job_mem_kb metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobMemKBMetric)))
			})

			It("returns a bosh_job_mem_percent metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobMemPercentMetric)))
			})

			It("returns a bosh_job_swap_kb metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobSwapKBMetric)))
			})

			It("returns a bosh_job_swap_percent metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobSwapPercentMetric)))
			})

			It("returns a bosh_job_system_disk_inode_percent metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobSystemDiskInodePercentMetric)))
			})

			It("returns a bosh_job_system_disk_percent metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobSystemDiskPercentMetric)))
			})

			It("returns a bosh_job_ephemeral_disk_inode_percent metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobEphemeralDiskInodePercentMetric)))
			})

			It("returns a bosh_job_ephemeral_disk_percent metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobEphemeralDiskPercentMetric)))
			})

			It("returns a bosh_job_persistent_disk_inode_percent metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobPersistentDiskInodePercentMetric)))
			})

			It("returns a bosh_job_persistent_disk_percent metric", func() {
				Eventually(metrics).Should(Receive(Equal(jobPersistentDiskPercentMetric)))
			})

			Context("and no deployment matches", func() {
				BeforeEach(func() {
					boshClient.FindDeploymentReturns(nil, errors.New("does not exists"))
				})

				It("does not return any metric", func() {
					Consistently(metrics).ShouldNot(Receive())
				})
			})
		})
	})
})
