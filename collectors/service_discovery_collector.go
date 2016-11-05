package collectors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/model"

	"github.com/cloudfoundry-community/bosh_exporter/filters"
)

const (
	boshProcessNameLabel = model.MetaLabelPrefix + "bosh_process"
)

type Processes map[string][]ProcessInfo

type ProcessInfo struct {
	DeploymentName string
	JobName        string
	JobIndex       int
	JobAZ          string
	JobIP          string
}

type TargetGroups []TargetGroup

type TargetGroup struct {
	Targets []string       `json:"targets"`
	Labels  model.LabelSet `json:"labels,omitempty"`
}

type ServiceDiscoveryCollector struct {
	namespace                                     string
	deploymentsFilter                             filters.DeploymentsFilter
	serviceDiscoveryFilename                      string
	processesFilter                               filters.RegexpFilter
	lastServiceDiscoveryScrapeTimestampDesc       *prometheus.Desc
	lastServiceDiscoveryScrapeDurationSecondsDesc *prometheus.Desc
	mu                                            *sync.Mutex
}

func NewServiceDiscoveryCollector(
	namespace string,
	deploymentsFilter filters.DeploymentsFilter,
	serviceDiscoveryFilename string,
	processesFilter filters.RegexpFilter,
) *ServiceDiscoveryCollector {
	lastServiceDiscoveryScrapeTimestampDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "last_service_discovery_scrape_timestamp"),
		"Number of seconds since 1970 since last scrape of Service Discovery from BOSH.",
		[]string{},
		nil,
	)

	lastServiceDiscoveryScrapeDurationSecondsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "last_service_discovery_scrape_duration_seconds"),
		"Duration of the last scrape of Service Discovery from BOSH.",
		[]string{},
		nil,
	)

	collector := &ServiceDiscoveryCollector{
		namespace:                                     namespace,
		deploymentsFilter:                             deploymentsFilter,
		serviceDiscoveryFilename:                      serviceDiscoveryFilename,
		processesFilter:                               processesFilter,
		lastServiceDiscoveryScrapeTimestampDesc:       lastServiceDiscoveryScrapeTimestampDesc,
		lastServiceDiscoveryScrapeDurationSecondsDesc: lastServiceDiscoveryScrapeDurationSecondsDesc,
		mu: &sync.Mutex{},
	}
	return collector
}

func (c ServiceDiscoveryCollector) Collect(ch chan<- prometheus.Metric) {
	var begun = time.Now()

	deployments := c.deploymentsFilter.GetDeployments()
	processes := make(Processes)

	var wg sync.WaitGroup
	for _, deployment := range deployments {
		wg.Add(1)
		go func(deployment director.Deployment, processes Processes) {
			defer wg.Done()
			c.getDeploymentProcesses(deployment, processes)
		}(deployment, processes)
	}
	wg.Wait()

	targetGroups := c.createTargetGroups(processes)

	if err := c.writeTargetGroupsToFile(targetGroups); err != nil {
		log.Error(err)
	}

	ch <- prometheus.MustNewConstMetric(
		c.lastServiceDiscoveryScrapeTimestampDesc,
		prometheus.GaugeValue,
		float64(time.Now().Unix()),
	)

	ch <- prometheus.MustNewConstMetric(
		c.lastServiceDiscoveryScrapeDurationSecondsDesc,
		prometheus.GaugeValue,
		time.Since(begun).Seconds(),
	)
}

func (c ServiceDiscoveryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.lastServiceDiscoveryScrapeTimestampDesc
	ch <- c.lastServiceDiscoveryScrapeDurationSecondsDesc
}

func (c ServiceDiscoveryCollector) getDeploymentProcesses(deployment director.Deployment, processes Processes) {
	log.Debugf("Reading VM info for deployment `%s`:", deployment.Name())
	vmInfos, err := deployment.VMInfos()
	if err != nil {
		log.Errorf("Error while reading VM info for deployment `%s`: %v", deployment.Name(), err)
		return
	}

	for _, vmInfo := range vmInfos {
		if len(vmInfo.IPs) >= 0 {
			for _, pi := range vmInfo.Processes {
				if !c.processesFilter.Enabled(pi.Name) {
					continue
				}

				processInfo := &ProcessInfo{
					DeploymentName: deployment.Name(),
					JobName:        vmInfo.JobName,
					JobAZ:          vmInfo.AZ,
					JobIP:          vmInfo.IPs[0],
				}

				c.mu.Lock()
				processes[pi.Name] = append(processes[pi.Name], *processInfo)
				c.mu.Unlock()
			}
		}
	}
}

func (c ServiceDiscoveryCollector) createTargetGroups(processes Processes) TargetGroups {
	targetGroups := TargetGroups{}

	for processName, processesInfo := range processes {
		targets := []string{}
		for _, processInfo := range processesInfo {
			targets = append(targets, processInfo.JobIP)
		}

		targetGroup := TargetGroup{
			Targets: targets,
			Labels: model.LabelSet{
				model.LabelName(boshProcessNameLabel): model.LabelValue(processName),
			},
		}
		targetGroups = append(targetGroups, targetGroup)
	}

	return targetGroups
}

func (c ServiceDiscoveryCollector) writeTargetGroupsToFile(targetGroups TargetGroups) error {
	targetGroupsJSON, err := json.Marshal(targetGroups)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while marshalling TargetGroups: %v", err))
	}

	dir, name := path.Split(c.serviceDiscoveryFilename)
	f, err := ioutil.TempFile(dir, name)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating temp file: %v", err))
	}

	_, err = f.Write(targetGroupsJSON)
	if err == nil {
		err = f.Sync()
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if permErr := os.Chmod(f.Name(), 0644); err == nil {
		err = permErr
	}
	if err == nil {
		err = os.Rename(f.Name(), c.serviceDiscoveryFilename)
	}

	if err != nil {
		os.Remove(f.Name())
	}

	return err
}
