// provides abstraction over reading a writing the lock file, as well as
// comparing it to the configured installed packages in the configuration
package lock

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/xnacly/mehr/types"
)

var path = LookUp()
var lock *types.LockFile = &types.LockFile{Packages: map[string]*types.Package{}}

func init() {
	toml.DecodeFile(path, lock)
}

// path to the lock file
func LookUp() string {
	confHome, _ := os.UserConfigDir()
	return filepath.Join(confHome, "mehr", "lock.toml")
}

// return all packages permanently installed on the system
func Permanent(config *types.Configuration, lock *types.LockFile) map[string]*types.Package {
	perm := map[string]*types.Package{}
	for k, v := range config.Packages {
		if _, ok := lock.Packages[k]; ok {
			perm[k] = v
		}
	}
	return perm
}

// return all temporary installed packages on the system
func Temporary(config *types.Configuration, lock *types.LockFile) map[string]*types.Package {
	temp := map[string]*types.Package{}
	for k, v := range config.Packages {
		if _, ok := lock.Packages[k]; !ok {
			temp[k] = v
		}
	}
	return temp
}

// reads the lockfile from disk on first call, every call afterwards omits the
// disk interaction and reads the cached lock file
func Get() *types.LockFile {
	return lock
}

// adds entry to the lock file, writes it to disk and updates the cached lock file
func Write(name string, entry *types.Package) error {
	lock.Packages[name] = entry
	val, err := os.Stat(path)
	if err == nil && val.IsDir() {
		return fmt.Errorf("Lockfile %q is a directory, how did that happen?", path)
	} else {
		os.MkdirAll(filepath.Dir(path), 0777)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Failed to open %q, %s", path, err)
	}
	defer file.Close()

	e := toml.NewEncoder(file)
	return e.Encode(lock)
}
