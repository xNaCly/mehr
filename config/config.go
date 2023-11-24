package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "embed"

	"github.com/BurntSushi/toml"
	"github.com/xnacly/mehr/pkgmgr"
	"github.com/xnacly/mehr/types"
)

// only filled once, afterwards used to cache the configuration - this allows
// us to skip file reading and toml decoding
var configCache *types.Configuration

//go:embed mehr.toml
var DefaultConfigFileContent []byte

func IsRoot() bool {
	// group id for elevated process is 0: "the login group for the superuser
	// must have GID 0" (https://en.wikipedia.org/wiki/Group_identifier)
	return os.Getegid() == 0
}

// returns path the mehr config is located
func LookUp() string {
	confHome, _ := os.UserConfigDir()
	return filepath.Join(confHome, "mehr", "mehr.toml")
}

func CreateDirIfNotExist() {
	path := LookUp()
	os.MkdirAll(filepath.Dir(path), 0777)
}

// reads configuration at 'path', decodes to Configuration struct
func Get(path string) (*types.Configuration, error) {
	if configCache != nil {
		return configCache, nil
	}

	c := &types.Configuration{}

	_, err := toml.DecodeFile(path, c)
	if err != nil {
		return c, fmt.Errorf("Failed to decode configuration file: %w", err)
	}

	err = ValidateConfig(c)
	if err != nil {
		return c, err
	}

	configCache = c
	return c, nil
}

// validates the struct and fills empty fields
func ValidateConfig(c *types.Configuration) error {
	if c.PackageManager == "" || c.PackageManager == "auto" {
		if mgr, ok := pkgmgr.Get(); ok {
			c.PackageManager = mgr.Name
		} else {
			return errors.New("No package manager found")
		}
	}
	if c.Packages == nil {
		c.Packages = make(map[string]map[string]*types.Package)
	} else {
		cpy := c.Packages["$"]
		delete(c.Packages, "$")
		if el, ok := c.Packages[c.PackageManager]; ok {
			for k, v := range el {
				cpy[k] = v
			}
		}
		c.Packages[c.PackageManager] = cpy
	}
	if c.SystemConfig == nil {
		c.SystemConfig = &types.SystemConfig{}
	}
	if c.SystemConfig.Path == "" {
		confHome, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("Failed to get user config dir: %w", err)
		}
		c.SystemConfig.Path = confHome
	}
	return nil
}
