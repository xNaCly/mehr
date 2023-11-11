package types

// configure a package
type PackageConfig struct {
	Content string `toml:"content"` // content to write to configuration file
	URL     string `toml:"url"`     // location of configuration to write to configuration file
	Path    string `toml:"path"`    // location of the file contents to copy to the configuration file
	Link    string `toml:"link"`    // location of the file to link the configuration file to
}

type SystemConfig struct {
	// path to store package configuration at, falls back to `$XDG_CONFIG_HOME` for
	// linux, %AppData% for windows and $HOME/Library/Application Support/ for
	// macOS
	Path string `toml:"path"`
	// specify configuration files from various sources
	Files map[string]*PackageConfig `toml:"files"`
}

// install a singular package
type Package struct {
	Version string `toml:"version"` // leave empty for latest
	URL     string `toml:"url"`     // url to download package from
}

type Configuration struct {
	PackageManager string              `toml:"package-manager"` // specify what package manager to use, path or empty for auto lookup
	Packages       map[string]*Package `toml:"packages"`        // packages to install
	SystemConfig   *SystemConfig       `toml:"config"`          // configure the system and the installed packages
}
