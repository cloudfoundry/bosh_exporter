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

type JobsCollector struct {
	namespace                         string
	deploymentsFilter                 filters.DeploymentsFilter
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
	jobProcessHealthyDesc             *prometheus.Desc
	jobProcessUptimeDesc              *prometheus.Desc
	jobProcessCPUTotalDesc            *prometheus.Desc
	jobProcessMemKBDesc               *prometheus.Desc
	jobProcessMemPercentDesc          *prometheus.Desc
	lastJobsScrapeTimestampDesc       *prometheus.Desc
	lastJobsScrapeDurationSecondsDesc *prometheus.Desc
}

func NewJobsCollector(
	namespace string,
	deploymentsFilter filters.DeploymentsFilter,
) *JobsCollector {
	jobHealthyDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "healthy"),
		"BOSH Job Healthy (1 for healthy, 0 for unhealthy).",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobLoadAvg01Desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "load_avg01"),
		"BOSH Job Load avg01.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobLoadAvg05Desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "load_avg05"),
		"BOSH Job Load avg05.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobLoadAvg15Desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "load_avg15"),
		"BOSH Job Load avg15.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobCPUSysDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "cpu_sys"),
		"BOSH Job CPU System.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobCPUUserDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "cpu_user"),
		"BOSH Job CPU User.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobCPUWaitDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "cpu_wait"),
		"BOSH Job CPU Wait.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobMemKBDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "mem_kb"),
		"BOSH Job Memory KB.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobMemPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "mem_percent"),
		"BOSH Job Memory Percent.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobSwapKBDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "swap_kb"),
		"BOSH Job Swap KB.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobSwapPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "swap_percent"),
		"BOSH Job Swap Percent.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobSystemDiskInodePercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "system_disk_inode_percent"),
		"BOSH Job System Disk Inode Percent.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobSystemDiskPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "system_disk_percent"),
		"BOSH Job System Disk Percent.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobEphemeralDiskInodePercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "ephemeral_disk_inode_percent"),
		"BOSH Job Ephemeral Disk Inode Percent.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobEphemeralDiskPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "ephemeral_disk_percent"),
		"BOSH Job Ephemeral Disk Percent.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobPersistentDiskInodePercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "persistent_disk_inode_percent"),
		"BOSH Job Persistent Disk Inode Percent.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobPersistentDiskPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job", "persistent_disk_percent"),
		"BOSH Job Persistent Disk Percent.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
		nil,
	)

	jobProcessHealthyDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job_process", "healthy"),
		"BOSH Job Process Healthy (1 for healthy, 0 for unhealthy).",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
		nil,
	)

	jobProcessUptimeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job_process", "uptime_seconds"),
		"BOSH Job Process Uptime in seconds.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
		nil,
	)

	jobProcessCPUTotalDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job_process", "cpu_total"),
		"BOSH Job Process CPU Total.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
		nil,
	)

	jobProcessMemKBDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job_process", "mem_kb"),
		"BOSH Job Process Memory KB.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
		nil,
	)

	jobProcessMemPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "job_process", "mem_percent"),
		"BOSH Job Process Memory Percent.",
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
		nil,
	)

	lastJobsScrapeTimestampDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "last_jobs_scrape_timestamp"),
		"Number of seconds since 1970 since last scrape of Job metrics from BOSH.",
		[]string{},
		nil,
	)

	lastJobsScrapeDurationSecondsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "last_jobs_scrape_duration_seconds"),
		"Duration of the last scrape of Job metrics from BOSH.",
		[]string{},
		nil,
	)

	collector := &JobsCollector{
		namespace:                         namespace,
		deploymentsFilter:                 deploymentsFilter,
		jobHealthyDesc:                    jobHealthyDesc,
		jobLoadAvg01Desc:                  jobLoadAvg01Desc,
		jobLoadAvg05Desc:                  jobLoadAvg05Desc,
		jobLoadAvg15Desc:                  jobLoadAvg15Desc,
		jobCPUSysDesc:                     jobCPUSysDesc,
		jobCPUUserDesc:                    jobCPUUserDesc,
		jobCPUWaitDesc:                    jobCPUWaitDesc,
		jobMemKBDesc:                      jobMemKBDesc,
		jobMemPercentDesc:                 jobMemPercentDesc,
		jobSwapKBDesc:                     jobSwapKBDesc,
		jobSwapPercentDesc:                jobSwapPercentDesc,
		jobSystemDiskInodePercentDesc:     jobSystemDiskInodePercentDesc,
		jobSystemDiskPercentDesc:          jobSystemDiskPercentDesc,
		jobEphemeralDiskInodePercentDesc:  jobEphemeralDiskInodePercentDesc,
		jobEphemeralDiskPercentDesc:       jobEphemeralDiskPercentDesc,
		jobPersistentDiskInodePercentDesc: jobPersistentDiskInodePercentDesc,
		jobPersistentDiskPercentDesc:      jobPersistentDiskPercentDesc,
		jobProcessHealthyDesc:             jobProcessHealthyDesc,
		jobProcessUptimeDesc:              jobProcessUptimeDesc,
		jobProcessCPUTotalDesc:            jobProcessCPUTotalDesc,
		jobProcessMemKBDesc:               jobProcessMemKBDesc,
		jobProcessMemPercentDesc:          jobProcessMemPercentDesc,
		lastJobsScrapeTimestampDesc:       lastJobsScrapeTimestampDesc,
		lastJobsScrapeDurationSecondsDesc: lastJobsScrapeDurationSecondsDesc,
	}
	return collector
}

