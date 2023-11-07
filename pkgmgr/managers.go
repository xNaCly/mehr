package pkgmgr

import (
	"fmt"

	"github.com/xnacly/mehr/types"
)

var Managers = []*PackageManager{
	{
		Name:    "apt",
		options: []string{"-y"}, // skips [Y/n] confirmation prompts
		install: &types.SubCommand{Name: "install"},
		update:  &types.SubCommand{Name: "update"},
		upgrade: &types.SubCommand{Name: "upgrade", Options: []string{
			"--only-upgrade", // only upgrades the specified packages
		}},
		remove: &types.SubCommand{Name: "remove"},
		formatPkgWithVersion: func(name, version string) string {
			return name + "=" + version
		},
	},
	{
		Name:    "pacman",
		options: []string{"--noconfirm"}, // skips [Y/n] confirmation prompts
		install: &types.SubCommand{Name: "-S"},
		update:  &types.SubCommand{Name: "-Sy"},
		upgrade: &types.SubCommand{Name: "-Su"},
		remove:  &types.SubCommand{Name: "-Rs"},
		formatPkgWithVersion: func(name, version string) string {
			return name + "=" + version
		},
	},
}

// returns the package manager for name if found, otherwise errors
func GetByName(name string) (*PackageManager, error) {
	for _, mgr := range Managers {
		if mgr.Name == name {
			path, ok := mgr.Exists()
			if !ok {
				return nil, fmt.Errorf("Package manager %q not found on operating system", mgr.Name)
			}
			mgr.Name = path
			return mgr, nil
		}
	}
	return nil, fmt.Errorf("Package manager %q not found in defined list of package managers", name)
}

// returns the currently available package manager on the system
func Get() (*PackageManager, bool) {
	for _, mgr := range Managers {
		if path, ok := mgr.Exists(); ok {
			mgr.Name = path
			return mgr, true
		}
	}
	return nil, false
}
