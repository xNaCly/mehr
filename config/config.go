package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
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
