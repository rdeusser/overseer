package knife

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/iamthemuffinman/overseer/pkg/workerpool"

	log "github.com/iamthemuffinman/logsip"
)

func (k *Knife) Bootstrap() error {
	knife := exec.Command("knife", fmt.Sprintf("bootstrap %s -E %s -r 'role[%s]' --sudo --use-sudo-password", k.Hostname, k.Environment, k.BaseRole))

	log.Infof("Executing: %s", strings.Join(knife.Args, " "))

	knife.Stdout = os.Stdout
	knife.Stderr = os.Stderr

	// Create a job to run the knife command
	job := workerpool.Job{Command: knife}

	// Push the job onto the queue
	workerpool.JobQueue <- job

	return nil
}
