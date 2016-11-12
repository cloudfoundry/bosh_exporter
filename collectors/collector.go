package collectors

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudfoundry-community/bosh_exporter/deployments"
)

type Collector interface {
	Collect(deployments []deployments.DeploymentInfo, ch chan<- prometheus.Metric) error
	Describe(ch chan<- *prometheus.Desc)
}
