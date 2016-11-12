package collectors

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/deployments"
)

type JobsCollector struct {
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

func NewJobsCollector(namespace string) *JobsCollector {
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

func (c *JobsCollector) Collect(deployments []deployments.DeploymentInfo, ch chan<- prometheus.Metric) error {
	var err error
	var begun = time.Now()

	for _, deployment := range deployments {
		err = c.reportJobMetrics(deployment, ch)
	}

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

	return err
}

func (c *JobsCollector) Describe(ch chan<- *prometheus.Desc) {
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

func (c *JobsCollector) reportJobMetrics(deployment deployments.DeploymentInfo, ch chan<- prometheus.Metric) error {
	var err error

	for _, instance := range deployment.Instances {
		deploymentName := deployment.Name
		jobName := instance.Name
		jobID := instance.ID
		jobIndex := instance.Index
		jobAZ := instance.AZ
		jobIP := ""
		if len(instance.IPs) > 0 {
			jobIP = instance.IPs[0]
		}

		err = c.jobHealthyMetrics(ch, instance.Healthy, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		err = c.jobLoadAvgMetrics(ch, instance.Vitals.Load, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		err = c.jobCPUMetrics(ch, instance.Vitals.CPU, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		err = c.jobMemMetrics(ch, instance.Vitals.Mem, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		err = c.jobSwapMetrics(ch, instance.Vitals.Swap, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		err = c.jobSystemDiskMetrics(ch, instance.Vitals.SystemDisk, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		err = c.jobEphemeralDiskMetrics(ch, instance.Vitals.EphemeralDisk, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		err = c.jobPersistentDiskMetrics(ch, instance.Vitals.PersistentDisk, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)

		for _, process := range instance.Processes {
			jobProcessName := process.Name

			err = c.jobProcessHealthyMetrics(ch, process.Healthy, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
			err = c.jobProcessUptimeMetrics(ch, process.Uptime, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
			err = c.jobProcessCPUMetrics(ch, process.CPU, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
			err = c.jobProcessMemMetrics(ch, process.Mem, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
		}
	}

	return err
}

func (c *JobsCollector) jobHealthyMetrics(
	ch chan<- prometheus.Metric,
	healthy bool,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var healthyMetric float64
	if healthy {
		healthyMetric = 1
	}

	ch <- prometheus.MustNewConstMetric(
		c.jobHealthyDesc,
		prometheus.GaugeValue,
		healthyMetric,
		deploymentName,
		jobName,
		jobID,
		jobIndex,
		jobAZ,
		jobIP,
	)

	return nil
}

func (c *JobsCollector) jobLoadAvgMetrics(
	ch chan<- prometheus.Metric,
	loadAvg []string,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var err error

	if len(loadAvg) == 3 {
		if loadAvg[0] != "" {
			loadAvg01, err := strconv.ParseFloat(loadAvg[0], 64)
			if err != nil {
				err = errors.New(fmt.Sprintf("Error while converting Load avg01 metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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
				err = errors.New(fmt.Sprintf("Error while converting Load avg05 metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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
				err = errors.New(fmt.Sprintf("Error while converting Load avg15 metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	return err
}

func (c *JobsCollector) jobCPUMetrics(
	ch chan<- prometheus.Metric,
	cpu deployments.CPU,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var err error

	if cpu.Sys != "" {
		cpuSys, err := strconv.ParseFloat(cpu.Sys, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting CPU Sys metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	if cpu.User != "" {
		cpuUser, err := strconv.ParseFloat(cpu.User, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting CPU User metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	if cpu.Wait != "" {
		cpuWait, err := strconv.ParseFloat(cpu.Wait, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting CPU Wait metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	return err
}

func (c *JobsCollector) jobMemMetrics(
	ch chan<- prometheus.Metric,
	mem deployments.Mem,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var err error

	if mem.KB != "" {
		memKB, err := strconv.ParseFloat(mem.KB, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting Mem KB metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	if mem.Percent != "" {
		memPercent, err := strconv.ParseFloat(mem.Percent, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting Mem Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	return err
}

func (c *JobsCollector) jobSwapMetrics(
	ch chan<- prometheus.Metric,
	swap deployments.Mem,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var err error

	if swap.KB != "" {
		swapKB, err := strconv.ParseFloat(swap.KB, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting Swap KB metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	if swap.Percent != "" {
		swapPercent, err := strconv.ParseFloat(swap.Percent, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting Swap Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	return err
}

func (c *JobsCollector) jobSystemDiskMetrics(
	ch chan<- prometheus.Metric,
	systemDisk deployments.Disk,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var err error

	if systemDisk.InodePercent != "" {
		systemDiskInodePercent, err := strconv.ParseFloat(systemDisk.InodePercent, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting System Disk Inode Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	if systemDisk.Percent != "" {
		systemDiskPercent, err := strconv.ParseFloat(systemDisk.Percent, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting System Disk Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	return err
}

func (c *JobsCollector) jobEphemeralDiskMetrics(
	ch chan<- prometheus.Metric,
	ephemeralDisk deployments.Disk,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var err error

	if ephemeralDisk.InodePercent != "" {
		ephemeralDiskInodePercent, err := strconv.ParseFloat(ephemeralDisk.InodePercent, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting Ephemeral Disk Inode Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	if ephemeralDisk.Percent != "" {
		ephemeralDiskPercent, err := strconv.ParseFloat(ephemeralDisk.Percent, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting Ephemeral Disk Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	return err
}

func (c *JobsCollector) jobPersistentDiskMetrics(
	ch chan<- prometheus.Metric,
	persistentDisk deployments.Disk,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var err error

	if persistentDisk.InodePercent != "" {
		persistentDiskInodePercent, err := strconv.ParseFloat(persistentDisk.InodePercent, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting Persistent Disk Inode Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	if persistentDisk.Percent != "" {
		persistentDiskPercent, err := strconv.ParseFloat(persistentDisk.Percent, 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error while converting Persistent Disk Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err))
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

	return err
}

func (c *JobsCollector) jobProcessHealthyMetrics(
	ch chan<- prometheus.Metric,
	healthy bool,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) error {
	var healthyMetric float64
	if healthy {
		healthyMetric = 1
	}

	ch <- prometheus.MustNewConstMetric(
		c.jobProcessHealthyDesc,
		prometheus.GaugeValue,
		healthyMetric,
		deploymentName,
		jobName,
		jobID,
		jobIndex,
		jobAZ,
		jobIP,
		jobProcessName,
	)

	return nil
}

func (c *JobsCollector) jobProcessUptimeMetrics(
	ch chan<- prometheus.Metric,
	uptime *uint64,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) error {
	if uptime != nil {
		ch <- prometheus.MustNewConstMetric(
			c.jobProcessUptimeDesc,
			prometheus.GaugeValue,
			float64(*uptime),
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		)
	}

	return nil
}

func (c *JobsCollector) jobProcessCPUMetrics(
	ch chan<- prometheus.Metric,
	cpu deployments.CPU,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) error {
	if cpu.Total != nil {
		ch <- prometheus.MustNewConstMetric(
			c.jobProcessCPUTotalDesc,
			prometheus.GaugeValue,
			float64(*cpu.Total),
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		)
	}

	return nil
}

func (c *JobsCollector) jobProcessMemMetrics(
	ch chan<- prometheus.Metric,
	mem deployments.MemInt,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) error {
	if mem.KB != nil {
		ch <- prometheus.MustNewConstMetric(
			c.jobProcessMemKBDesc,
			prometheus.GaugeValue,
			float64(*mem.KB),
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		)
	}

	if mem.Percent != nil {
		ch <- prometheus.MustNewConstMetric(
			c.jobProcessMemPercentDesc,
			prometheus.GaugeValue,
			*mem.Percent,
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		)
	}

	return nil
}
