package knife

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/iamthemuffinman/overseer/pkg/workerpool"

	log "github.com/iamthemuffinman/logsip"
)

func (k *Knife) AddToRunList() error {
	runList := strings.Join(k.RunList, ",")
	knife := exec.Command("knife", fmt.Sprintf("node run_list add %q %q", k.Hostname, runList))

	log.Infof("Executing: %s", strings.Join(knife.Args, " "))

	knife.Stdout = os.Stdout
	knife.Stderr = os.Stderr

	// Create a job to run the knife command
	job := workerpool.Job{Command: knife}

	// Push the job onto the queue
	workerpool.JobQueue <- job

	return nil
}
