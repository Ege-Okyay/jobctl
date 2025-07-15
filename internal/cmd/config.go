package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Ege-Okyay/jobctl/internal/config"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var ConfigCommand types.Command

func init() {
	ConfigCommand = types.Command{
		Name:        "config",
		Description: "Manage the jobctl config file",
		Usage:       "config [show|set <path>|validate|edit|reset]",
		Flags:       []types.Flag{},
		Execute:     configHandler,
	}
}

func configHandler(args []string) {
	if len(args) == 0 {
		fmt.Println(ConfigCommand.Usage)
		return
	}

	switch args[0] {
	case "show":
		path, err := config.ConfigPath()
		if err != nil {
			util.ErrorMessage(fmt.Sprint("Error:", err))
		} else {
			util.SuccessMessage(fmt.Sprint("Config path: ", path))
		}

	case "set":
		if len(args) < 2 {
			fmt.Println("Usage: jobctl config set <path>")
			return
		}

		newPath := args[1]

		if err := config.SetConfigPath(newPath); err != nil {
			util.ErrorMessage(fmt.Sprint("Failed to set config path:", err))
		} else {
			util.SuccessMessage(fmt.Sprint("Config path updated to:", newPath))
		}

	case "validate":
		path, err := config.ConfigPath()
		if err != nil {
			util.ErrorMessage(fmt.Sprint("Error location config:", err))
			return
		}

		if _, err := config.LoadConfig(path); err != nil {
			util.ErrorMessage(fmt.Sprint("Invalid config:", err))
		} else {
			util.SuccessMessage("Config is valid")
		}

	case "edit":
		path, err := config.ConfigPath()
		if err != nil {
			util.ErrorMessage(fmt.Sprint("Error locating config:", err))
			return
		}

		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("notepad.exe", path)
		case "darwin":
			cmd = exec.Command("open", "-a", "TextEdit", path)
		default:
			// Respect the EDITOR environment variable on Linux/Unix.
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "nano" // Fallback editor.
			}

			cmd = exec.Command(editor, path)
		}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			util.ErrorMessage(fmt.Sprint("Failed to open editor:", err))
		}

	case "reset":
		// Revert to the default configuration path by removing the custom path marker.
		home, err := os.UserHomeDir()
		if err != nil {
			util.ErrorMessage(fmt.Sprintf("Error finding home directory: %v", err))
			return
		}

		marker := filepath.Join(home, config.GetUserConfigMarker())
		if err := os.Remove(marker); err != nil {
			if os.IsNotExist(err) {
				util.SuccessMessage("No custom config path was set")
			} else {
				util.ErrorMessage(fmt.Sprintf("Failed to reset config path: %v", err))
			}
		} else {
			util.SuccessMessage("Custom config path cleared; using default location")
		}

	default:
		fmt.Println("unknown subcommand:", args[0])
		fmt.Print(ConfigCommand.Usage)
	}
}
