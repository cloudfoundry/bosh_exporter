package collectors

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"

	"github.com/cloudfoundry/bosh_exporter/deployments"
	"github.com/cloudfoundry/bosh_exporter/filters"
)

const (
	boshDeploymentNameLabel     = model.MetaLabelPrefix + "bosh_deployment"
	boshDeploymentReleasesLabel = model.MetaLabelPrefix + "bosh_deployment_releases"
	boshJobProcessNameLabel     = model.MetaLabelPrefix + "bosh_job_process_name"
	boshJobProcessReleaseLabel  = model.MetaLabelPrefix + "bosh_job_process_release"
)

type LabelGroups map[LabelGroupKey]*LabelGroupValue

type LabelGroupKey struct {
	DeploymentName string
	ProcessName    string
}
type LabelGroupValue struct {
	Targets            []string
	ProcessRelease     string
	DeploymentReleases []string
}

func NewLabelGroupValue(deployment deployments.DeploymentInfo, process deployments.Process) *LabelGroupValue {
	lgv := &LabelGroupValue{}
	for _, release := range deployment.Releases {
		ri := release.ToString()
		lgv.DeploymentReleases = append(lgv.DeploymentReleases, ri)
		// warning: works only if release job name == process name
		if release.HasJobName(process.Name) {
			lgv.ProcessRelease = ri
		}
	}
	return lgv
}
func (labelGroupValue *LabelGroupValue) addTarget(ip string) {
	labelGroupValue.Targets = append(labelGroupValue.Targets, ip)
}

func (labelGroupValue *LabelGroupValue) exportReleasesAsString() string {
	var releases []string
	releases = append(releases, labelGroupValue.DeploymentReleases...)
	return strings.Join(releases, ",")
}

func (c *ServiceDiscoveryCollector) createLabels(key LabelGroupKey, value *LabelGroupValue) model.LabelSet {
	return model.LabelSet{
		boshDeploymentNameLabel:     model.LabelValue(key.DeploymentName),
		boshDeploymentReleasesLabel: model.LabelValue(value.exportReleasesAsString()),
		boshJobProcessNameLabel:     model.LabelValue(key.ProcessName),
		boshJobProcessReleaseLabel:  model.LabelValue(value.ProcessRelease),
	}
}

type TargetGroups []TargetGroup

type TargetGroup struct {
	Targets []string       `json:"targets"`
	Labels  model.LabelSet `json:"labels,omitempty"`
}

type ServiceDiscoveryCollector struct {
	serviceDiscoveryFilename                        string
	azsFilter                                       *filters.AZsFilter
	processesFilter                                 *filters.RegexpFilter
	cidrsFilter                                     *filters.CidrFilter
	lastServiceDiscoveryScrapeTimestampMetric       prometheus.Gauge
	lastServiceDiscoveryScrapeDurationSecondsMetric prometheus.Gauge
	mu                                              *sync.Mutex
}

func NewServiceDiscoveryCollector(
	namespace string,
	environment string,
	boshName string,
	boshUUID string,
	serviceDiscoveryFilename string,
	azsFilter *filters.AZsFilter,
	processesFilter *filters.RegexpFilter,
	cidrsFilter *filters.CidrFilter,
) *ServiceDiscoveryCollector {
	metrics := NewServiceDiscoveryCollectorMetrics(namespace, environment, boshName, boshUUID)
	collector := &ServiceDiscoveryCollector{
		serviceDiscoveryFilename: serviceDiscoveryFilename,
		azsFilter:                azsFilter,
		processesFilter:          processesFilter,
		cidrsFilter:              cidrsFilter,
		lastServiceDiscoveryScrapeTimestampMetric:       metrics.NewLastServiceDiscoveryScrapeTimestampMetric(),
		lastServiceDiscoveryScrapeDurationSecondsMetric: metrics.NewLastServiceDiscoveryScrapeDurationSecondsMetric(),
		mu: &sync.Mutex{},
	}
	return collector
}

func (c *ServiceDiscoveryCollector) Collect(deployments []deployments.DeploymentInfo, ch chan<- prometheus.Metric) error {
	var begun = time.Now()

	labelGroups := c.createLabelGroups(deployments)
	targetGroups := c.createTargetGroups(labelGroups)

	err := c.writeTargetGroupsToFile(targetGroups)

	c.lastServiceDiscoveryScrapeTimestampMetric.Set(float64(time.Now().Unix()))
	c.lastServiceDiscoveryScrapeTimestampMetric.Collect(ch)

	c.lastServiceDiscoveryScrapeDurationSecondsMetric.Set(time.Since(begun).Seconds())
	c.lastServiceDiscoveryScrapeDurationSecondsMetric.Collect(ch)

	return err
}

func (c *ServiceDiscoveryCollector) Describe(ch chan<- *prometheus.Desc) {
	c.lastServiceDiscoveryScrapeTimestampMetric.Describe(ch)
	c.lastServiceDiscoveryScrapeDurationSecondsMetric.Describe(ch)
}

func (c *ServiceDiscoveryCollector) getLabelGroupKey(
	deployment deployments.DeploymentInfo,
	process deployments.Process,
) LabelGroupKey {
	return LabelGroupKey{
		DeploymentName: deployment.Name,
		ProcessName:    process.Name,
	}
}

func (c *ServiceDiscoveryCollector) createLabelGroups(deployments []deployments.DeploymentInfo) LabelGroups {
	labelGroups := LabelGroups{}

	for _, deployment := range deployments {
		for _, instance := range deployment.Instances {
			ip, found := c.cidrsFilter.Select(instance.IPs)
			if !found || !c.azsFilter.Enabled(instance.AZ) {
				continue
			}

			for _, process := range instance.Processes {
				if !c.processesFilter.Enabled(process.Name) {
					continue
				}
				key := c.getLabelGroupKey(deployment, process)
				if _, found := labelGroups[key]; !found {
					labelGroups[key] = NewLabelGroupValue(deployment, process)
				}
				labelGroups[key].addTarget(ip)
			}
		}
	}

	return labelGroups
}

func (c *ServiceDiscoveryCollector) createTargetGroups(labelGroups LabelGroups) TargetGroups {
	targetGroups := TargetGroups{}

	for key, value := range labelGroups {
		targetGroups = append(targetGroups, TargetGroup{
			Labels:  c.createLabels(key, value),
			Targets: value.Targets,
		})
	}

	return targetGroups
}

func (c *ServiceDiscoveryCollector) writeTargetGroupsToFile(targetGroups TargetGroups) error {
	targetGroupsJSON, err := json.Marshal(targetGroups)
	if err != nil {
		return fmt.Errorf("error while marshalling TargetGroups: %v", err)
	}

	dir, name := path.Split(c.serviceDiscoveryFilename)
	f, err := os.CreateTemp(dir, name)
	if err != nil {
		return fmt.Errorf("error creating temp file: %v", err)
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
		_ = os.Remove(f.Name())
	}

	return err
}
