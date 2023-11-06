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

type cache struct {
	config *types.Configuration
	path   string
}

// only filled once, afterwards used to cache the configuration - this allows
// us to skip file reading and toml decoding
var configCache *cache = &cache{}

//go:embed mehr.toml
var DefaultConfigFileContent []byte

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
	if configCache != nil && configCache.path == path {
		return configCache.config, nil
	}
	c, err := ValidateConfig(types.Configuration{})
	if err != nil {
		return c, err
	}
	out, err := os.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("Failed to read configuration file, using fallback: %w", err)
	}
	_, err = toml.Decode(string(out), c)
	if err != nil {
		return c, fmt.Errorf("Failed to decode configuration file: %w", err)
	}
	configCache.config = c
	configCache.path = path
	return c, nil
}

// validates the struct and fills empty fields
func ValidateConfig(c types.Configuration) (*types.Configuration, error) {
	if c.PackageManager == "" || c.PackageManager == "auto" {
		if mgr, ok := pkgmgr.Get(); ok {
			c.PackageManager = mgr.Name
		} else {
			return nil, errors.New("No package manager found")
		}
	}
	if c.SystemConfig == nil {
		c.SystemConfig = &types.SystemConfig{}
	}
	if c.SystemConfig.Path == "" {
		confHome, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("Failed to get user config dir: %w", err)
		}
		c.SystemConfig.Path = confHome
	}
	return &c, nil
}
