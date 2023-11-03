package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	packagemanagers "github.com/xnacly/mehr/package_managers"
)

// returns path the mehr config is located
func LookUp() string {
	confHome, _ := os.UserConfigDir()
	return filepath.Join(confHome, "mehr", "mehr.toml")
}

// reads configuration at 'path', decodes to Configuration struct
func Get(path string) (*Configuration, error) {
	c := &Configuration{}
	c.Ok()
	out, err := os.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("Failed to read from configuration file: %w", err)
	}
	_, err = toml.Decode(string(out), c)
	if err != nil {
		return c, fmt.Errorf("Failed to decode configuration file: %w", err)
	}
	return c, nil
}

type PackageConfig struct {
	Content string `toml:"content"`
}

type Package struct {
	Version string                   `toml:"version"` // leave empty for latest
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
