package deployments

type DeploymentInfo struct {
	Name      string
	Instances []Instance
	Releases  []Release
	Stemcells []Stemcell
}

type Instance struct {
	AgentID            string
	VMID               string
	Name               string
	ID                 string
	Index              string
	Bootstrap          bool
	IPs                []string
	AZ                 string
	VMType             string
	ResourcePool       string
	ResurrectionPaused bool
	Healthy            bool
	Processes          []Process
	Vitals             Vitals
}

type Process struct {
	Name    string
	Uptime  *uint64
	Healthy bool
	CPU     CPU
	Mem     MemInt
}

type Vitals struct {
	CPU            CPU
	Mem            Mem
	Swap           Mem
	Uptime         *uint64
	Load           []string
	SystemDisk     Disk
	EphemeralDisk  Disk
	PersistentDisk Disk
}

type CPU struct {
	Total *float64
	Sys   string
	User  string
	Wait  string
}

type Mem struct {
	KB      string
	Percent string
}

type MemInt struct {
	KB      *uint64
	Percent *float64
}

type Disk struct {
	InodePercent string
	Percent      string
}

type Release struct {
	Name    string
	Version string
}

type Stemcell struct {
	Name    string
	Version string
	OSName  string
}
