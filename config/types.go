package config

import (
	"errors"
	"fmt"
	"os"

	packagemanagers "github.com/xnacly/mehr/package_managers"
)

type PackageConfig struct {
	Content string `toml:"content"`
}

type Package struct {
	Version string                   `toml:"version"` // leave empty for latest
	URL     string                   `toml:"url"`     // url to download package from
	Config  map[string]PackageConfig `toml:"config"`
}

type Configuration struct {
	// path to store package configuration at, falls back to `$XDG_CONFIG_HOME` for
	// linux, %AppData% for windows and $HOME/Library/Application Support/ for
	// macOS
	Path           string             `toml:"config-path"`
	PackageManager string             `toml:"package-manager"` // path or empty for auto lookup
	Packages       map[string]Package `toml:"packages"`
}

// validates the struct and fills empty fields
func (c *Configuration) Ok() error {
	if c.PackageManager == "" || c.PackageManager == "auto" {
		if mgr, ok := packagemanagers.Get(); ok {
			c.PackageManager = mgr.Name
		} else {
			return errors.New("No package manager found")
		}
	}
	if c.Path == "" {
		confHome, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("Failed to get user config dir: %w", err)
		}
		c.Path = confHome
	}
	return nil
}
