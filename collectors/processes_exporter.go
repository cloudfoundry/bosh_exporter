package collectors

import (
	"strconv"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type processesCollector struct {
	namespace             string
	boshClient            director.Director
	processHealthyDesc    *prometheus.Desc
	processUptimeDesc     *prometheus.Desc
	processCPUTotalDesc   *prometheus.Desc
	processMemKBDesc      *prometheus.Desc
	processMemPercentDesc *prometheus.Desc
}

func NewProcessesCollector(
	namespace string,
	boshClient director.Director,
) *processesCollector {
	processHealthyDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_process_healthy"),
		"BOSH Job Process Healthy.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "process_name"},
		nil,
	)

	processUptimeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_process_uptime_seconds"),
		"BOSH Job Process Uptime in seconds.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "process_name"},
		nil,
	)

	processCPUTotalDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_process_cpu_total"),
		"BOSH Job Process CPU Total.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "process_name"},
		nil,
	)

	processMemKBDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_process_mem_kb"),
		"BOSH Job Process Memory KB.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "process_name"},
		nil,
	)

	processMemPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_process_mem_percent"),
		"BOSH Job Process Memory Percent.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "process_name"},
		nil,
	)

	collector := &processesCollector{
		namespace:             namespace,
		boshClient:            boshClient,
		processHealthyDesc:    processHealthyDesc,
		processUptimeDesc:     processUptimeDesc,
		processCPUTotalDesc:   processCPUTotalDesc,
		processMemKBDesc:      processMemKBDesc,
		processMemPercentDesc: processMemPercentDesc,
	}
	return collector
}

func (c processesCollector) Collect(ch chan<- prometheus.Metric) {
	deployments, err := c.boshClient.Deployments()
	if err != nil {
		log.Errorf("Error while reading deployments: %v", err)
		return
	}

	for _, deployment := range deployments {
		vmInfos, err := deployment.VMInfos()
		if err != nil {
			log.Errorf("Error while reading VM info for deployment `%s`: %v", deployment.Name(), err)
			continue
		}

		for _, vmInfo := range vmInfos {
			deploymentName := deployment.Name()
			jobName := vmInfo.JobName
			jobIndex := strconv.Itoa(int(*vmInfo.Index))
			jobAZ := vmInfo.AZ

			for _, processInfo := range vmInfo.Processes {
				processName := processInfo.Name

				c.processHealthyMetrics(ch, processInfo.IsRunning(), deploymentName, jobName, jobIndex, jobAZ, processName)
				c.processUptimeMetrics(ch, processInfo.Uptime, deploymentName, jobName, jobIndex, jobAZ, processName)
				c.processCPUMetrics(ch, processInfo.CPU, deploymentName, jobName, jobIndex, jobAZ, processName)
				c.processMemMetrics(ch, processInfo.Mem, deploymentName, jobName, jobIndex, jobAZ, processName)
			}
		}
	}
}

func (c processesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.processHealthyDesc
	ch <- c.processUptimeDesc
	ch <- c.processCPUTotalDesc
	ch <- c.processMemKBDesc
	ch <- c.processMemPercentDesc
}

func (c processesCollector) processHealthyMetrics(
	ch chan<- prometheus.Metric,
	processRunning bool,
	deploymentName string,
	jobName string,
	jobIndex string,
	jobAZ string,
	processName string,
) {
	var runningMetric float64
	if processRunning {
		runningMetric = 1
	}

	ch <- prometheus.MustNewConstMetric(
		c.processHealthyDesc,
		prometheus.GaugeValue,
		runningMetric,
		deploymentName,
		jobName,
		jobIndex,
		jobAZ,
		processName,
	)
}

func (c processesCollector) processUptimeMetrics(
	ch chan<- prometheus.Metric,
	uptime director.VMInfoVitalsUptime,
	deploymentName string,
	jobName string,
	jobIndex string,
	jobAZ string,
	processName string,
) {
	if uptime.Seconds != nil {
		ch <- prometheus.MustNewConstMetric(
			c.processUptimeDesc,
			prometheus.GaugeValue,
			float64(*uptime.Seconds),
			deploymentName,
			jobName,
			jobIndex,
			jobAZ,
			processName,
		)
	}
}

func (c processesCollector) processCPUMetrics(
	ch chan<- prometheus.Metric,
	cpuMetrics director.VMInfoVitalsCPU,
	deploymentName string,
	jobName string,
	jobIndex string,
	jobAZ string,
	processName string,
) {
	if cpuMetrics.Total != nil {
		ch <- prometheus.MustNewConstMetric(
			c.processCPUTotalDesc,
			prometheus.GaugeValue,
			float64(*cpuMetrics.Total),
			deploymentName,
			jobName,
			jobIndex,
			jobAZ,
			processName,
		)
	}
}

func (c processesCollector) processMemMetrics(
	ch chan<- prometheus.Metric,
	memMetrics director.VMInfoVitalsMemIntSize,
	deploymentName string,
	jobName string,
	jobIndex string,
	jobAZ string,
	processName string,
) {
	if memMetrics.KB != nil {
		ch <- prometheus.MustNewConstMetric(
			c.processMemKBDesc,
			prometheus.GaugeValue,
			float64(*memMetrics.KB),
			deploymentName,
			jobName,
			jobIndex,
			jobAZ,
			processName,
		)
	}

	if memMetrics.Percent != nil {
		ch <- prometheus.MustNewConstMetric(
			c.processMemPercentDesc,
			prometheus.GaugeValue,
			*memMetrics.Percent,
			deploymentName,
			jobName,
			jobIndex,
			jobAZ,
			processName,
		)
	}
}
