// abstracts package managers away in a operating system independent way
package pkgmgr

import (
	"os"
	"os/exec"
	"strings"

	"github.com/xnacly/mehr/lock"
	"github.com/xnacly/mehr/log"
	"github.com/xnacly/mehr/types"
)

var defaultBuffer strings.Builder

type PackageManager struct {
	Name                 string                                   // name of the executable
	install              *types.SubCommand                        // command to be executed for installing packages
	upgrade              *types.SubCommand                        // command to be executed for updating packages
	remove               *types.SubCommand                        // command to be executed for removing packages
	update               *types.SubCommand                        // command to be executed for updating source / fetching new package data
	noSudo               bool                                     // some package managers do not require sudo, true if so
	options              []string                                 // options for all sub commands
	formatPkgWithVersion func(name string, version string) string // used to format every package before attempting to install it
}

func (p *PackageManager) createCmd(c *types.SubCommand, packages ...string) error {
	if c.Name == "" {
		return nil
	}
	args := []string{p.Name, c.Name}
	args = append(args, p.options...)
	args = append(args, c.Options...)
	args = append(args, packages...)

	log.Infof("running %q", strings.Join(args, " "))
	var cmd *exec.Cmd
	if p.noSudo {
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		cmd = exec.Command("sudo", args...)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *PackageManager) Install(packages map[string]*types.Package) (error, int) {
	pkgs := make([]string, 0)
	for k, v := range packages {
		if _, ok := lock.Get().Packages[k]; ok {
			log.Warnf("Package %s@%s already installed, skipping", k, v.Version)
			continue
		}
		if v.Version != "" {
			pkgs = append(pkgs, p.formatPkgWithVersion(k, v.Version))
		} else {
			pkgs = append(pkgs, k)
		}
	}
	if len(pkgs) == 0 {
		return nil, 0
	}
	err := p.createCmd(p.install, pkgs...)
	if err != nil {
		return err, 0
	}
	for k, v := range packages {
		err := lock.AddPackage(k, p.Name, v)
		if err != nil {
			return err, 0
		}
	}
	return nil, len(pkgs)
}

func (p *PackageManager) Upgrade(packages map[string]*types.Package) (error, int) {
	pkgs := make([]string, 0)
	for k, v := range packages {
		if v.Version != "" {
			pkgs = append(pkgs, p.formatPkgWithVersion(k, v.Version))
		} else {
			pkgs = append(pkgs, k)
		}
	}
	if len(pkgs) == 0 {
		return nil, 0
	}
	err := p.createCmd(p.upgrade, pkgs...)
	if err != nil {
		return err, 0
	}
	for k, v := range packages {
		err := lock.AddPackage(k, p.Name, v)
		if err != nil {
			return err, 0
		}
	}
	return nil, len(pkgs)
}

func (p *PackageManager) Remove(packages ...string) (error, int) {
	if len(packages) == 0 {
		return nil, 0
	}
	err := p.createCmd(p.remove, packages...)
	if err != nil {
		return err, 0
	}
	for _, pkg := range packages {
		err := lock.RemovePackage(pkg, p.Name)
		if err != nil {
			return err, 0
		}
	}
	return nil, len(packages)
}

func (p *PackageManager) Update() error {
	return p.createCmd(p.update)
}

func (p *PackageManager) Exists() bool {
	_, err := exec.LookPath(p.Name)
	return err == nil
}
