package knife

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"text/template"

	log "github.com/iamthemuffinman/logsip"
	"github.com/mitchellh/go-homedir"
)

type Knife struct {
	Hostname    string
	Hosts       string
	Username    string
	HomeDir     string
	Environment string
	Role        string
	Cookbook    string
	ChefServer  string
}

const knifeTemplate = `
current_dir = File.dirname(__FILE__)
log_level                :info
log_location             STDOUT
node_name                "{{.Username}}"
client_key               "#{current_dir}/{{.Username}}.pem"
validation_client_name   ""
validation_key           "#{current_dir}/{{.Username}}.pem"
chef_server_url          "{{.ChefServer}}"
cookbook_path            ["{{.HomeDir}}/src/chef/os-cookbooks", "{{.HomeDir}}/src/chef/site-cookbooks"]
knife[:ssh_user] = ""
knife[:ssh_password] = ""
`

func New() *Knife {
	user, err := user.Current()
	if err != nil {
		return nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return nil
	}

	return &Knife{
		Hostname:    "",
		Username:    user.Username,
		HomeDir:     home,
		Environment: "",
		Role:        "",
		Cookbook:    "",
		ChefServer:  "",
	}
}

func (k *Knife) GenerateAndInitialize() error {
	if err := os.Mkdir(fmt.Sprintf("%s/.chef", k.HomeDir), 0644); err != nil {
		return err
	}

	knifeConfig, err := os.Create(fmt.Sprintf("%s/.chef/knife.rb", k.HomeDir))
	if err != nil {
		return err
	}

	compiled, err := template.New("knife_template").Parse(knifeTemplate)
	if err != nil {
		return err
	}

	if err := compiled.Execute(knifeConfig, k); err != nil {
		return err
	}

	knife := exec.Command("knife", "ssl", "fetch")

	log.Infof("Executing: %s", strings.Join(knife.Args, " "))

	knife.Stdout = os.Stdout
	knife.Stderr = os.Stderr

	if err := knife.Run(); err != nil {
		return err
	}

	knife = exec.Command("knife", "ssl", "check")

	log.Infof("Executing: %s", strings.Join(knife.Args, " "))

	knife.Stdout = os.Stdout
	knife.Stderr = os.Stderr

	if err := knife.Run(); err != nil {
		return err
	}

	return nil
}
