package collectors_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/director/fakes"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/collectors"
)

var _ = Describe("ProcessesCollector", func() {
	var (
		processHealthyDesc    *prometheus.Desc
		processUptimeDesc     *prometheus.Desc
		processCPUTotalDesc   *prometheus.Desc
		processMemKBDesc      *prometheus.Desc
		processMemPercentDesc *prometheus.Desc

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
	})

	Describe("Describe", func() {
		var (
			descriptions chan *prometheus.Desc
		)

		BeforeEach(func() {
			descriptions = make(chan *prometheus.Desc)
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

	Describe("Collect", func() {
		var (
			deploymentName       = "fake-deployment-name"
			jobName              = "fake-job-name"
			jobIndex             = 0
			jobAZ                = "fake-job-az"
			processName          = "fake-process-name"
			processState         = "running"
			processUptimeSeconds = uint64(3600)
			processCPUTotal      = float64(0.5)
			processMemKB         = uint64(2000)
			processMemPercent    = float64(20)

			deployments []director.Deployment
			vmInfos     []director.VMInfo
			vmProcesses []director.VMInfoProcess

			metrics                 chan prometheus.Metric
			processHealthyMetric    prometheus.Metric
			processUptimeMetric     prometheus.Metric
			processCPUTotalMetric   prometheus.Metric
			processMemKBMetric      prometheus.Metric
			processMemPercentMetric prometheus.Metric
		)

		BeforeEach(func() {
			processHealthyMetric = prometheus.MustNewConstMetric(
				processHealthyDesc,
				prometheus.GaugeValue,
				float64(1),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				processName,
			)

			processUptimeMetric = prometheus.MustNewConstMetric(
				processUptimeDesc,
				prometheus.GaugeValue,
				float64(processUptimeSeconds),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				processName,
			)

			processCPUTotalMetric = prometheus.MustNewConstMetric(
				processCPUTotalDesc,
				prometheus.GaugeValue,
				processCPUTotal,
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				processName,
			)

			processMemKBMetric = prometheus.MustNewConstMetric(
				processMemKBDesc,
				prometheus.GaugeValue,
				float64(processMemKB),
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				processName,
			)

			processMemPercentMetric = prometheus.MustNewConstMetric(
				processMemPercentDesc,
				prometheus.GaugeValue,
				processMemPercent,
				deploymentName,
				jobName,
				strconv.Itoa(jobIndex),
				jobAZ,
				processName,
			)
		})

		JustBeforeEach(func() {
			vmProcesses = []director.VMInfoProcess{
				{
					Name:   processName,
					State:  processState,
					CPU:    director.VMInfoVitalsCPU{Total: &processCPUTotal},
					Mem:    director.VMInfoVitalsMemIntSize{KB: &processMemKB, Percent: &processMemPercent},
					Uptime: director.VMInfoVitalsUptime{Seconds: &processUptimeSeconds},
				},
			}

			vmInfos = []director.VMInfo{
				{
					JobName:   jobName,
					Index:     &jobIndex,
					AZ:        jobAZ,
					Processes: vmProcesses,
				},
			}

			deployments = []director.Deployment{
				&fakes.FakeDeployment{
					NameStub:    func() string { return deploymentName },
					VMInfosStub: func() ([]director.VMInfo, error) { return vmInfos, nil },
				},
			}

			boshClient.DeploymentsReturns(deployments, nil)

			metrics = make(chan prometheus.Metric)
			go processesCollector.Collect(metrics)
		})

		It("returns a bosh_job_process_healthy metric", func() {
			Eventually(metrics).Should(Receive(Equal(processHealthyMetric)))
		})

		Context("when the process is not running", func() {
			BeforeEach(func() {
				processState = "failing"

				processHealthyMetric = prometheus.MustNewConstMetric(
					processHealthyDesc,
					prometheus.GaugeValue,
					float64(0),
					deploymentName,
					jobName,
					strconv.Itoa(jobIndex),
					jobAZ,
					processName,
				)
			})

			It("returns a bosh_job_process_healthy metric", func() {
				Eventually(metrics).Should(Receive(Equal(processHealthyMetric)))
			})
		})

		It("returns a bosh_job_process_uptime_seconds metric", func() {
			Eventually(metrics).Should(Receive(Equal(processUptimeMetric)))
		})

		It("returns a bosh_job_process_cpu_total metric", func() {
			Eventually(metrics).Should(Receive(Equal(processCPUTotalMetric)))
		})

		It("returns a bosh_job_process_mem_kb metric metric", func() {
			Eventually(metrics).Should(Receive(Equal(processMemKBMetric)))
		})

		It("returns a bosh_job_process_mem_percent metric", func() {
			Eventually(metrics).Should(Receive(Equal(processMemPercentMetric)))
		})
	})
})
