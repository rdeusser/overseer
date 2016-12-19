package hammer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/iamthemuffinman/overseer/pkg/hostspec"
	"github.com/iamthemuffinman/overseer/pkg/workerpool"

	log "github.com/iamthemuffinman/logsip"
)

// This will all be replaced by https://github.com/iamthemuffinman/go-foreman
// when it can support all of the options here.

type Hammer struct {
	Username          string
	Password          string
	Hostname          string
	Organization      string
	Location          string
	Hostgroup         string
	Environment       string
	PartitionTableID  int
	OperatingSystemID int
	Medium            string
	ArchitectureID    int
	DomainID          int
	ComputeProfile    string
	ComputeResource   string
	Host              Host
}

type Host struct {
	CPUs   int
	Cores  int
	Memory int
	Disks  []*hostspec.Disk
}

func (h *Hammer) joinVolumes() string {
	var volumes []string

	for _, disk := range h.Host.Disks {
		volumes = append(volumes, fmt.Sprintf("--volume size_gb=%d", disk.Size))
	}

	return strings.Join(volumes, ", ")
}

func (h *Hammer) joinComputeAttributes() string {
	computeAttributes := fmt.Sprintf("start=1,cpus=%d,corespersocket=%d,memory_mb=%d", h.Host.CPUs, h.Host.Cores, h.Host.Memory)

	return computeAttributes
}

func (h *Hammer) Execute() error {

	// Build massive hammer command
	// I can't wait for this to be replace by go-foreman
	hammer := exec.Command("hammer", fmt.Sprintf(`-u %q -p %q host create --name %q --organization %q --location %q
	--hostgroup-title %q --environment %q --partition-table-id %q --operatingsystem-id %q --medium %q --architecture-id %q
	--domain-id %q --subnet %q --compute-profile %q --compute-attributes %q %q --compute-resource %q`,
		h.Username, h.Password, h.Hostname, h.Organization, h.Location, h.Location, h.Environment, h.PartitionTableID,
		h.OperatingSystemID, h.Medium, h.ArchitectureID, h.DomainID, h.Location, h.ComputeProfile, h.joinComputeAttributes(),
		h.joinVolumes(), h.ComputeResource))

	hammer.Stdout = os.Stdout
	hammer.Stderr = os.Stderr

	log.Infof("Executing: %s", strings.Join(hammer.Args, " "))

	// Create a job to run the hammer command
	job := workerpool.Job{Command: hammer}

	// Push the job onto the queue
	workerpool.JobQueue <- job

	return nil
}
