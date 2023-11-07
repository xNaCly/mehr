package types

type LockFile struct {
	Packages map[string]*Package `toml:"packages"`
}
