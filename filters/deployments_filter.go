package filters

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/prometheus/common/log"
)

type DeploymentsFilter struct {
	filters    []string
	boshClient director.Director
}

func NewDeploymentsFilter(filters []string, boshClient director.Director) *DeploymentsFilter {
	return &DeploymentsFilter{filters: filters, boshClient: boshClient}
}

func (f *DeploymentsFilter) GetDeployments() ([]director.Deployment, error) {
	var err error
	var deployments []director.Deployment

	if len(f.filters) > 0 {
		log.Debugf("Filtering deployments by `%v`...", f.filters)
		for _, deploymentName := range f.filters {
			deployment, err := f.boshClient.FindDeployment(strings.Trim(deploymentName, " "))
			if err != nil {
				return deployments, fmt.Errorf("error while reading deployment `%s`: %v", deploymentName, err)
			}
			deployments = append(deployments, deployment)
		}
	} else {
		log.Debugf("Reading deployments...")
		deployments, err = f.boshClient.Deployments()
		if err != nil {
			return deployments, fmt.Errorf("error while reading deployments: %v", err)
		}
	}

	return deployments, nil
}
