package deployments_test

import (
	"errors"
	"strconv"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/director/directorfakes"
	"github.com/cppforlife/go-semi-semantic/version"
	"github.com/prometheus/common/log"

	"github.com/cloudfoundry/bosh_exporter/filters"

	"github.com/cloudfoundry/bosh_exporter/deployments"
)

func init() {
	_ = log.Base().SetLevel("fatal")
}

var _ = ginkgo.Describe("Fetcher", func() {
	var (
		err                error
		boshDeployments    []string
		boshClient         *directorfakes.FakeDirector
		deploymentsFilter  *filters.DeploymentsFilter
		deploymentsFetcher *deployments.Fetcher
	)

	ginkgo.BeforeEach(func() {
		boshDeployments = []string{}
		boshClient = &directorfakes.FakeDirector{}
	})

	ginkgo.JustBeforeEach(func() {
		deploymentsFilter = filters.NewDeploymentsFilter(boshDeployments, boshClient)
		deploymentsFetcher = deployments.NewFetcher(*deploymentsFilter)
	})

	ginkgo.Describe("Deployments", func() {
		var (
			deploymentName                = "fake-deployment-name"
			agentID                       = "fake-agent-id"
			jobName                       = "fake-job-name"
			jobID                         = "fake-job-id"
			jobIndex                      = 0
			jobBootstrap                  = true
			jobIP                         = "1.2.3.4"
			jobAZ                         = "fake-job-az"
			jobVMType                     = "fake-job-vm-type"
			jobResourcePool               = "fake-job-resource-pool"
			jobResurrectionPause          = true
			jobVMID                       = "fake-job-vmid"
			processState                  = "running"
			jobUptimeSeconds              = uint64(3600)
			jobLoadAvg01                  = float64(0.01)
			jobLoadAvg05                  = float64(0.05)
			jobLoadAvg15                  = float64(0.15)
			jobCPUSys                     = float64(0.5)
			jobCPUUser                    = float64(1.0)
			jobCPUWait                    = float64(1.5)
			jobMemKB                      = 1000
			jobMemPercent                 = 10
			jobSwapKB                     = 2000
			jobSwapPercent                = 20
			jobSystemDiskInodePercent     = 10
			jobSystemDiskPercent          = 20
			jobEphemeralDiskInodePercent  = 30
			jobEphemeralDiskPercent       = 40
			jobPersistentDiskInodePercent = 50
			jobPersistentDiskPercent      = 60
			jobProcessName                = "fake-process-name"
			jobProcessState               = "running"
			jobProcessUptimeSeconds       = uint64(3600)
			jobProcessCPUTotal            = float64(0.5)
			jobProcessMemKB               = uint64(2000)
			jobProcessMemPercent          = float64(20)
			releaseName                   = "fake-release-name"
			releaseVersion                = "1.2.3"
			releaseJob1Name               = "fake-release-job1-name"
			releaseJob2Name               = "fake-release-job2-name"
			releasePackage1Name           = "fake-release-package1-name"
			releasePackage2Name           = "fake-release-package2-name"
			stemcellName                  = "fake-stemcell-name"
			stemcellVersion               = "4.5.6"
			stemcellOSName                = "fake-stemcell-os-name"

			processes  []director.VMInfoProcess
			vitals     director.VMInfoVitals
			instances  []director.VMInfo
			release    director.Release
			releases   []director.Release
			stemcell   director.Stemcell
			stemcells  []director.Stemcell
			depls      []director.Deployment
			deployment director.Deployment

			deploymentsInfo         []deployments.DeploymentInfo
			expectedDeploymentsInfo []deployments.DeploymentInfo
		)

		ginkgo.BeforeEach(func() {
			processes = []director.VMInfoProcess{
				{
					Name:   jobProcessName,
					State:  jobProcessState,
					CPU:    director.VMInfoVitalsCPU{Total: &jobProcessCPUTotal},
					Mem:    director.VMInfoVitalsMemIntSize{KB: &jobProcessMemKB, Percent: &jobProcessMemPercent},
					Uptime: director.VMInfoVitalsUptime{Seconds: &jobProcessUptimeSeconds},
				},
			}

			vitals = director.VMInfoVitals{
				CPU: director.VMInfoVitalsCPU{
					Sys:  strconv.FormatFloat(jobCPUSys, 'E', -1, 64),
					User: strconv.FormatFloat(jobCPUUser, 'E', -1, 64),
					Wait: strconv.FormatFloat(jobCPUWait, 'E', -1, 64),
				},
				Mem: director.VMInfoVitalsMemSize{
					KB:      strconv.Itoa(jobMemKB),
					Percent: strconv.Itoa(jobMemPercent),
				},
				Swap: director.VMInfoVitalsMemSize{
					KB:      strconv.Itoa(jobSwapKB),
					Percent: strconv.Itoa(jobSwapPercent),
				},
				Uptime: director.VMInfoVitalsUptime{
					Seconds: &jobUptimeSeconds,
				},
				Load: []string{
					strconv.FormatFloat(jobLoadAvg01, 'E', -1, 64),
					strconv.FormatFloat(jobLoadAvg05, 'E', -1, 64),
					strconv.FormatFloat(jobLoadAvg15, 'E', -1, 64),
				},
				Disk: map[string]director.VMInfoVitalsDiskSize{
					"system": {
						InodePercent: strconv.Itoa(jobSystemDiskInodePercent),
						Percent:      strconv.Itoa(jobSystemDiskPercent),
					},
					"ephemeral": {
						InodePercent: strconv.Itoa(jobEphemeralDiskInodePercent),
						Percent:      strconv.Itoa(jobEphemeralDiskPercent),
					},
					"persistent": {
						InodePercent: strconv.Itoa(jobPersistentDiskInodePercent),
						Percent:      strconv.Itoa(jobPersistentDiskPercent),
					},
				},
			}

			instances = []director.VMInfo{
				{
					AgentID:            agentID,
					JobName:            jobName,
					ID:                 jobID,
					Index:              &jobIndex,
					Bootstrap:          jobBootstrap,
					ProcessState:       processState,
					IPs:                []string{jobIP},
					AZ:                 jobAZ,
					VMType:             jobVMType,
					ResourcePool:       jobResourcePool,
					ResurrectionPaused: jobResurrectionPause,
					VMID:               jobVMID,
					Vitals:             vitals,
					Processes:          processes,
				},
			}

			release = &directorfakes.FakeRelease{
				NameStub:    func() string { return releaseName },
				VersionStub: func() version.Version { return version.MustNewVersionFromString(releaseVersion) },
				JobsStub: func() ([]director.Job, error) {
					return []director.Job{{Name: releaseJob1Name}, {Name: releaseJob2Name}}, nil
				},
				PackagesStub: func() ([]director.Package, error) {
					return []director.Package{{Name: releasePackage1Name}, {Name: releasePackage2Name}}, nil
				},
			}
			releases = []director.Release{release}

			stemcell = &directorfakes.FakeStemcell{
				NameStub:    func() string { return stemcellName },
				VersionStub: func() version.Version { return version.MustNewVersionFromString(stemcellVersion) },
				OSNameStub:  func() string { return stemcellOSName },
			}
			stemcells = []director.Stemcell{stemcell}

			deployment = &directorfakes.FakeDeployment{
				NameStub:          func() string { return deploymentName },
				InstanceInfosStub: func() ([]director.VMInfo, error) { return instances, nil },
				ReleasesStub:      func() ([]director.Release, error) { return releases, nil },
				StemcellsStub:     func() ([]director.Stemcell, error) { return stemcells, nil },
			}

			depls = []director.Deployment{deployment}
			boshClient.DeploymentsReturns(depls, nil)

			expectedDeploymentsInfo = []deployments.DeploymentInfo{
				{
					Name: deploymentName,
					Instances: []deployments.Instance{
						{
							AgentID:            agentID,
							Name:               jobName,
							ID:                 jobID,
							Index:              strconv.Itoa(jobIndex),
							Bootstrap:          jobBootstrap,
							IPs:                []string{jobIP},
							AZ:                 jobAZ,
							VMType:             jobVMType,
							ResourcePool:       jobResourcePool,
							ResurrectionPaused: jobResurrectionPause,
							Healthy:            true,
							Processes: []deployments.Process{
								{
									Name:    jobProcessName,
									Uptime:  &jobProcessUptimeSeconds,
									Healthy: true,
									CPU:     deployments.CPU{Total: &jobProcessCPUTotal},
									Mem:     deployments.MemInt{KB: &jobProcessMemKB, Percent: &jobProcessMemPercent},
								},
							},
							Vitals: deployments.Vitals{
								CPU: deployments.CPU{
									Sys:  strconv.FormatFloat(jobCPUSys, 'E', -1, 64),
									User: strconv.FormatFloat(jobCPUUser, 'E', -1, 64),
									Wait: strconv.FormatFloat(jobCPUWait, 'E', -1, 64),
								},
								Mem: deployments.Mem{
									KB:      strconv.Itoa(jobMemKB),
									Percent: strconv.Itoa(jobMemPercent),
								},
								Swap: deployments.Mem{
									KB:      strconv.Itoa(jobSwapKB),
									Percent: strconv.Itoa(jobSwapPercent),
								},
								Uptime: &jobUptimeSeconds,
								Load: []string{
									strconv.FormatFloat(jobLoadAvg01, 'E', -1, 64),
									strconv.FormatFloat(jobLoadAvg05, 'E', -1, 64),
									strconv.FormatFloat(jobLoadAvg15, 'E', -1, 64),
								},
								SystemDisk: deployments.Disk{
									InodePercent: strconv.Itoa(jobSystemDiskInodePercent),
									Percent:      strconv.Itoa(jobSystemDiskPercent),
								},
								EphemeralDisk: deployments.Disk{
									InodePercent: strconv.Itoa(jobEphemeralDiskInodePercent),
									Percent:      strconv.Itoa(jobEphemeralDiskPercent),
								},
								PersistentDisk: deployments.Disk{
									InodePercent: strconv.Itoa(jobPersistentDiskInodePercent),
									Percent:      strconv.Itoa(jobPersistentDiskPercent),
								},
							},
						},
					},
					Releases: []deployments.Release{
						{Name: releaseName, Version: releaseVersion,
							JobNames:     []string{releaseJob1Name, releaseJob2Name},
							PackageNames: []string{releasePackage1Name, releasePackage2Name},
						},
					},
					Stemcells: []deployments.Stemcell{
						{Name: stemcellName, Version: stemcellVersion, OSName: stemcellOSName},
					},
				},
			}
		})

		ginkgo.JustBeforeEach(func() {
			deploymentsInfo, err = deploymentsFetcher.Deployments()
		})

		ginkgo.It("returns the deployments", func() {
			gomega.Expect(deploymentsInfo).To(gomega.Equal(expectedDeploymentsInfo))
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.Context("when instance has no VMID", func() {
			ginkgo.BeforeEach(func() {
				instances[0].VMID = ""
				deployment = &directorfakes.FakeDeployment{
					NameStub:      func() string { return deploymentName },
					ReleasesStub:  func() ([]director.Release, error) { return releases, nil },
					StemcellsStub: func() ([]director.Stemcell, error) { return stemcells, nil },
				}
				depls = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(depls, nil)
			})

			ginkgo.It("does not return the instance", func() {
				gomega.Expect(deploymentsInfo[0].Instances).To(gomega.BeEmpty())
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when there are no deployments", func() {
			ginkgo.BeforeEach(func() {
				boshClient.DeploymentsReturns([]director.Deployment{}, nil)
			})

			ginkgo.It("does not return deployments", func() {
				gomega.Expect(deploymentsInfo).To(gomega.BeEmpty())
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when it fails to get the deployment", func() {
			ginkgo.BeforeEach(func() {
				boshClient.DeploymentsReturns([]director.Deployment{}, errors.New("no deployments"))
			})

			ginkgo.It("does not return deployments", func() {
				gomega.Expect(deploymentsInfo).To(gomega.BeEmpty())
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when there are no instances", func() {
			ginkgo.BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:      func() string { return deploymentName },
					ReleasesStub:  func() ([]director.Release, error) { return releases, nil },
					StemcellsStub: func() ([]director.Stemcell, error) { return stemcells, nil },
				}
				depls = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(depls, nil)
			})

			ginkgo.It("does not return instances", func() {
				gomega.Expect(deploymentsInfo[0].Instances).To(gomega.BeEmpty())
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when it fails to get the deployment instances", func() {
			ginkgo.BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:          func() string { return deploymentName },
					InstanceInfosStub: func() ([]director.VMInfo, error) { return nil, errors.New("no instances") },
				}
				depls = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(depls, nil)
			})

			ginkgo.It("does not return deployments", func() {
				gomega.Expect(deploymentsInfo).To(gomega.BeEmpty())
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when there are no releases", func() {
			ginkgo.BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:          func() string { return deploymentName },
					InstanceInfosStub: func() ([]director.VMInfo, error) { return instances, nil },
					StemcellsStub:     func() ([]director.Stemcell, error) { return stemcells, nil },
				}
				depls = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(depls, nil)
			})

			ginkgo.It("does not return releases", func() {
				gomega.Expect(deploymentsInfo[0].Releases).To(gomega.BeEmpty())
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when it fails to get the deployment releases", func() {
			ginkgo.BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:     func() string { return deploymentName },
					ReleasesStub: func() ([]director.Release, error) { return nil, errors.New("no releases") },
				}
				depls = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(depls, nil)
			})

			ginkgo.It("does not return deployments", func() {
				gomega.Expect(deploymentsInfo).To(gomega.BeEmpty())
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when there are no stemcells", func() {
			ginkgo.BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:          func() string { return deploymentName },
					InstanceInfosStub: func() ([]director.VMInfo, error) { return instances, nil },
					ReleasesStub:      func() ([]director.Release, error) { return releases, nil },
				}
				depls = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(depls, nil)
			})

			ginkgo.It("does not return stemcells", func() {
				gomega.Expect(deploymentsInfo[0].Stemcells).To(gomega.BeEmpty())
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when it fails to get the deployment stemcells", func() {
			ginkgo.BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:      func() string { return deploymentName },
					StemcellsStub: func() ([]director.Stemcell, error) { return nil, errors.New("no stemcells") },
				}
				depls = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(depls, nil)
			})

			ginkgo.It("does not return deployments", func() {
				gomega.Expect(deploymentsInfo).To(gomega.BeEmpty())
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})
	})
})
