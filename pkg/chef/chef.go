package chef

import (
	"io/ioutil"

	"github.com/go-chef/chef"
	"github.com/iamthemuffinman/overseer/pkg/util"
)

type Chef struct {
}

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

func AddToRunList(name, runListItem string) chef.Node {
	node := chef.NewNode(name)
	node.RunList = append(node.RunList, runListItem)
	return node
}

func ReadKey(keyPath string) (string, error) {
	path, err := util.ExpandPath(keyPath)
	if err != nil {
		return "", err
	}

	key, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

func UpdateNode(name, keyPath, chefServer string, runList []string) error {
	key, err := ReadKey(keyPath)
	if err != nil {
		return err
	}

	client, err := NewClient(key, chefServer)
	if err != nil {
		return err
	}

	var node chef.Node
	for _, item := range runList {
		node = AddToRunList(name, item)
	}

	client.Nodes.Put(node)
	return nil
}
