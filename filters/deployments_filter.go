package filters

import (
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/prometheus/common/log"
)

type DeploymentsFilter struct {
	filter     []string
	boshClient director.Director
}

func NewDeploymentsFilter(
	filter []string,
	boshClient director.Director,
) *DeploymentsFilter {
	return &DeploymentsFilter{
		filter:     filter,
		boshClient: boshClient,
	}
}

func (f DeploymentsFilter) GetDeployments() []director.Deployment {
	var err error
	var deployments []director.Deployment

	if len(f.filter) > 0 {
		log.Debugf("Filtering deployments by `%v`...", f.filter)
		for _, deploymentName := range f.filter {
			deployment, err := f.boshClient.FindDeployment(deploymentName)
			if err != nil {
				log.Errorf("Error while reading deployment `%s`: %v", deploymentName, err)
				continue
			}
			deployments = append(deployments, deployment)
		}
	} else {
		log.Debugf("Reading deployments...")
		deployments, err = f.boshClient.Deployments()
		if err != nil {
			log.Errorf("Error while reading deployments: %v", err)
			return deployments
		}
	}

	return deployments
}
