package cmd

import (
	"fmt"

	"github.com/Ege-Okyay/jobctl/internal/logic"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var DeleteCommand types.Command

func init() {
	DeleteCommand = types.Command{
		Name:        "delete",
		Description: "Delete a scheduled job",
		Usage:       "delete <job_name>",
		Flags:       []types.Flag{},
		Execute:     deleteHandler,
	}
}

func deleteHandler(args []string) {
	if len(args) != 1 {
		util.ErrorMessage(fmt.Sprint("Usage: ", DeleteCommand.Usage))
		return
	}

	name := args[0]
	cfgPath := util.ResolvePaths().ConfigPath

	if err := logic.DeleteJob(name, cfgPath); err != nil {
		util.ErrorMessage(fmt.Sprintf("Failed to delete %q: %v", name, err))
	} else {
		util.SuccessMessage(fmt.Sprintf("Job %q deleted", name))
	}
}
