package collectors

import (
	"strconv"
	"sync"
	"time"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/cloudfoundry-community/bosh_exporter/filters"
)

type ProcessesCollector struct {
	namespace                              string
	deploymentsFilter                      filters.DeploymentsFilter
	processHealthyDesc                     *prometheus.Desc
	processUptimeDesc                      *prometheus.Desc
	processCPUTotalDesc                    *prometheus.Desc
	processMemKBDesc                       *prometheus.Desc
	processMemPercentDesc                  *prometheus.Desc
	lastProcessesScrapeDurationSecondsDesc *prometheus.Desc
}

func NewProcessesCollector(
	namespace string,
	deploymentsFilter filters.DeploymentsFilter,
) *ProcessesCollector {
	processHealthyDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_process_healthy"),
		"BOSH Job Process Healthy.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip", "bosh_process"},
		nil,
	)

	processUptimeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_process_uptime_seconds"),
		"BOSH Job Process Uptime in seconds.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip", "bosh_process"},
		nil,
	)

	processCPUTotalDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_process_cpu_total"),
		"BOSH Job Process CPU Total.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip", "bosh_process"},
		nil,
	)

	processMemKBDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_process_mem_kb"),
		"BOSH Job Process Memory KB.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip", "bosh_process"},
		nil,
	)

	processMemPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_process_mem_percent"),
		"BOSH Job Process Memory Percent.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip", "bosh_process"},
		nil,
	)

	lastProcessesScrapeDurationSecondsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "last_job_processes_scrape_duration_seconds"),
		"Duration of the last scrape of Job Processes metrics from BOSH.",
		[]string{},
		nil,
	)

	collector := &ProcessesCollector{
		namespace:                              namespace,
		deploymentsFilter:                      deploymentsFilter,
		processHealthyDesc:                     processHealthyDesc,
		processUptimeDesc:                      processUptimeDesc,
		processCPUTotalDesc:                    processCPUTotalDesc,
		processMemKBDesc:                       processMemKBDesc,
		processMemPercentDesc:                  processMemPercentDesc,
		lastProcessesScrapeDurationSecondsDesc: lastProcessesScrapeDurationSecondsDesc,
	}
	return collector
}

func (c ProcessesCollector) Collect(ch chan<- prometheus.Metric) {
	var begun = time.Now()

	deployments := c.deploymentsFilter.GetDeployments()

	var wg sync.WaitGroup
	for _, deployment := range deployments {
		wg.Add(1)
		go func(deployment director.Deployment, ch chan<- prometheus.Metric) {
			defer wg.Done()
			c.reportProcessesMetrics(deployment, ch)
		}(deployment, ch)
	}
	wg.Wait()

	ch <- prometheus.MustNewConstMetric(
		c.lastProcessesScrapeDurationSecondsDesc,
		prometheus.GaugeValue,
		time.Since(begun).Seconds(),
	)
}

func (c ProcessesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.processHealthyDesc
	ch <- c.processUptimeDesc
	ch <- c.processCPUTotalDesc
	ch <- c.processMemKBDesc
	ch <- c.processMemPercentDesc
	ch <- c.lastProcessesScrapeDurationSecondsDesc
}

func (c ProcessesCollector) reportProcessesMetrics(
	deployment director.Deployment,
	ch chan<- prometheus.Metric,
) {
	log.Debugf("Reading VM info for deployment `%s`:", deployment.Name())
	vmInfos, err := deployment.VMInfos()
	if err != nil {
		log.Errorf("Error while reading VM info for deployment `%s`: %v", deployment.Name(), err)
		return
	}

	for _, vmInfo := range vmInfos {
		deploymentName := deployment.Name()
		jobName := vmInfo.JobName
		jobIndex := strconv.Itoa(int(*vmInfo.Index))
		jobAZ := vmInfo.AZ
		jobIP := ""
		if len(vmInfo.IPs) > 0 {
			jobIP = vmInfo.IPs[0]
		}

		for _, processInfo := range vmInfo.Processes {
			processName := processInfo.Name

			c.processHealthyMetrics(ch, processInfo.IsRunning(), deploymentName, jobName, jobIndex, jobAZ, jobIP, processName)
			c.processUptimeMetrics(ch, processInfo.Uptime, deploymentName, jobName, jobIndex, jobAZ, jobIP, processName)
			c.processCPUMetrics(ch, processInfo.CPU, deploymentName, jobName, jobIndex, jobAZ, jobIP, processName)
			c.processMemMetrics(ch, processInfo.Mem, deploymentName, jobName, jobIndex, jobAZ, jobIP, processName)
		}
	}
}

func (c ProcessesCollector) processHealthyMetrics(
	ch chan<- prometheus.Metric,
	processRunning bool,
	deploymentName string,
	jobName string,
	jobIndex string,
	jobAZ string,
	jobIP string,
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
		jobIP,
		processName,
	)
}

func (c ProcessesCollector) processUptimeMetrics(
	ch chan<- prometheus.Metric,
	uptime director.VMInfoVitalsUptime,
	deploymentName string,
	jobName string,
	jobIndex string,
	jobAZ string,
	jobIP string,
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
			jobIP,
			processName,
		)
	}
}

func (c ProcessesCollector) processCPUMetrics(
	ch chan<- prometheus.Metric,
	cpuMetrics director.VMInfoVitalsCPU,
	deploymentName string,
	jobName string,
	jobIndex string,
	jobAZ string,
	jobIP string,
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
			jobIP,
			processName,
		)
	}
}

func (c ProcessesCollector) processMemMetrics(
	ch chan<- prometheus.Metric,
	memMetrics director.VMInfoVitalsMemIntSize,
	deploymentName string,
	jobName string,
	jobIndex string,
	jobAZ string,
	jobIP string,
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
			jobIP,
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
			jobIP,
			processName,
		)
	}
}
