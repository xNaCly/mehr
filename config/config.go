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
	// TODO: XDG_CONFIG_HOME is undefined in sudo env, modified to
	// /root/.config/ meaning our config does not live at the correct path.
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
	c, err := ValidateConfig(&types.Configuration{})
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

	configCache = c
	return c, nil
}

// validates the struct and fills empty fields
func ValidateConfig(c *types.Configuration) (*types.Configuration, error) {
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
	return c, nil
}
