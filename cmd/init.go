package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xnacly/mehr/config"
	l "github.com/xnacly/mehr/log"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise mehr",
	Args:  cobra.MaximumNArgs(1),
	Long: `Initialise mehr by creating a new mehr.toml at the default configuration location or the specified directory with:
Initialise the whole system:

    init 

Create the default configuration file at the specified path:

    init <path>

    init .

    init testdir 
    `,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.LookUp()

		if len(args) == 1 {
			t, err := filepath.Abs(filepath.Join(args[0], "mehr.toml"))
			if err != nil {
				l.Errorf("Failed to create an absolute path representation of %q: %s", configPath, err)
				return
			}
			configPath = t
		}

		if val, err := os.Stat(configPath); err == nil {
			if val.IsDir() {
				l.Error("Configuration file is a directory")
				return
			}
			// we ignore the error here because the flag is false if an error
			// occurred
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				l.Errorf("Configuration file %q already exists, use '--force' to override this check", configPath)
				return
			}
			l.Warnf("Got force, overwriting already existing configuration file")
		}

		err := os.MkdirAll(filepath.Dir(configPath), 0777)
		if err != nil {
			l.Errorf("Failed to create all directories to configuration file: %q", err)
			return
		}
		err = os.WriteFile(configPath, config.DefaultConfigFileContent, 0777)
		if err != nil {
			l.Errorf("Failed to write default configuration to configuration file: %q", err)
			return
		}
		l.Infof("Wrote default configuration file to %q", configPath)
	},
}
