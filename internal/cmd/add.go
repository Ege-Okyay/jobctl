package cmd

import (
	"fmt"
	"strconv"

	"github.com/Ege-Okyay/jobctl/internal/config"
	"github.com/Ege-Okyay/jobctl/internal/logic"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var AddCommand types.Command

func init() {
	AddCommand = types.Command{
		Name:        "add",
		Description: "Add a new scheduled job",
		Usage:       "add --name <name> --interval <seconds> --command <cmd> [--retries <n>] [--timeout <seconds>]",
		Flags: []types.Flag{
			{Name: "--name", Description: "Job name", Required: true},
			{Name: "--interval", Description: "Interval in seconds", Required: true},
			{Name: "--command", Description: "Shell command or URL", Required: true},
			{Name: "--retries", Description: "Number of automatic re-attempts on failure (optional)", Required: false},
			{Name: "--timeout", Description: "Maximum execution time in seconds (optional)", Required: false},
		},
		Execute: addHandler,
	}
}

func addHandler(args []string) {
	flags := util.ParseFlags(args)

	intervalStr := flags["--interval"]
	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval <= 0 {
		util.ErrorMessage(fmt.Sprintf("Invalid --interval %q: must be a positive integer\n", intervalStr))
		return
	}

	rawCmd := flags["--command"]
	if rawCmd == "" {
		util.ErrorMessage("Error: --command cannot be empty")
		return
	}

	// The command string might be quoted; unquote it for accurate execution.
	cmdStr, err := strconv.Unquote(rawCmd)
	if err != nil {
		cmdStr = rawCmd // Use the raw string if unquoting fails.
	}

	retries := util.AtoiOrDefault(flags["--retries"], 0)
	timeout := util.AtoiOrDefault(flags["--timeout"], 0)

	job := types.JobConfig{
		Name:     flags["--name"],
		Interval: interval,
		Command:  cmdStr,
		Retries:  retries,
		Timeout:  timeout,
		Enabled:  true, // New jobs are enabled by default.
	}

	configPath, err := config.ConfigPath()
	if err != nil {
		util.ErrorMessage("Error: couldn't find the config path")
		return
	}

	if err := logic.AddJob(job, configPath); err != nil {
		util.ErrorMessage(fmt.Sprint("Error adding job:", err))
	} else {
		util.SuccessMessage(fmt.Sprint("Job added successfully:", job.Name))
	}
}
