package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"sync"
	"time"

	"github.com/iamthemuffinman/overseer/config"
	"github.com/iamthemuffinman/overseer/pkg/buildspec"
	"github.com/iamthemuffinman/overseer/pkg/hammer"
	"github.com/iamthemuffinman/overseer/pkg/hostspec"
	"github.com/iamthemuffinman/overseer/pkg/knife"
	"github.com/iamthemuffinman/overseer/pkg/workerpool"

	"github.com/iamthemuffinman/cli"
	log "github.com/iamthemuffinman/logsip"
	"github.com/mitchellh/go-homedir"
	flag "github.com/ogier/pflag"
)

type ProvisionVirtualCommand struct {
	UI         cli.Ui
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

		specfile := c.FlagSet.StringP("buildspec", "h", "", "Provide a buildspec for your host(s) (i.e. indy.prod.kafka)")

		// Parse everything after 3 arguments (i.e overseer provision virtual STARTHERE)
		c.FlagSet.Parse(os.Args[3:])

		// GTFO if a buildspec wasn't specified
		if *specfile == "" {
			log.Fatal("You must specify a buildspec")
		}

		home, err := getHomeDir()
		if err != nil {
			log.Fatalf("unable to retrieve users home directory: %s", err)
		}

		// Parse overseer's config file which contains usernames and passwords
		cspec, err := config.ParseFile(fmt.Sprintf("%s/.overseer/overseer.conf", home))
		if err != nil {
			log.Fatalf("unable to parse overseer config: %s", err)
		}

		// Here is where we essentially parse the entire buildspecs directory to find
		// the buildspec specified on the command line.
		bspec, err := buildspec.ParseDir("/etc/overseer/buildspecs", *specfile)
		if err != nil {
			log.Fatalf("unable to parse buildspec: %s", err)
		}

		hammerCmd := &hammer.Hammer{
			Username:          cspec.Foreman.Username,
			Password:          cspec.Foreman.Password,
			Hostname:          "",
			Organization:      bspec.Foreman.Organization,
			Location:          bspec.Foreman.Location,
			Hostgroup:         bspec.Foreman.Hostgroup,
			Environment:       bspec.Foreman.Environment,
			PartitionTableID:  bspec.Foreman.PartitionTableID,
			OperatingSystemID: bspec.Foreman.OperatingSystemID,
			Medium:            bspec.Foreman.Medium,
			ArchitectureID:    bspec.Foreman.ArchitectureID,
			DomainID:          bspec.Foreman.DomainID,
			ComputeProfile:    bspec.Foreman.ComputeProfile,
			ComputeResource:   bspec.Foreman.ComputeResource,
			Host: hammer.Host{
				CPUs:   bspec.Virtual.CPUs,
				Cores:  bspec.Virtual.Cores,
				Memory: bspec.Virtual.Memory,
				Disks:  bspec.Vsphere.Devices.Disks,
			},
		}

		knifeCmd := &knife.Knife{
			Hostname:    "",
			Environment: bspec.Chef.Environment,
			RunList:     bspec.Chef.RunList,
		}

		var wg sync.WaitGroup

		// If there are arguments, then the user has specified a host on the
		// command line rather than using a hostspec
		if len(c.FlagSet.Args()) > 0 {
			log.Errorf("Please use a hostspec instead of specifying hosts on the command line")
			os.Exit(1)
		} else {
			// Parse the hostspec in the current directory to get a list of hosts
			hspec, err := hostspec.ParseFile("./hostspec")
			if err != nil {
				log.Fatalf("couldn't find your hostspec: %s", err)
			}

			// Range over all the hosts in the hostspec
			for _, host := range hspec.Hosts {
				hammerCmd.Hostname = host
				// Execute is a method that will send the command to a job queue
				// to be processed by a goroutine. This way we can build more
				// hosts at the same time by executing hammer in parallel.
				if err := hammerCmd.Execute(); err != nil {
					log.Fatalf("error executing hammer: %s", err)
				}

				wg.Add(1)
				go func(host string) {
					defer wg.Done()
					for {
						// GetBuildStatus will return 0 if Foreman says the host has been
						// build successfully. We'll wait until all hosts have been built
						// successfully and then we'll execute hammer.
						status, err := hammerCmd.GetBuildStatus()
						if err != nil {
							log.Fatalf("error executing hammer: %s", err)
						}

						if status == 0 {
							log.Infof("%s built successfully!", host)
							break
						} else {
							time.Sleep(1 * time.Minute)
						}
					}
				}(hammerCmd.Hostname)
			}

			wg.Wait()

			for _, host := range hspec.Hosts {
				knifeCmd.Hostname = host

				wg.Add(1)
				go func(host string) {
					defer wg.Done()
					// Add all recipes/cookbooks/roles to the run list
					// of each node
					if err := knifeCmd.AddToRunList(); err != nil {
						log.Fatalf("error executing knife: %s", err)
					}
				}(hammerCmd.Hostname)
			}

			wg.Wait()

			log.Info("All hosts successfully created and chef'd!")
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

// Get user's home directory so we can pass it to the config parser
func getHomeDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		// If for some reason the above doesn't work, let's see what the standard library
		// can do for us here. If this doesn't work, something is wrong and we should
		// cut out at this point.
		currentUser, err := user.Current()
		if err != nil {
			return "", err
		}

		return currentUser.HomeDir, nil
	}
	return home, nil
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
