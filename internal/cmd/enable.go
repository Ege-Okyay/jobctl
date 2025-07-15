package cmd

import (
	"fmt"

	"github.com/Ege-Okyay/jobctl/internal/logic"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var EnableCommand types.Command

func init() {
	EnableCommand = types.Command{
		Name:        "enable",
		Description: "Enable a scheduled job",
		Usage:       "enable <job_name>",
		Flags:       []types.Flag{},
		Execute:     enableHandler,
	}
}

func enableHandler(args []string) {
	if len(args) != 1 {
		util.ErrorMessage(fmt.Sprint("Usage: ", EnableCommand.Usage))
		return
	}

	name := args[0]
	cfgPath := util.ResolvePaths().ConfigPath

	if err := logic.EnableJob(name, cfgPath); err != nil {
		util.ErrorMessage(fmt.Sprintf("Failed to enable %q: %v", name, err))
	} else {
		util.SuccessMessage(fmt.Sprintf("Job %q enabled", name))
	}
}
