package packagemanagers

var Managers = []*PackageManager{
	{
		Name:    "apt",
		Options: []string{"-y"}, // skips [Y/n] confirmation prompts
		install: &SubCommand{Name: "install"},
		update:  &SubCommand{Name: "update"},
		upgrade: &SubCommand{Name: "upgrade"},
		remove:  &SubCommand{Name: "remove"},
	},
	{
		Name:    "pacman",
		install: &SubCommand{Name: "-S", Options: []string{"--no-confirm"}},
		update:  &SubCommand{Name: "-Sy"},
		upgrade: &SubCommand{Name: "-Su", Options: []string{"--no-confirm"}},
		remove:  &SubCommand{Name: "-Rs", Options: []string{"--no-confirm"}},
	},
	{
		Name: "npm",
		install: &SubCommand{Name: "install", Options: []string{"-g"}},
		update:  &SubCommand{Name: "update", Options: []string{"-g"}},
		remove:  &SubCommand{Name: "uninstall", Options: []string{"-g"}},
	},
	{
		Name: "pnpm",
		install: &SubCommand{Name: "add", Options: []string{"-g"}},
		update:  &SubCommand{Name: "update", Options: []string{"-g"}},
		remove:  &SubCommand{Name: "remove", Options: []string{"-g"}},
	},
}

func Get() (*PackageManager, bool) {
	for _, mgr := range Managers {
		if path, ok := mgr.Exists(); ok {
			mgr.Name = path
			return mgr, true
		}
	}
	return nil, false
}