func (c JobsCollector) Collect(ch chan<- prometheus.Metric) {
	var begun = time.Now()

	deployments := c.deploymentsFilter.GetDeployments()

	var wg sync.WaitGroup
	for _, deployment := range deployments {
		wg.Add(1)
		go func(deployment director.Deployment, ch chan<- prometheus.Metric) {
			defer wg.Done()
			c.reportJobMetrics(deployment, ch)
		}(deployment, ch)
	}
	wg.Wait()

	ch <- prometheus.MustNewConstMetric(
		c.lastJobsScrapeTimestampDesc,
		prometheus.GaugeValue,
		float64(time.Now().Unix()),
	)

	ch <- prometheus.MustNewConstMetric(
		c.lastJobsScrapeDurationSecondsDesc,
		prometheus.GaugeValue,
		time.Since(begun).Seconds(),
	)
}

func (c JobsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.jobHealthyDesc
	ch <- c.jobLoadAvg01Desc
	ch <- c.jobLoadAvg05Desc
	ch <- c.jobLoadAvg15Desc
	ch <- c.jobCPUSysDesc
	ch <- c.jobCPUUserDesc
	ch <- c.jobCPUWaitDesc
	ch <- c.jobMemKBDesc
	ch <- c.jobMemPercentDesc
	ch <- c.jobSwapKBDesc
	ch <- c.jobSwapPercentDesc
	ch <- c.jobSystemDiskInodePercentDesc
	ch <- c.jobSystemDiskPercentDesc
	ch <- c.jobEphemeralDiskInodePercentDesc
	ch <- c.jobEphemeralDiskPercentDesc
	ch <- c.jobPersistentDiskInodePercentDesc
	ch <- c.jobPersistentDiskPercentDesc
	ch <- c.jobProcessHealthyDesc
	ch <- c.jobProcessUptimeDesc
	ch <- c.jobProcessCPUTotalDesc
	ch <- c.jobProcessMemKBDesc
	ch <- c.jobProcessMemPercentDesc
	ch <- c.lastJobsScrapeTimestampDesc
	ch <- c.lastJobsScrapeDurationSecondsDesc
}

