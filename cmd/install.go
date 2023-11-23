package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xnacly/mehr/config"
	l "github.com/xnacly/mehr/log"
	"github.com/xnacly/mehr/pkgmgr"
	"github.com/xnacly/mehr/types"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install packages",
	Example: `
    install neovim
    install neovim kitty
    install 
    `,
	Long: `Install a single, multiple or all configured packages via: 

Install a single package:

    install [package]

Install multiple packages:

    install [package...]

Install all packages in the current configuration:

    install

Packages installed via this command are not automatically added to the
configuration and are therefore referred to as termporary packages. Restore
your system state to your configuration via 'sync' - removing temporary
packages not found in the configuration.

See 'mehr sync help' for more information.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.LookUp()

		conf, err := config.Get(configPath)
		if err != nil {
			l.Errorf("Failed to get config: %q", err)
			return
		}

		if len(args) == 0 {
			// installs packages from configuration
			for mgr, pkgs := range conf.Packages {
				if len(pkgs) == 0 {
					continue
				}
				var manager *pkgmgr.PackageManager
				if mgr == "$" {
					var ok bool
					manager, ok = pkgmgr.Get()
					if !ok {
						l.Error("Failed to find a package manager")
						return
					}
				} else {
					manager, err = pkgmgr.GetByName(mgr)
					if err != nil {
						l.Error(err)
						return
					}
				}

				err, amount := manager.Install(pkgs)
				if err != nil {
					l.Errorf("failed to install packages: %s", err)
				} else if amount > 0 {
					l.Infof("Installed %d packages", len(conf.Packages))
				} else {
					l.Infof("Did nothing, exiting")
				}
			}
		} else {
			var manager *pkgmgr.PackageManager
			if conf.PackageManager == "auto" || conf.PackageManager == "" {
				mgr, ok := pkgmgr.Get()

				if !ok {
					l.Error("Failed to find a package manager")
					return
				}

				manager = mgr
			} else {
				var err error
				mgr, err := pkgmgr.GetByName(conf.PackageManager)
				if err != nil {
					l.Error(err)
					return
				}
				manager = mgr
			}

			// installs packages with default package manager
			pkgs := map[string]*types.Package{}
			for _, pkg := range args {
				pkgs[pkg] = &types.Package{}
			}
			err, amount := manager.Install(pkgs)
			if err != nil {
				l.Errorf("failed to install packages: %s", err)
			} else if amount > 0 {
				l.Infof("Installed %d packages", len(conf.Packages))
			}
		}
	},
}
