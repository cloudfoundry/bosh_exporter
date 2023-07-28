package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

type JobsCollectorMetrics struct {
	namespace   string
	environment string
	boshName    string
	boshUUID    string
}

func NewJobsCollectorMetrics(
	namespace string,
	environment string,
	boshName string,
	boshUUID string,
) *JobsCollectorMetrics {
	return &JobsCollectorMetrics{
		namespace:   namespace,
		environment: environment,
		boshName:    boshName,
		boshUUID:    boshUUID,
	}
}

func (m *JobsCollectorMetrics) NewLastJobsScrapeDurationSecondsMetric() prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "",
			Name:      "last_jobs_scrape_duration_seconds",
			Help:      "Duration of the last scrape of Job metrics from BOSH.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
	)
}

func (m *JobsCollectorMetrics) NewLastJobsScrapeTimestampMetric() prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "",
			Name:      "last_jobs_scrape_timestamp",
			Help:      "Number of seconds since 1970 since last scrape of Job metrics from BOSH.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
	)
}

func (m *JobsCollectorMetrics) NewJobProcessMemPercentMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job_process",
			Name:      "mem_percent",
			Help:      "BOSH Job Process Memory Percent.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
	)
}

func (m *JobsCollectorMetrics) NewJobProcessMemKBMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job_process",
			Name:      "mem_kb",
			Help:      "BOSH Job Process Memory KB.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
	)
}

func (m *JobsCollectorMetrics) NewJobProcessCPUTotalMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job_process",
			Name:      "cpu_total",
			Help:      "BOSH Job Process CPU Total.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
	)
}

func (m *JobsCollectorMetrics) NewJobProcessUptimeMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job_process",
			Name:      "uptime_seconds",
			Help:      "BOSH Job Process Uptime in seconds.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
	)
}

func (m *JobsCollectorMetrics) NewJobProcessHealthyMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job_process",
			Name:      "healthy",
			Help:      "BOSH Job Process Healthy (1 for healthy, 0 for unhealthy).",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name"},
	)
}

func (m *JobsCollectorMetrics) NewJobProcessInfoMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job_process",
			Name:      "info",
			Help:      "BOSH Job Process Info with a constant '1' value. Release can be found only if process name is the same as release job name.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip", "bosh_job_process_name", "bosh_job_process_release_name", "bosh_job_process_release_version"},
	)
}

func (m *JobsCollectorMetrics) NewJobPersistentDiskPercentMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "persistent_disk_percent",
			Help:      "BOSH Job Persistent Disk Percent.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobPersistentDiskInodePercentMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "persistent_disk_inode_percent",
			Help:      "BOSH Job Persistent Disk Inode Percent.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobEphemeralDiskPercentMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "ephemeral_disk_percent",
			Help:      "BOSH Job Ephemeral Disk Percent.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobEphemeralDiskInodePercentMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "ephemeral_disk_inode_percent",
			Help:      "BOSH Job Ephemeral Disk Inode Percent.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobSystemDiskPercentMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "system_disk_percent",
			Help:      "BOSH Job System Disk Percent.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobSystemDiskInodePercentMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "system_disk_inode_percent",
			Help:      "BOSH Job System Disk Inode Percent.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobSwapPercentMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "swap_percent",
			Help:      "BOSH Job Swap Percent.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobSwapKBMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "swap_kb",
			Help:      "BOSH Job Swap KB.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobMemPercentMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "mem_percent",
			Help:      "BOSH Job Memory Percent.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobMemKBMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "mem_kb",
			Help:      "BOSH Job Memory KB.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobCPUWaitMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "cpu_wait",
			Help:      "BOSH Job CPU Wait.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobCPUUserMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "cpu_user",
			Help:      "BOSH Job CPU User.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobCPUSysMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "cpu_sys",
			Help:      "BOSH Job CPU System.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobLoadAvg15Metric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "load_avg15",
			Help:      "BOSH Job Load avg15.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobLoadAvg05Metric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "load_avg05",
			Help:      "BOSH Job Load avg05.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobLoadAvg01Metric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "load_avg01",
			Help:      "BOSH Job Load avg01.",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}

func (m *JobsCollectorMetrics) NewJobHealthyMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: "job",
			Name:      "healthy",
			Help:      "BOSH Job Healthy (1 for healthy, 0 for unhealthy).",
			ConstLabels: prometheus.Labels{
				"environment": m.environment,
				"bosh_name":   m.boshName,
				"bosh_uuid":   m.boshUUID,
			},
		},
		[]string{"bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_index", "bosh_job_az", "bosh_job_ip"},
	)
}
