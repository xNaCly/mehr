package pkgmgr

import (
	"os"
	"os/exec"
	"strings"

	"github.com/xnacly/mehr/log"
	"github.com/xnacly/mehr/types"
)

var defaultBuffer strings.Builder

type PackageManager struct {
	Name                 string            // name of the executable
	install              *types.SubCommand // command to be executed for installing packages
	upgrade              *types.SubCommand // command to be executed for updating packages
	remove               *types.SubCommand // command to be executed for removing packages
	update               *types.SubCommand // command to be executed for updating source / fetching new package data
	options              []string          // options for all sub commands
	buffer               *strings.Builder  // buffer to write stdout to
	formatPkgWithVersion func(name string, version string) string
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

func (p *PackageManager) Install(packages map[string]*types.Package) error {
	pkgs := make([]string, len(packages))
	i := 0
	for k, v := range packages {
		if v.Version != "" {
			pkgs[i] = p.formatPkgWithVersion(k, v.Version)
		} else {
			pkgs[i] = k
		}
		i++
	}
	return p.createCmd(p.install, pkgs...).Run()
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
	return p.buffer.String()
}
