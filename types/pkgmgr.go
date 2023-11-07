package types

type Manager interface {
	Install(packages []Package) (error, int)
	Upgrade(packages []string) error
	Remove(packages []string) error
	Exists() bool
	Update() error
}

type SubCommand struct {
	Name    string
	Options []string
}
