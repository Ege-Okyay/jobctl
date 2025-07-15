package cmd

import (
	"fmt"

	"github.com/Ege-Okyay/jobctl/internal/logic"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var DisableCommand types.Command

func init() {
	DisableCommand = types.Command{
		Name:        "disable",
		Description: "Disable a scheduled job",
		Usage:       "disable <job_name>",
		Flags:       []types.Flag{},
		Execute:     disableHandler,
	}
}

func disableHandler(args []string) {
	if len(args) != 1 {
		util.ErrorMessage(fmt.Sprint("Usage: ", DisableCommand.Usage))
		return
	}

	name := args[0]
	cfgPath := util.ResolvePaths().ConfigPath

	if err := logic.DisableJob(name, cfgPath); err != nil {
		util.ErrorMessage(fmt.Sprintf("Failed to disable %q: %v", name, err))
	} else {
		util.SuccessMessage(fmt.Sprintf("Job %q disabled", name))
	}
}
