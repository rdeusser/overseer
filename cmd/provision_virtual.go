package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/iamthemuffinman/overseer/config"
	"github.com/iamthemuffinman/overseer/pkg/buildspec"
	"github.com/iamthemuffinman/overseer/pkg/hammer"
	"github.com/iamthemuffinman/overseer/pkg/hostspec"
	"github.com/iamthemuffinman/overseer/pkg/workerpool"

	"github.com/iamthemuffinman/cli"
	log "github.com/iamthemuffinman/logsip"
	"github.com/mitchellh/go-homedir"
	flag "github.com/ogier/pflag"
)

type ProvisionVirtualCommand struct {
	Ui         cli.Ui
	FlagSet    *flag.FlagSet
	ShutdownCh <-chan struct{}
}

func (c *ProvisionVirtualCommand) Run(args []string) int {
	if len(args) == 0 {
		return cli.RunResultHelp
	}

	for _, arg := range args {
		if arg == "-h" || arg == "-help" || arg == "--help" {
			return cli.RunResultHelp
		}
	}

	// Okay, we're ready to start doing some work at this point.
	// Let's create the pool of workers so they can start listening
	// for jobs that are put into the JobQueue.
	dispatcher := workerpool.NewDispatcher()
	dispatcher.Run()

	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		c.FlagSet = flag.NewFlagSet("virtual", flag.ExitOnError)

		specfile := c.FlagSet.StringP("hostspec", "h", "", "Provide a specfile name for your host(s) (i.e. indy.prod.kafka)")

		// Parse everything after 3 arguments (i.e overseer provision virtual STARTHERE)
		c.FlagSet.Parse(os.Args[3:])

		// GTFO if a hostspec wasn't specified
		if *specfile == "" {
			log.Fatal("You must specify a hostspec")
		}

		// Get user's home directory so we can pass it to the config parser
		home, err := homedir.Dir()
		if err != nil {
			// If for some reason the above doesn't work, let's see what the standard library
			// can do for us here. If this doesn't work, something is wrong and we should
			// cut out at this point.
			currentUser, err := user.Current()
			if err != nil {
				log.Fatalf("unable to get the home directory of the user running this process")
			}

			home = currentUser.HomeDir
		}

		// Parse overseer's config file which contains usernames and passwords
		conf, err := config.ParseFile(fmt.Sprintf("%s/.overseer/overseer.conf", home))
		if err != nil {
			log.Fatalf("unable to parse overseer config: %s", err)
		}

		// Here is where we essentially parse the entire hostspecs directory to find
		// the hostspec specified on the command line.
		hspec, err := hostspec.ParseDir("/etc/overseer/hostspecs", *specfile)
		if err != nil {
			log.Fatalf("unable to parse hostspec: %s", err)
		}

		// If there are arguments, then the user has specified a host on the
		// command line rather than using a buildspec
		if len(c.FlagSet.Args()) > 0 {
			log.Errorf("Please use a buildspec instead of specifying hosts on the command line")
			os.Exit(1)
		} else {
			// Parse the buildspec in the current directory to get a list of hosts
			bspec, err := buildspec.ParseFile("./buildspec")
			if err != nil {
				log.Fatalf("couldn't find your buildspec: %s", err)
			}

			// Same thing as above - range over all the hosts in the buildspec
			for _, host := range bspec.Hosts {
				cmd := &hammer.Hammer{
					Username:          conf.Foreman.Username,
					Password:          conf.Foreman.Password,
					Hostname:          host,
					Organization:      hspec.Foreman.Organization,
					Location:          hspec.Foreman.Location,
					Hostgroup:         hspec.Foreman.Hostgroup,
					Environment:       hspec.Foreman.Environment,
					PartitionTableID:  hspec.Foreman.PartitionTableID,
					OperatingSystemID: hspec.Foreman.OperatingSystemID,
					Medium:            hspec.Foreman.Medium,
					ArchitectureID:    hspec.Foreman.ArchitectureID,
					DomainID:          hspec.Foreman.DomainID,
					ComputeProfile:    hspec.Foreman.ComputeProfile,
					ComputeResource:   hspec.Foreman.ComputeResource,
					Host: hammer.Host{
						CPUs:   hspec.Virtual.CPUs,
						Cores:  hspec.Virtual.Cores,
						Memory: hspec.Virtual.Memory,
						Disks:  hspec.Vsphere.Devices.Disks,
					},
				}

				// Execute is a method that will send the command to a job queue
				// to be processed by a goroutine. This way we can build more
				// hosts at the same time by executing hammer in parallel.
				if err := cmd.Execute(); err != nil {
					log.Fatalf("error executing hammer: %s", err)
				}

				// Run chef/knife stuff here
			}
		}
	}()

	select {
	case <-c.ShutdownCh:
		log.Info("Interrupt received. Gracefully shutting down...")

		// Stop execution here
		// need to either find out or do something here about removing data for all hosts
		// or just the current host

		select {
		case <-c.ShutdownCh:
			log.Warn("Two interrupts received - exiting immediately. Some things may not have finished and no cleanup will be attempted.")
			return 1
		case <-doneCh:
		}
	case <-doneCh:
	}

	return 0
}

func (c *ProvisionVirtualCommand) Help() string {
	return c.helpProvisionVirtual()
}

func (c *ProvisionVirtualCommand) Synopsis() string {
	return "Provision virtual infrastructure"
}

func (c *ProvisionVirtualCommand) helpProvisionVirtual() string {
	helpText := `
Usage: overseer provision virtual [OPTIONS] [HOSTS]
`
	return strings.TrimSpace(helpText)
}
