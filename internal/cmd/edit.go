package cmd

import (
	"fmt"

	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/logic"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var EditCommand types.Command

func init() {
	EditCommand = types.Command{
		Name:        "edit",
		Description: "Edit settings of an existing job",
		Usage:       "edit <job_name> [--interval <n>] [--command <cmd>] [--retries <n>] [--timeout <n>]",
		Flags: []types.Flag{
			{Name: "--interval", Description: "New interval in seconds", Required: false},
			{Name: "--command", Description: "New shell command or URL", Required: false},
			{Name: "--retries", Description: "New retry count", Required: false},
			{Name: "--timeout", Description: "New per-run timeout in seconds", Required: false},
		},
		Execute: editHandler,
	}
}

// editHandler applies partial updates to an existing job's configuration.
func editHandler(args []string) {
	if len(args) < 1 {
		util.ErrorMessage(fmt.Sprint("Usage: ", EditCommand.Usage))
		return
	}

	name := args[0]
	flags := util.ParseFlags(args[1:])

	job, err := db.GetJob(name)
	if err != nil {
		util.ErrorMessage(fmt.Sprintf("Error looking up %q: %v", name, err))
		return
	}

	if job == nil {
		util.ErrorMessage(fmt.Sprintf("No job named %q", name))
		return
	}

	if raw, ok := flags["--interval"]; ok {
		v := util.AtoiOrDefault(raw, 0)
		if v <= 0 {
			util.ErrorMessage(fmt.Sprintf("Invalid --interval %q", raw))
			return
		}

		job.Interval = v
	}

	if cmd, ok := flags["--command"]; ok {
		job.Command = cmd
	}

	if raw, ok := flags["--retries"]; ok {
		v := util.AtoiOrDefault(raw, -1)
		if v < 0 {
			util.ErrorMessage(fmt.Sprintf("Invalid --retries %q", raw))
			return
		}

		job.Retries = v
	}

	if raw, ok := flags["--timeout"]; ok {
		v := util.AtoiOrDefault(raw, -1)
		if v < 0 {
			util.ErrorMessage(fmt.Sprintf("Invalid --timeout %q", raw))
			return
		}

		job.Timeout = v
	}

	if err := logic.EditJob(name, *job, util.ResolvePaths().ConfigPath); err != nil {
		util.ErrorMessage(fmt.Sprintf("Failed to edit %q: %v", name, err))
	} else {
		util.SuccessMessage(fmt.Sprintf("Job %q updated", name))
	}
}
