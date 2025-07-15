package cli

import (
	"fmt"
	"os"

	"github.com/Ege-Okyay/jobctl/internal/cmd"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var commands = map[string]types.Command{
	"add":     cmd.AddCommand,
	"debug":   cmd.DebugCommand,
	"config":  cmd.ConfigCommand,
	"list":    cmd.ListCommand,
	"next":    cmd.NextCommand,
	"pause":   cmd.PauseCommand,
	"resume":  cmd.ResumeCommand,
	"status":  cmd.StatusCommand,
	"dry-run": cmd.DryRunCommand,
	"enable":  cmd.EnableCommand,
	"disable": cmd.DisableCommand,
	"delete":  cmd.DeleteCommand,
	"run":     cmd.RunCommand,
	"history": cmd.HistoryCommand,
	"inspect": cmd.InspectCommand,
	"edit":    cmd.EditCommand,
}

// Setup is the main entry point for the CLI.
// It parses command-line arguments, finds the appropriate command, validates flags,
// and executes the command's handler function.
func Setup() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "help" {
		util.PrintHelp(commands)
		return
	}

	cmdName := args[0]
	cmd, exists := commands[cmdName]
	if !exists {
		util.HandleUnknownCommand(commands, cmdName)
		return
	}

	if len(args) > 1 && util.IsHelpFlag(args[1]) {
		util.PrintCommandHelp(cmd)
		return
	}

	flags := util.ParseFlags(args[1:])

	if err := util.ValidateFlags(flags, cmd.Flags); err != nil {
		fmt.Printf("Error: %v\n\n%s\n", err, cmd.Usage)
		return
	}

	cmd.Execute(args[1:])
}
