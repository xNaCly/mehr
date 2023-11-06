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
		upgrade: &types.SubCommand{Name: "upgrade"},
		remove:  &types.SubCommand{Name: "remove"},
	},
	{
		Name:    "pacman",
		install: &types.SubCommand{Name: "-S", Options: []string{"--noconfirm"}},
		update:  &types.SubCommand{Name: "-Sy"},
		upgrade: &types.SubCommand{Name: "-Su", Options: []string{"--noconfirm"}},
		remove:  &types.SubCommand{Name: "-Rs", Options: []string{"--noconfirm"}},
		formatPkgWithVersion: func(name, version string) string {
			return name + "=" + version
		},
	},
	{
		Name: "npm",
		install: &types.SubCommand{Name: "install", Options: []string{"-g"}},
		update:  &types.SubCommand{Name: "update", Options: []string{"-g"}},
		remove:  &types.SubCommand{Name: "uninstall", Options: []string{"-g"}},
	},
	{
		Name: "pnpm",
		install: &types.SubCommand{Name: "add", Options: []string{"-g"}},
		update:  &types.SubCommand{Name: "update", Options: []string{"-g"}},
		remove:  &types.SubCommand{Name: "remove", Options: []string{"-g"}},
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
