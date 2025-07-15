package cmd

import (
	"strings"

	"github.com/Ege-Okyay/jobctl/internal/logger"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var DebugCommand types.Command

func init() {
	DebugCommand = types.Command{
		Name:        "debug",
		Description: "Enable or disable on-screen debug logs",
		Usage:       "debug <on|off> # show or suppress debug logs on terminal",
		Flags:       []types.Flag{},
		Execute:     debugHandler,
	}
}

func debugHandler(args []string) {
	if len(args) != 1 {
		util.ErrorMessage("Error: expected exactly one argument: on or off")
		return
	}

	arg := strings.ToLower(args[0])
	switch arg {
	case "on":
		logger.DebugEnabled = true
		util.SuccessMessage("Debug mode: ON")
	case "off":
		logger.DebugEnabled = false
		util.ErrorMessage("Debug mode: OFF")
	default:
		util.ErrorMessage("Error: argument must be 'on' or 'off'")
	}
}
