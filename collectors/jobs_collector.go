package collectors

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bosh-prometheus/bosh_exporter/deployments"
	"github.com/bosh-prometheus/bosh_exporter/filters"
)

type JobsCollector struct {
	azsFilter                           *filters.AZsFilter
	cidrsFilter                         *filters.CidrFilter
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
}

func NewJobsCollector(
	namespace string,
	environment string,
	boshName string,
	boshUUID string,
	azsFilter *filters.AZsFilter,
	cidrsFilter *filters.CidrFilter,
) *JobsCollector {
	jobHealthyMetric := prometheus.NewGaugeVec(
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

	jobLoadAvg01Metric := prometheus.NewGaugeVec(
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

	jobLoadAvg05Metric := prometheus.NewGaugeVec(
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

	jobLoadAvg15Metric := prometheus.NewGaugeVec(
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

	jobCPUSysMetric := prometheus.NewGaugeVec(
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

	jobCPUUserMetric := prometheus.NewGaugeVec(
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

	jobCPUWaitMetric := prometheus.NewGaugeVec(
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

	jobMemKBMetric := prometheus.NewGaugeVec(
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

	jobMemPercentMetric := prometheus.NewGaugeVec(
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

	jobSwapKBMetric := prometheus.NewGaugeVec(
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

	jobSwapPercentMetric := prometheus.NewGaugeVec(
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

	jobSystemDiskInodePercentMetric := prometheus.NewGaugeVec(
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

	jobSystemDiskPercentMetric := prometheus.NewGaugeVec(
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

	jobEphemeralDiskInodePercentMetric := prometheus.NewGaugeVec(
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

	jobEphemeralDiskPercentMetric := prometheus.NewGaugeVec(
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

	jobPersistentDiskInodePercentMetric := prometheus.NewGaugeVec(
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

	jobPersistentDiskPercentMetric := prometheus.NewGaugeVec(
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

	jobProcessHealthyMetric := prometheus.NewGaugeVec(
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

	jobProcessUptimeMetric := prometheus.NewGaugeVec(
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

	jobProcessCPUTotalMetric := prometheus.NewGaugeVec(
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

	jobProcessMemKBMetric := prometheus.NewGaugeVec(
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

	jobProcessMemPercentMetric := prometheus.NewGaugeVec(
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

	lastJobsScrapeTimestampMetric := prometheus.NewGauge(
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

	lastJobsScrapeDurationSecondsMetric := prometheus.NewGauge(
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

	collector := &JobsCollector{
		azsFilter:                           azsFilter,
		cidrsFilter:                         cidrsFilter,
		jobHealthyMetric:                    jobHealthyMetric,
		jobLoadAvg01Metric:                  jobLoadAvg01Metric,
		jobLoadAvg05Metric:                  jobLoadAvg05Metric,
		jobLoadAvg15Metric:                  jobLoadAvg15Metric,
		jobCPUSysMetric:                     jobCPUSysMetric,
		jobCPUUserMetric:                    jobCPUUserMetric,
		jobCPUWaitMetric:                    jobCPUWaitMetric,
		jobMemKBMetric:                      jobMemKBMetric,
		jobMemPercentMetric:                 jobMemPercentMetric,
		jobSwapKBMetric:                     jobSwapKBMetric,
		jobSwapPercentMetric:                jobSwapPercentMetric,
		jobSystemDiskInodePercentMetric:     jobSystemDiskInodePercentMetric,
		jobSystemDiskPercentMetric:          jobSystemDiskPercentMetric,
		jobEphemeralDiskInodePercentMetric:  jobEphemeralDiskInodePercentMetric,
		jobEphemeralDiskPercentMetric:       jobEphemeralDiskPercentMetric,
		jobPersistentDiskInodePercentMetric: jobPersistentDiskInodePercentMetric,
		jobPersistentDiskPercentMetric:      jobPersistentDiskPercentMetric,
		jobProcessHealthyMetric:             jobProcessHealthyMetric,
		jobProcessUptimeMetric:              jobProcessUptimeMetric,
		jobProcessCPUTotalMetric:            jobProcessCPUTotalMetric,
		jobProcessMemKBMetric:               jobProcessMemKBMetric,
		jobProcessMemPercentMetric:          jobProcessMemPercentMetric,
		lastJobsScrapeTimestampMetric:       lastJobsScrapeTimestampMetric,
		lastJobsScrapeDurationSecondsMetric: lastJobsScrapeDurationSecondsMetric,
	}
	return collector
}

func (c *JobsCollector) Collect(deployments []deployments.DeploymentInfo, ch chan<- prometheus.Metric) error {
	var err error
	var begun = time.Now()

	c.jobHealthyMetric.Reset()
	c.jobLoadAvg01Metric.Reset()
	c.jobLoadAvg05Metric.Reset()
	c.jobLoadAvg15Metric.Reset()
	c.jobCPUSysMetric.Reset()
	c.jobCPUUserMetric.Reset()
	c.jobCPUWaitMetric.Reset()
	c.jobMemKBMetric.Reset()
	c.jobMemPercentMetric.Reset()
	c.jobSwapKBMetric.Reset()
	c.jobSwapPercentMetric.Reset()
	c.jobSystemDiskInodePercentMetric.Reset()
	c.jobSystemDiskPercentMetric.Reset()
	c.jobEphemeralDiskInodePercentMetric.Reset()
	c.jobEphemeralDiskPercentMetric.Reset()
	c.jobPersistentDiskInodePercentMetric.Reset()
	c.jobPersistentDiskPercentMetric.Reset()
	c.jobProcessHealthyMetric.Reset()
	c.jobProcessUptimeMetric.Reset()
	c.jobProcessCPUTotalMetric.Reset()
	c.jobProcessMemKBMetric.Reset()
	c.jobProcessMemPercentMetric.Reset()

	for _, deployment := range deployments {
		err = c.reportJobMetrics(deployment)
	}

	c.jobHealthyMetric.Collect(ch)
	c.jobLoadAvg01Metric.Collect(ch)
	c.jobLoadAvg05Metric.Collect(ch)
	c.jobLoadAvg15Metric.Collect(ch)
	c.jobCPUSysMetric.Collect(ch)
	c.jobCPUUserMetric.Collect(ch)
	c.jobCPUWaitMetric.Collect(ch)
	c.jobMemKBMetric.Collect(ch)
	c.jobMemPercentMetric.Collect(ch)
	c.jobSwapKBMetric.Collect(ch)
	c.jobSwapPercentMetric.Collect(ch)
	c.jobSystemDiskInodePercentMetric.Collect(ch)
	c.jobSystemDiskPercentMetric.Collect(ch)
	c.jobEphemeralDiskInodePercentMetric.Collect(ch)
	c.jobEphemeralDiskPercentMetric.Collect(ch)
	c.jobPersistentDiskInodePercentMetric.Collect(ch)
	c.jobPersistentDiskPercentMetric.Collect(ch)
	c.jobProcessHealthyMetric.Collect(ch)
	c.jobProcessUptimeMetric.Collect(ch)
	c.jobProcessCPUTotalMetric.Collect(ch)
	c.jobProcessMemKBMetric.Collect(ch)
	c.jobProcessMemPercentMetric.Collect(ch)

	c.lastJobsScrapeTimestampMetric.Set(float64(time.Now().Unix()))
	c.lastJobsScrapeTimestampMetric.Collect(ch)

	c.lastJobsScrapeDurationSecondsMetric.Set(time.Since(begun).Seconds())
	c.lastJobsScrapeDurationSecondsMetric.Collect(ch)

	return err
}

func (c *JobsCollector) Describe(ch chan<- *prometheus.Desc) {
	c.jobHealthyMetric.Describe(ch)
	c.jobLoadAvg01Metric.Describe(ch)
	c.jobLoadAvg05Metric.Describe(ch)
	c.jobLoadAvg15Metric.Describe(ch)
	c.jobCPUSysMetric.Describe(ch)
	c.jobCPUUserMetric.Describe(ch)
	c.jobCPUWaitMetric.Describe(ch)
	c.jobMemKBMetric.Describe(ch)
	c.jobMemPercentMetric.Describe(ch)
	c.jobSwapKBMetric.Describe(ch)
	c.jobSwapPercentMetric.Describe(ch)
	c.jobSystemDiskInodePercentMetric.Describe(ch)
	c.jobSystemDiskPercentMetric.Describe(ch)
	c.jobEphemeralDiskInodePercentMetric.Describe(ch)
	c.jobEphemeralDiskPercentMetric.Describe(ch)
	c.jobPersistentDiskInodePercentMetric.Describe(ch)
	c.jobPersistentDiskPercentMetric.Describe(ch)
	c.jobProcessHealthyMetric.Describe(ch)
	c.jobProcessUptimeMetric.Describe(ch)
	c.jobProcessCPUTotalMetric.Describe(ch)
	c.jobProcessMemKBMetric.Describe(ch)
	c.jobProcessMemPercentMetric.Describe(ch)
	c.lastJobsScrapeTimestampMetric.Describe(ch)
	c.lastJobsScrapeDurationSecondsMetric.Describe(ch)
}

func (c *JobsCollector) reportJobMetrics(deployment deployments.DeploymentInfo) error {
	var endErr error

	for _, instance := range deployment.Instances {
		if !c.azsFilter.Enabled(instance.AZ) {
			continue
		}

		deploymentName := deployment.Name
		jobName := instance.Name
		jobID := instance.ID
		jobIndex := instance.Index
		jobAZ := instance.AZ
		jobIP, _ := c.cidrsFilter.Select(instance.IPs)

		c.jobHealthyMetrics(instance.Healthy, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)

		err := c.jobLoadAvgMetrics(instance.Vitals.Load, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		if err != nil {
			endErr = err
		}

		err = c.jobCPUMetrics(instance.Vitals.CPU, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		if err != nil {
			endErr = err
		}

		err = c.jobMemMetrics(instance.Vitals.Mem, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		if err != nil {
			endErr = err
		}

		err = c.jobSwapMetrics(instance.Vitals.Swap, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		if err != nil {
			endErr = err
		}

		err = c.jobSystemDiskMetrics(instance.Vitals.SystemDisk, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		if err != nil {
			endErr = err
		}

		err = c.jobEphemeralDiskMetrics(instance.Vitals.EphemeralDisk, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		if err != nil {
			endErr = err
		}

		err = c.jobPersistentDiskMetrics(instance.Vitals.PersistentDisk, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP)
		if err != nil {
			endErr = err
		}

		for _, process := range instance.Processes {
			jobProcessName := process.Name
			c.jobProcessHealthyMetrics(process.Healthy, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
			c.jobProcessUptimeMetrics(process.Uptime, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
			c.jobProcessCPUMetrics(process.CPU, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
			c.jobProcessMemMetrics(process.Mem, deploymentName, jobName, jobID, jobIndex, jobAZ, jobIP, jobProcessName)
		}
	}

	return endErr
}

func (c *JobsCollector) jobHealthyMetrics(
	healthy bool,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) {
	var healthyMetric float64
	if healthy {
		healthyMetric = 1
	}

	c.jobHealthyMetric.WithLabelValues(
		deploymentName,
		jobName,
		jobID,
		jobIndex,
		jobAZ,
		jobIP,
	).Set(healthyMetric)
}

func (c *JobsCollector) jobLoadAvgMetrics(
	loadAvg []string,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var (
		err  error
		load float64
	)

	if len(loadAvg) == 3 {
		if loadAvg[0] != "" {
			load, err = strconv.ParseFloat(loadAvg[0], 64)
			if err != nil {
				err = fmt.Errorf("error while converting Load avg01 metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
			} else {
				c.jobLoadAvg01Metric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				).Set(load)
			}
		}

		if loadAvg[1] != "" {
			load, err = strconv.ParseFloat(loadAvg[1], 64)
			if err != nil {
				err = fmt.Errorf("error while converting Load avg05 metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
			} else {
				c.jobLoadAvg05Metric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				).Set(load)
			}
		}

		if loadAvg[2] != "" {
			load, err = strconv.ParseFloat(loadAvg[2], 64)
			if err != nil {
				err = fmt.Errorf("error while converting Load avg15 metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
			} else {
				c.jobLoadAvg15Metric.WithLabelValues(
					deploymentName,
					jobName,
					jobID,
					jobIndex,
					jobAZ,
					jobIP,
				).Set(load)
			}
		}
	}

	return err
}

func (c *JobsCollector) jobCPUMetrics(
	cpu deployments.CPU,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var (
		err  error
		load float64
	)

	if cpu.Sys != "" {
		load, err = strconv.ParseFloat(cpu.Sys, 64)
		if err != nil {
			err = fmt.Errorf("error while converting CPU Sys metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobCPUSysMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(load)
		}
	}

	if cpu.User != "" {
		load, err = strconv.ParseFloat(cpu.User, 64)
		if err != nil {
			err = fmt.Errorf("error while converting CPU User metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobCPUUserMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(load)
		}
	}

	if cpu.Wait != "" {
		load, err = strconv.ParseFloat(cpu.Wait, 64)
		if err != nil {
			err = fmt.Errorf("error while converting CPU Wait metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobCPUWaitMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(load)
		}
	}

	return err
}

func (c *JobsCollector) jobMemMetrics(
	mem deployments.Mem,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var (
		err   error
		value float64
	)

	if mem.KB != "" {
		value, err = strconv.ParseFloat(mem.KB, 64)
		if err != nil {
			err = fmt.Errorf("error while converting Mem KB metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobMemKBMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(value)
		}
	}

	if mem.Percent != "" {
		value, err = strconv.ParseFloat(mem.Percent, 64)
		if err != nil {
			err = fmt.Errorf("error while converting Mem Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobMemPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(value)
		}
	}

	return err
}

func (c *JobsCollector) jobSwapMetrics(
	swap deployments.Mem,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var (
		err   error
		value float64
	)

	if swap.KB != "" {
		value, err = strconv.ParseFloat(swap.KB, 64)
		if err != nil {
			err = fmt.Errorf("error while converting Swap KB metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobSwapKBMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(value)
		}
	}

	if swap.Percent != "" {
		value, err = strconv.ParseFloat(swap.Percent, 64)
		if err != nil {
			err = fmt.Errorf("error while converting Swap Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobSwapPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(value)
		}
	}

	return err
}

func (c *JobsCollector) jobSystemDiskMetrics(
	systemDisk deployments.Disk,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var (
		err   error
		value float64
	)

	if systemDisk.InodePercent != "" {
		value, err = strconv.ParseFloat(systemDisk.InodePercent, 64)
		if err != nil {
			err = fmt.Errorf("error while converting System Disk Inode Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobSystemDiskInodePercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(value)
		}
	}

	if systemDisk.Percent != "" {
		value, err = strconv.ParseFloat(systemDisk.Percent, 64)
		if err != nil {
			err = fmt.Errorf("error while converting System Disk Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobSystemDiskPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(value)
		}
	}

	return err
}

func (c *JobsCollector) jobEphemeralDiskMetrics(
	ephemeralDisk deployments.Disk,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var (
		err   error
		value float64
	)

	if ephemeralDisk.InodePercent != "" {
		value, err = strconv.ParseFloat(ephemeralDisk.InodePercent, 64)
		if err != nil {
			err = fmt.Errorf("error while converting Ephemeral Disk Inode Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobEphemeralDiskInodePercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(value)
		}
	}

	if ephemeralDisk.Percent != "" {
		value, err = strconv.ParseFloat(ephemeralDisk.Percent, 64)
		if err != nil {
			err = fmt.Errorf("error while converting Ephemeral Disk Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobEphemeralDiskPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(value)
		}
	}

	return err
}

func (c *JobsCollector) jobPersistentDiskMetrics(
	persistentDisk deployments.Disk,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
) error {
	var (
		err   error
		value float64
	)

	if persistentDisk.InodePercent != "" {
		value, err = strconv.ParseFloat(persistentDisk.InodePercent, 64)
		if err != nil {
			err = fmt.Errorf("error while converting Persistent Disk Inode Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobPersistentDiskInodePercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(value)
		}
	}

	if persistentDisk.Percent != "" {
		value, err = strconv.ParseFloat(persistentDisk.Percent, 64)
		if err != nil {
			err = fmt.Errorf("error while converting Persistent Disk Percent metric for deployment `%s` and job `%s`: %v", deploymentName, jobName, err)
		} else {
			c.jobPersistentDiskPercentMetric.WithLabelValues(
				deploymentName,
				jobName,
				jobID,
				jobIndex,
				jobAZ,
				jobIP,
			).Set(value)
		}
	}

	return err
}

func (c *JobsCollector) jobProcessHealthyMetrics(
	healthy bool,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) {
	var healthyMetric float64
	if healthy {
		healthyMetric = 1
	}

	c.jobProcessHealthyMetric.WithLabelValues(
		deploymentName,
		jobName,
		jobID,
		jobIndex,
		jobAZ,
		jobIP,
		jobProcessName,
	).Set(healthyMetric)
}

func (c *JobsCollector) jobProcessUptimeMetrics(
	uptime *uint64,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) {
	if uptime != nil {
		c.jobProcessUptimeMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		).Set(float64(*uptime))
	}
}

func (c *JobsCollector) jobProcessCPUMetrics(
	cpu deployments.CPU,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) {
	if cpu.Total != nil {
		c.jobProcessCPUTotalMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		).Set(*cpu.Total)
	}
}

func (c *JobsCollector) jobProcessMemMetrics(
	mem deployments.MemInt,
	deploymentName string,
	jobName string,
	jobID string,
	jobIndex string,
	jobAZ string,
	jobIP string,
	jobProcessName string,
) {
	if mem.KB != nil {
		c.jobProcessMemKBMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		).Set(float64(*mem.KB))
	}

	if mem.Percent != nil {
		c.jobProcessMemPercentMetric.WithLabelValues(
			deploymentName,
			jobName,
			jobID,
			jobIndex,
			jobAZ,
			jobIP,
			jobProcessName,
		).Set(*mem.Percent)
	}
}