func (c JobsCollector) reportJobMetrics(
	deployment director.Deployment,
	ch chan<- prometheus.Metric,
) {
	log.Debugf("Reading VM info for deployment `%s`:", deployment.Name())
	instanceInfos, err := deployment.InstanceInfos()
	if err != nil {
		log.Errorf("Error while reading VM info for deployment `%s`: %v", deployment.Name(), err)
		return
	}

	for _, instanceInfo := range instanceInfos {
		if instanceInfo.VMID == "" {
			continue
		}

		deploymentName := deployment.Name()
		jobName := instanceInfo.JobName
		jobID := instanceInfo.ID
		jobIndex := strconv.Itoa(int(*instanceInfo.Index))
		jobAZ := instanceInfo.AZ
		jobIP := ""
		if len(instanceInfo.IPs) > 0 {
			jobIP = instanceInfo.IPs[0]
		}

		c.jobHealthyMetrics(ch, instanceInfo.IsRunning(), deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		c.jobLoadAvgMetrics(ch, instanceInfo.Vitals.Load, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		c.jobCPUMetrics(ch, instanceInfo.Vitals.CPU, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		c.jobMemMetrics(ch, instanceInfo.Vitals.Mem, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		c.jobSwapMetrics(ch, instanceInfo.Vitals.Swap, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		c.jobSystemDiskMetrics(ch, instanceInfo.Vitals.SystemDisk(), deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		c.jobEphemeralDiskMetrics(ch, instanceInfo.Vitals.EphemeralDisk(), deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		c.jobPersistentDiskMetrics(ch, instanceInfo.Vitals.PersistentDisk(), deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)

		for _, jobProcessInfo := range instanceInfo.Processes {
			jobProcessName := jobProcessInfo.Name

			c.jobProcessHealthyMetrics(ch, jobProcessInfo.IsRunning(), deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
			c.jobProcessUptimeMetrics(ch, jobProcessInfo.Uptime, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
			c.jobProcessCPUMetrics(ch, jobProcessInfo.CPU, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
			c.jobProcessMemMetrics(ch, jobProcessInfo.Mem, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
		}
	}
}

func (c JobsCollector) jobHealthyMetrics(
	ch chan<- prometheus.Metric,
	vmRunning bool,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) {
	var runningMetric float64
	if vmRunning {
		runningMetric = 1
	}

	ch <- prometheus.MustNewConstMetric(
		c.jobHealthyDesc,
		prometheus.GaugeValue,
		runningMetric,
		deploymentName,
		jobName,
		jobID,
		jobIndex,
		jobAZ,
		jobIP,
	)
}

func (c JobsCollector) jobLoadAvgMetrics(
	ch chan<- prometheus.Metric,
	loadAvg []string,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) {
	if len(loadAvg) == 3 {
		if loadAvg[0] != "" {
			loadAvg01, err := strconv.ParseFloat(loadAvg[0], 64)
			if err != nil {
				log.Errorf("Error while converting Load avg01 metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.jobLoadAvg01Desc,
					prometheus.GaugeValue,
					float64(loadAvg01),
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				)
			}
		}

		if loadAvg[1] != "" {
			loadAvg05, err := strconv.ParseFloat(loadAvg[1], 64)
			if err != nil {
				log.Errorf("Error while converting Load avg05 metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.jobLoadAvg05Desc,
					prometheus.GaugeValue,
					float64(loadAvg05),
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				)
			}
		}

		if loadAvg[2] != "" {
			loadAvg15, err := strconv.ParseFloat(loadAvg[2], 64)
			if err != nil {
				log.Errorf("Error while converting Load avg15 metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.jobLoadAvg15Desc,
					prometheus.GaugeValue,
					float64(loadAvg15),
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				)
			}
		}
	}
}

func (c JobsCollector) jobCPUMetrics(
	ch chan<- prometheus.Metric,
	cpuMetrics director.VMInfoVitalsCPU,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) {
	if cpuMetrics.Sys != "" {
		cpuSys, err := strconv.ParseFloat(cpuMetrics.Sys, 64)
		if err != nil {
			log.Errorf("Error while converting CPU Sys metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobCPUSysDesc,
				prometheus.GaugeValue,
				cpuSys,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}

	if cpuMetrics.User != "" {
		cpuUser, err := strconv.ParseFloat(cpuMetrics.User, 64)
		if err != nil {
			log.Errorf("Error while converting CPU User metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobCPUUserDesc,
				prometheus.GaugeValue,
				cpuUser,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}

	if cpuMetrics.Wait != "" {
		cpuWait, err := strconv.ParseFloat(cpuMetrics.Wait, 64)
		if err != nil {
			log.Errorf("Error while converting CPU Wait metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobCPUWaitDesc,
				prometheus.GaugeValue,
				cpuWait,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}
}

func (c JobsCollector) jobMemMetrics(
	ch chan<- prometheus.Metric,
	memMetrics director.VMInfoVitalsMemSize,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) {
	if memMetrics.KB != "" {
		memKB, err := strconv.ParseFloat(memMetrics.KB, 64)
		if err != nil {
			log.Errorf("Error while converting Mem KB metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobMemKBDesc,
				prometheus.GaugeValue,
				memKB,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}

	if memMetrics.Percent != "" {
		memPercent, err := strconv.ParseFloat(memMetrics.Percent, 64)
		if err != nil {
			log.Errorf("Error while converting Mem Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobMemPercentDesc,
				prometheus.GaugeValue,
				memPercent,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}
}

func (c JobsCollector) jobSwapMetrics(
	ch chan<- prometheus.Metric,
	swapMetrics director.VMInfoVitalsMemSize,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) {
	if swapMetrics.KB != "" {
		swapKB, err := strconv.ParseFloat(swapMetrics.KB, 64)
		if err != nil {
			log.Errorf("Error while converting Swap KB metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobSwapKBDesc,
				prometheus.GaugeValue,
				swapKB,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}

	if swapMetrics.Percent != "" {
		swapPercent, err := strconv.ParseFloat(swapMetrics.Percent, 64)
		if err != nil {
			log.Errorf("Error while converting Swap Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobSwapPercentDesc,
				prometheus.GaugeValue,
				swapPercent,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}
}

func (c JobsCollector) jobSystemDiskMetrics(
	ch chan<- prometheus.Metric,
	systemDiskMetrics director.VMInfoVitalsDiskSize,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) {
	if systemDiskMetrics.InodePercent != "" {
		systemDiskInodePercent, err := strconv.ParseFloat(systemDiskMetrics.InodePercent, 64)
		if err != nil {
			log.Errorf("Error while converting System Disk Inode Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobSystemDiskInodePercentDesc,
				prometheus.GaugeValue,
				systemDiskInodePercent,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}

	if systemDiskMetrics.Percent != "" {
		systemDiskPercent, err := strconv.ParseFloat(systemDiskMetrics.Percent, 64)
		if err != nil {
			log.Errorf("Error while converting System Disk Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobSystemDiskPercentDesc,
				prometheus.GaugeValue,
				systemDiskPercent,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}
}

func (c JobsCollector) jobEphemeralDiskMetrics(
	ch chan<- prometheus.Metric,
	ephemeralDiskMetrics director.VMInfoVitalsDiskSize,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) {
	if ephemeralDiskMetrics.InodePercent != "" {
		ephemeralDiskInodePercent, err := strconv.ParseFloat(ephemeralDiskMetrics.InodePercent, 64)
		if err != nil {
			log.Errorf("Error while converting Ephemeral Disk Inode Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {

			ch <- prometheus.MustNewConstMetric(
				c.jobEphemeralDiskInodePercentDesc,
				prometheus.GaugeValue,
				ephemeralDiskInodePercent,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}

	if ephemeralDiskMetrics.Percent != "" {
		ephemeralDiskPercent, err := strconv.ParseFloat(ephemeralDiskMetrics.Percent, 64)
		if err != nil {
			log.Errorf("Error while converting Ephemeral Disk Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobEphemeralDiskPercentDesc,
				prometheus.GaugeValue,
				ephemeralDiskPercent,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}
}

func (c JobsCollector) jobPersistentDiskMetrics(
	ch chan<- prometheus.Metric,
	persistentDiskMetrics director.VMInfoVitalsDiskSize,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) {
	if persistentDiskMetrics.InodePercent != "" {
		persistentDiskInodePercent, err := strconv.ParseFloat(persistentDiskMetrics.InodePercent, 64)
		if err != nil {
			log.Errorf("Error while converting Persistent Disk Inode Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobPersistentDiskInodePercentDesc,
				prometheus.GaugeValue,
				persistentDiskInodePercent,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}

	if persistentDiskMetrics.Percent != "" {
		persistentDiskPercent, err := strconv.ParseFloat(persistentDiskMetrics.Percent, 64)
		if err != nil {
			log.Errorf("Error while converting Persistent Disk Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.jobPersistentDiskPercentDesc,
				prometheus.GaugeValue,
				persistentDiskPercent,
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}
}

func (c JobsCollector) jobProcessHealthyMetrics(
	ch chan<- prometheus.Metric,
	processRunning bool,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) {
	var runningMetric float64
	if processRunning {
		runningMetric = 1
	}

	ch <- prometheus.MustNewConstMetric(
		c.jobProcessHealthyDesc,
		prometheus.GaugeValue,
		runningMetric,
		deploymentName,
		jobName,
		jobID,
		jobIndex,
		jobAZ,
		jobIP,
		jobProcessName,
	)
}

func (c JobsCollector) jobProcessUptimeMetrics(
	ch chan<- prometheus.Metric,
	uptime director.VMInfoVitalsUptime,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) {
	if uptime.Seconds != nil {
		ch <- prometheus.MustNewConstMetric(
			c.jobProcessUptimeDesc,
			prometheus.GaugeValue,
			float64(*uptime.Seconds),
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		)
	}
}

func (c JobsCollector) jobProcessCPUMetrics(
	ch chan<- prometheus.Metric,
	cpuMetrics director.VMInfoVitalsCPU,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) {
	if cpuMetrics.Total != nil {
		ch <- prometheus.MustNewConstMetric(
			c.jobProcessCPUTotalDesc,
			prometheus.GaugeValue,
			float64(*cpuMetrics.Total),
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		)
	}
}

func (c JobsCollector) jobProcessMemMetrics(
	ch chan<- prometheus.Metric,
	memMetrics director.VMInfoVitalsMemIntSize,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) {
	if memMetrics.KB != nil {
		ch <- prometheus.MustNewConstMetric(
			c.jobProcessMemKBDesc,
			prometheus.GaugeValue,
			float64(*memMetrics.KB),
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		)
	}

	if memMetrics.Percent != nil {
		ch <- prometheus.MustNewConstMetric(
			c.jobProcessMemPercentDesc,
			prometheus.GaugeValue,
			*memMetrics.Percent,
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		)
	}
}
