package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xnacly/mehr/config"
	"github.com/xnacly/mehr/lock"
	l "github.com/xnacly/mehr/log"
)

func init() {
	listCmd.PersistentFlags().BoolP("temporary", "t", false, "only print temporary installed packages")
	listCmd.PersistentFlags().BoolP("permanent", "p", false, "only print installed packages also listed in the default configuration")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages",
	Example: `
    list -t
    list -p
    `,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.LookUp()

		conf, err := config.Get(configPath)
		if err != nil {
			l.Errorf("Failed to get config: %q", err)
			return
		}

		all := lock.All(lock.Get())

		if temp, err := cmd.Flags().GetBool("temporary"); err == nil && temp {
			printPackages(lock.Temporary(conf, lock.Get()))
			return
		} else if perm, err := cmd.Flags().GetBool("permanent"); err == nil && perm {
			printPackages(lock.Permanent(conf, lock.Get()))
			return
		} else {
			printPackages(all)
		}
	},
}

func printPackages(packages map[string][]string) {
	i := 1
	for manager, pkgs := range packages {
		if len(pkgs) == 0 {
			continue
		}
		fmt.Println(manager)
		for _, name := range pkgs {
			fmt.Print("(", i, ") ")
			fmt.Println(name)
			i++
		}
	}
}
