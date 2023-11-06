package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xnacly/mehr/config"
	"github.com/xnacly/mehr/lock"
	l "github.com/xnacly/mehr/log"
	"github.com/xnacly/mehr/types"
)

func init() {
	listCmd.PersistentFlags().BoolP("temporary", "t", false, "only print temporary installed packages")
	listCmd.PersistentFlags().BoolP("permanent", "p", false, "only print installed packages also listed in the default configuration")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			configPath = config.LookUp()
		}

		conf, err := config.Get(configPath)
		if err != nil {
			l.Errorf("Failed to get config: %q", err)
			return
		}

		all := lock.Get().Packages

		if temp, err := cmd.Flags().GetBool("temporary"); err == nil && temp {
			printPackages(lock.Temporary(conf, lock.Get()))
			return
		} else if perm, err := cmd.Flags().GetBool("permanent"); err == nil && perm {
			printPackages(lock.Permanent(conf, lock.Get()))
			return
		} else {
			temp := lock.Temporary(conf, lock.Get())
			l.Infof("Found %d packages installed on your system, %d of them temporary", len(all), len(temp))
			printPackages(all)
		}
	},
}

func printPackages(packages map[string]*types.Package) {
	for k, v := range packages {
		if v.Version != "" {
			fmt.Println(k + "@" + v.Version)
		} else {
			fmt.Println(k)
		}
	}
}
