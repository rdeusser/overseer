package util

import (
	"strings"

	"github.com/mitchellh/go-homedir"
)

func ExpandPath(path string) (string, error) {
	if path[:1] != "~" {
		return path, nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	expandedKeyPath := strings.Replace(path, "~", home, 1)

	return expandedKeyPath, nil
}
