package packagemanagers

var Managers = []*PackageManager{
	{
		Name:    "apt",
		install: &SubCommand{Name: "install"},
		update:  &SubCommand{Name: "update"},
		upgrade: &SubCommand{Name: "upgrade"},
		remove:  &SubCommand{Name: "remove"},
	},
	{
		Name:    "pacman",
		install: &SubCommand{Name: "", Options: []string{"-S"}},
		update:  &SubCommand{Name: "", Options: []string{"-Sy"}},
		upgrade: &SubCommand{Name: "", Options: []string{"-Su"}},
		remove:  &SubCommand{Name: "", Options: []string{"-Rs"}},
	},
}

func Get() (*PackageManager, bool) {
	for _, mgr := range Managers {
		if mgr.Exists() {
			return mgr, true
		}
	}
	return nil, false
}
