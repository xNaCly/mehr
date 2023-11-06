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
	options              []string                                 // options for all sub commands
	buffer               *strings.Builder                         // buffer to write stdout to
	formatPkgWithVersion func(name string, version string) string // used to format every package before attempting to install it
}

func (p *PackageManager) createCmd(c *types.SubCommand, packages ...string) *exec.Cmd {
	args := make([]string, 0)
	args = append(args, c.Name)
	args = append(args, p.options...)
	args = append(args, c.Options...)
	args = append(args, packages...)
	cmd := exec.Command(p.Name, args...)

	cmd.Stdin = os.Stdin
	if p.buffer == nil {
		p.buffer = &defaultBuffer
	}
	cmd.Stderr = p.buffer
	cmd.Stdout = p.buffer
	log.Infof("running '%s %s'", p.Name, strings.Join(args, " "))
	return cmd
}

func (p *PackageManager) Install(packages map[string]*types.Package) (error, int) {
	pkgs := make([]string, 0)
	for k, v := range packages {
		if _, ok := lock.Get().Packages[k]; ok {
			log.Warnf("Package %s already installed, skipping", k)
			continue
		}
		if v.Version != "" {
			pkgs = append(pkgs, p.formatPkgWithVersion(k, v.Version))
		} else {
			pkgs = append(pkgs, k)
		}
		err := lock.Write(k, v)
		if err != nil {
			return err, 0
		}
	}
	if len(pkgs) == 0 {
		return nil, 0
	}
	return p.createCmd(p.install, pkgs...).Run(), len(pkgs)
}

func (p *PackageManager) Upgrade(packages ...string) error {
	return p.createCmd(p.upgrade, packages...).Run()
}

func (p *PackageManager) Remove(packages ...string) error {
	return p.createCmd(p.remove, packages...).Run()
}

func (p *PackageManager) Update() error {
	return p.createCmd(p.update).Run()
}

func (p *PackageManager) Exists() (string, bool) {
	path, err := exec.LookPath(p.Name)
	return path, err == nil
}

func (p *PackageManager) Output() string {
	if p.buffer == nil {
		return ""
	}
	return p.buffer.String()
}
