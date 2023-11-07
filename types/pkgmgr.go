package types

type Manager interface {
	Install(packages map[string]*Package) (error, int)
	Upgrade(packages map[string]*Package) (error, int)
	Remove(packages ...string) (error, int)
	Exists() bool
	Update() error
}

type SubCommand struct {
	Name    string
	Options []string
}
