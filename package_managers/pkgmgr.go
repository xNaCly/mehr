package packagemanagers

import (
	"os"
	"os/exec"
)

type Manager interface {
	Install(packages []string) error
	Upgrade(packages []string) error
	Remove(packages []string) error
	Exists() bool
	Update() error
}

type SubCommand struct {
	Name    string
	Options []string
}

type PackageManager struct {
	Name    string      // name of the executable
	install *SubCommand // command to be executed for installing packages
	upgrade *SubCommand // command to be executed for updating packages
	remove  *SubCommand // command to be executed for removing packages
	update  *SubCommand // command to be executed for updating source / fetching new package data
	Options []string    // options for all sub commands
}

func (p *PackageManager) createCmd(c *SubCommand, packages ...string) *exec.Cmd {
	args := make([]string, 0)
	args = append(args, c.Name)
	args = append(args, p.Options...)
	args = append(args, c.Options...)
	args = append(args, packages...)
	cmd := exec.Command(p.Name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func (p *PackageManager) Install(packages ...string) error {
	return p.createCmd(p.install, packages...).Run()
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
