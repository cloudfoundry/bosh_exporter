package deployments_test

import (
	"errors"
	"flag"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-cli/director/directorfakes"
	"github.com/cppforlife/go-semi-semantic/version"

	"github.com/cloudfoundry-community/bosh_exporter/filters"

	. "github.com/cloudfoundry-community/bosh_exporter/deployments"
)

func init() {
	flag.Set("log.level", "fatal")
}

var _ = Describe("Fetcher", func() {
	var (
		err                error
		boshDeployments    []string
		boshClient         *directorfakes.FakeDirector
		deploymentsFilter  *filters.DeploymentsFilter
		deploymentsFetcher *Fetcher
	)

	BeforeEach(func() {
		boshDeployments = []string{}
		boshClient = &directorfakes.FakeDirector{}
	})

	JustBeforeEach(func() {
		deploymentsFilter = filters.NewDeploymentsFilter(boshDeployments, boshClient)
		deploymentsFetcher = NewFetcher(*deploymentsFilter)
	})

	Describe("Deployments", func() {
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
			stemcellName                  = "fake-stemcell-name"
			stemcellVersion               = "4.5.6"
			stemcellOSName                = "fake-stemcell-os-name"

			processes   []director.VMInfoProcess
			vitals      director.VMInfoVitals
			instances   []director.VMInfo
			release     director.Release
			releases    []director.Release
			stemcell    director.Stemcell
			stemcells   []director.Stemcell
			deployments []director.Deployment
			deployment  director.Deployment

			deploymentsInfo         []DeploymentInfo
			expectedDeploymentsInfo []DeploymentInfo
		)

		BeforeEach(func() {
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
					"system": director.VMInfoVitalsDiskSize{
						InodePercent: strconv.Itoa(int(jobSystemDiskInodePercent)),
						Percent:      strconv.Itoa(int(jobSystemDiskPercent)),
					},
					"ephemeral": director.VMInfoVitalsDiskSize{
						InodePercent: strconv.Itoa(int(jobEphemeralDiskInodePercent)),
						Percent:      strconv.Itoa(int(jobEphemeralDiskPercent)),
					},
					"persistent": director.VMInfoVitalsDiskSize{
						InodePercent: strconv.Itoa(int(jobPersistentDiskInodePercent)),
						Percent:      strconv.Itoa(int(jobPersistentDiskPercent)),
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

			deployments = []director.Deployment{deployment}
			boshClient.DeploymentsReturns(deployments, nil)

			expectedDeploymentsInfo = []DeploymentInfo{
				DeploymentInfo{
					Name: deploymentName,
					Instances: []Instance{
						Instance{
							AgentID:            agentID,
							Name:               jobName,
							ID:                 jobID,
							Index:              strconv.Itoa(int(jobIndex)),
							Bootstrap:          jobBootstrap,
							IPs:                []string{jobIP},
							AZ:                 jobAZ,
							VMType:             jobVMType,
							ResourcePool:       jobResourcePool,
							ResurrectionPaused: jobResurrectionPause,
							Healthy:            true,
							Processes: []Process{
								Process{
									Name:    jobProcessName,
									Uptime:  &jobProcessUptimeSeconds,
									Healthy: true,
									CPU:     CPU{Total: &jobProcessCPUTotal},
									Mem:     MemInt{KB: &jobProcessMemKB, Percent: &jobProcessMemPercent},
								},
							},
							Vitals: Vitals{
								CPU: CPU{
									Sys:  strconv.FormatFloat(jobCPUSys, 'E', -1, 64),
									User: strconv.FormatFloat(jobCPUUser, 'E', -1, 64),
									Wait: strconv.FormatFloat(jobCPUWait, 'E', -1, 64),
								},
								Mem: Mem{
									KB:      strconv.Itoa(jobMemKB),
									Percent: strconv.Itoa(jobMemPercent),
								},
								Swap: Mem{
									KB:      strconv.Itoa(jobSwapKB),
									Percent: strconv.Itoa(jobSwapPercent),
								},
								Uptime: &jobUptimeSeconds,
								Load: []string{
									strconv.FormatFloat(jobLoadAvg01, 'E', -1, 64),
									strconv.FormatFloat(jobLoadAvg05, 'E', -1, 64),
									strconv.FormatFloat(jobLoadAvg15, 'E', -1, 64),
								},
								SystemDisk: Disk{
									InodePercent: strconv.Itoa(int(jobSystemDiskInodePercent)),
									Percent:      strconv.Itoa(int(jobSystemDiskPercent)),
								},
								EphemeralDisk: Disk{
									InodePercent: strconv.Itoa(int(jobEphemeralDiskInodePercent)),
									Percent:      strconv.Itoa(int(jobEphemeralDiskPercent)),
								},
								PersistentDisk: Disk{
									InodePercent: strconv.Itoa(int(jobPersistentDiskInodePercent)),
									Percent:      strconv.Itoa(int(jobPersistentDiskPercent)),
								},
							},
						},
					},
					Releases: []Release{
						Release{Name: releaseName, Version: releaseVersion},
					},
					Stemcells: []Stemcell{
						Stemcell{Name: stemcellName, Version: stemcellVersion, OSName: stemcellOSName},
					},
				},
			}
		})

		JustBeforeEach(func() {
			deploymentsInfo, err = deploymentsFetcher.Deployments()
		})

		It("returns the deployments", func() {
			Expect(deploymentsInfo).To(Equal(expectedDeploymentsInfo))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when instance has no VMID", func() {
			BeforeEach(func() {
				instances[0].VMID = ""
				deployment = &directorfakes.FakeDeployment{
					NameStub:      func() string { return deploymentName },
					ReleasesStub:  func() ([]director.Release, error) { return releases, nil },
					StemcellsStub: func() ([]director.Stemcell, error) { return stemcells, nil },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("does not return the instance", func() {
				Expect(deploymentsInfo[0].Instances).To(BeEmpty())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when there are no deployments", func() {
			BeforeEach(func() {
				boshClient.DeploymentsReturns([]director.Deployment{}, nil)
			})

			It("does not return deployments", func() {
				Expect(deploymentsInfo).To(BeEmpty())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when it fails to get the deployment", func() {
			BeforeEach(func() {
				boshClient.DeploymentsReturns([]director.Deployment{}, errors.New("no deployments"))
			})

			It("does not return deployments", func() {
				Expect(deploymentsInfo).To(BeEmpty())
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when there are no instances", func() {
			BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:      func() string { return deploymentName },
					ReleasesStub:  func() ([]director.Release, error) { return releases, nil },
					StemcellsStub: func() ([]director.Stemcell, error) { return stemcells, nil },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("does not return instances", func() {
				Expect(deploymentsInfo[0].Instances).To(BeEmpty())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when it fails to get the deployment instances", func() {
			BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:          func() string { return deploymentName },
					InstanceInfosStub: func() ([]director.VMInfo, error) { return nil, errors.New("no instances") },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("does not return deployments", func() {
				Expect(deploymentsInfo).To(BeEmpty())
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when there are no releases", func() {
			BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:          func() string { return deploymentName },
					InstanceInfosStub: func() ([]director.VMInfo, error) { return instances, nil },
					StemcellsStub:     func() ([]director.Stemcell, error) { return stemcells, nil },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("does not return releases", func() {
				Expect(deploymentsInfo[0].Releases).To(BeEmpty())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when it fails to get the deployment releases", func() {
			BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:     func() string { return deploymentName },
					ReleasesStub: func() ([]director.Release, error) { return nil, errors.New("no releases") },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("does not return deployments", func() {
				Expect(deploymentsInfo).To(BeEmpty())
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when there are no stemcells", func() {
			BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:          func() string { return deploymentName },
					InstanceInfosStub: func() ([]director.VMInfo, error) { return instances, nil },
					ReleasesStub:      func() ([]director.Release, error) { return releases, nil },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("does not return stemcells", func() {
				Expect(deploymentsInfo[0].Stemcells).To(BeEmpty())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when it fails to get the deployment stemcells", func() {
			BeforeEach(func() {
				deployment = &directorfakes.FakeDeployment{
					NameStub:      func() string { return deploymentName },
					StemcellsStub: func() ([]director.Stemcell, error) { return nil, errors.New("no stemcells") },
				}
				deployments = []director.Deployment{deployment}
				boshClient.DeploymentsReturns(deployments, nil)
			})

			It("does not return deployments", func() {
				Expect(deploymentsInfo).To(BeEmpty())
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
