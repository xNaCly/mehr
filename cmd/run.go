package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xnacly/mehr/config"
	l "github.com/xnacly/mehr/log"
	"github.com/xnacly/mehr/types"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute mehr commands",
	Args:  cobra.ExactArgs(1),
	Long: `Run a defined command in its working directory and with its env variables:

Command definiton:

    [command]
    [command.l]
    cmd = "ls -la"
    cwd = ".."
    [command.l.env]
    VISUAL = "less"
`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.LookUp()
		conf, err := config.Get(configPath)
		if err != nil {
			l.Errorf("Failed to get config: %q", err)
			return
		}

		if args[0] == "help" {
			for k, v := range conf.Commands {
				fmt.Printf("%s: \n\tCommand: %q\n", k, v.Cmd)
				if len(v.Env) > 0 {
					fmt.Print("\tEnv:")
					for variable, val := range v.Env {
						fmt.Printf("\n\t\t%s: %q", variable, val)
					}
				}
				fmt.Println()
			}
			return
		} else if command, ok := conf.Commands[args[0]]; ok {
			err := runCommand(args[0], &command)
			if err != nil {
				l.Errorf("Failed to run command %q: %s", command.Cmd, err)
			}
		} else {
			l.Errorf("Command %q undefined, use 'help' for listing available commands", args[0])
			return
		}
	},
}

func runCommand(name string, cmd *types.Command) error {
	l.Infof("Got command: %q with cwd: %q and env: %q", cmd.Cmd, cmd.Cwd, cmd.Env)

	args := strings.Split(cmd.Cmd, " ")
	command := exec.Command(args[0], args[1:]...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if cmd.Cwd != "" {
		path, err := filepath.Abs(cmd.Cwd)
		if err != nil {
			return err
		}
		command.Dir = path
		l.Infof("Set command working directory to %q", path)
	}

	var env []string

	if cmd.PurgeOsEnv {
		l.Infof("Got %q.PurgeOsEnv, removing variables from child command, only keeping defined env variables", cmd.Cmd)
		if len(cmd.Env) > 0 {
			env = make([]string, 0, len(cmd.Env))
		}
	} else {
		env = os.Environ()
	}

	for k, v := range cmd.Env {
		env = append(env, k+"="+v)
		l.Infof("Setting env variable: '%s=%s'", k, v)
	}

	command.Env = env

	err := command.Run()
	if err == nil {
		l.Infof("Ran command %q successfully", name)
	}
	return err
}
