package collectors

import (
	"strconv"
	"sync"
	"time"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type JobsCollector struct {
	namespace                         string
	boshDeployments                   []string
	boshClient                        director.Director
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
	lastJobsScrapeDurationSecondsDesc *prometheus.Desc
}

func NewJobsCollector(
	namespace string,
	boshDeployments []string,
	boshClient director.Director,
) *JobsCollector {
	jobHealthyDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_healthy"),
		"BOSH Job Healthy.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobLoadAvg01Desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_load_avg01"),
		"BOSH Job Load avg01.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobLoadAvg05Desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_load_avg05"),
		"BOSH Job Load avg05.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobLoadAvg15Desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_load_avg15"),
		"BOSH Job Load avg15.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobCPUSysDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_cpu_sys"),
		"BOSH Job CPU System.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobCPUUserDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_cpu_user"),
		"BOSH Job CPU User.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobCPUWaitDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_cpu_wait"),
		"BOSH Job CPU Wait.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobMemKBDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_mem_kb"),
		"BOSH Job Memory KB.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobMemPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_mem_percent"),
		"BOSH Job Memory Percent.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobSwapKBDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_swap_kb"),
		"BOSH Job Swap KB.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobSwapPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_swap_percent"),
		"BOSH Job Swap Percent.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobSystemDiskInodePercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_system_disk_inode_percent"),
		"BOSH Job System Disk Inode Percent.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobSystemDiskPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_system_disk_percent"),
		"BOSH Job System Disk Percent.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobEphemeralDiskInodePercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_ephemeral_disk_inode_percent"),
		"BOSH Job Ephemeral Disk Inode Percent.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobEphemeralDiskPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_ephemeral_disk_percent"),
		"BOSH Job Ephemeral Disk Percent.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobPersistentDiskInodePercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_persistent_disk_inode_percent"),
		"BOSH Job Persistent Disk Inode Percent.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
		nil,
	)

	jobPersistentDiskPercentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "bosh", "job_persistent_disk_percent"),
		"BOSH Job Persistent Disk Percent.",
		[]string{"bosh_deployment", "bosh_job", "bosh_index", "bosh_az", "bosh_ip"},
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
		boshDeployments:                   boshDeployments,
		boshClient:                        boshClient,
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
		lastJobsScrapeDurationSecondsDesc: lastJobsScrapeDurationSecondsDesc,
	}
	return collector
}

func (c JobsCollector) Collect(ch chan<- prometheus.Metric) {
	var err error
	var deployments []director.Deployment
	var begun = time.Now()

	if len(c.boshDeployments) > 0 {
		log.Debugf("Filtering deployments by `%v`...", c.boshDeployments)
		for _, deploymentName := range c.boshDeployments {
			deployment, err := c.boshClient.FindDeployment(deploymentName)
			if err != nil {
				log.Errorf("Error while reading deployment `%s`: %v", deploymentName, err)
				continue
			}
			deployments = append(deployments, deployment)
		}
	} else {
		log.Debugf("Reading deployments...")
		deployments, err = c.boshClient.Deployments()
		if err != nil {
			log.Errorf("Error while reading deployments: %v", err)
			return
		}
	}

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
	ch <- c.lastJobsScrapeDurationSecondsDesc
}

func (c JobsCollector) reportJobMetrics(
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

		c.jobHealthyMetrics(ch, vmInfo.IsRunning(), deploymentName, jobName, jobIndex, jobAZ, jobIP)
		c.jobLoadAvgMetrics(ch, vmInfo.Vitals.Load, deploymentName, jobName, jobIndex, jobAZ, jobIP)
		c.jobCPUMetrics(ch, vmInfo.Vitals.CPU, deploymentName, jobName, jobIndex, jobAZ, jobIP)
		c.jobMemMetrics(ch, vmInfo.Vitals.Mem, deploymentName, jobName, jobIndex, jobAZ, jobIP)
		c.jobSwapMetrics(ch, vmInfo.Vitals.Swap, deploymentName, jobName, jobIndex, jobAZ, jobIP)
		c.jobSystemDiskMetrics(ch, vmInfo.Vitals.SystemDisk(), deploymentName, jobName, jobIndex, jobAZ, jobIP)
		c.jobEphemeralDiskMetrics(ch, vmInfo.Vitals.EphemeralDisk(), deploymentName, jobName, jobIndex, jobAZ, jobIP)
		c.jobPersistentDiskMetrics(ch, vmInfo.Vitals.PersistentDisk(), deploymentName, jobName, jobIndex, jobAZ, jobIP)
	}
}

func (c JobsCollector) jobHealthyMetrics(
	ch chan<- prometheus.Metric,
	vmRunning bool,
	deploymentName string,
	jobName string,
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
				jobIndex,
				jobAZ,
				jobIP,
			)
		}
	}
}
