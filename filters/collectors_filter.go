package filters

import (
	"errors"
	"fmt"
)

const (
	DeploymentsCollector = "Deployments"
	JobsCollector        = "Jobs"
)

type CollectorsFilter struct {
	collectorsEnabled map[string]bool
}

func NewCollectorsFilter(filter []string) (*CollectorsFilter, error) {
	collectorsEnabled := make(map[string]bool)

	for _, collectorName := range filter {
		switch collectorName {
		case DeploymentsCollector:
			collectorsEnabled[DeploymentsCollector] = true
		case JobsCollector:
			collectorsEnabled[JobsCollector] = true
		default:
			return &CollectorsFilter{}, errors.New(fmt.Sprintf("Collector filter `%s` is not supported", collectorName))
		}
	}

	return &CollectorsFilter{collectorsEnabled: collectorsEnabled}, nil
}

func (f *CollectorsFilter) Enabled(collectorName string) bool {
	if len(f.collectorsEnabled) > 0 {
		if f.collectorsEnabled[collectorName] {
			return true
		}

		return false
	}

	return true
}
