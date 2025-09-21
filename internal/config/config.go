// Package config
package config

import (
	"errors"
	"os"
	"path"
)

const (
	AppName    = "tmux-sessionpane"
	AppVersion = "v0.2.1"
)

var (
	BaseConfigPath = path.Join(os.Getenv("HOME"), ".config", AppName)
)

func Init() error {
	_, err := os.Stat(BaseConfigPath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(BaseConfigPath, 0o755)
	}

	return err
}
