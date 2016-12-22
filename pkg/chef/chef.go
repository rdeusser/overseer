package chef

import (
	"io/ioutil"
	"strings"

	"github.com/go-chef/chef"
	"github.com/mitchellh/go-homedir"
)

func NewClient(key, chefServer string) (*chef.Client, error) {
	client, err := chef.NewClient(&chef.Config{
		Name:    "overseer",
		Key:     key,
		BaseURL: chefServer,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewHost(name, env string, runList []string) *chef.Node {
	return &chef.Node{
		Name:        name,
		Environment: env,
		ChefType:    "node",
		JsonClass:   "Chef::Node",
		RunList:     runList,
	}
}

func ReadValidationKey(validationKey string) (string, error) {
	if validationKey[:1] != "~" {
		key, err := ioutil.ReadFile(validationKey)
		if err != nil {
			return "", err
		}

		return string(key), nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	expandedKeyPath := strings.Replace(validationKey, "~", home, 1)

	key, err := ioutil.ReadFile(expandedKeyPath)
	if err != nil {
		return "", err
	}

	return string(key), nil
}
