package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// userConfigMarker is the name of a file in the user's home directory
// that, if present, contains the absolute path to the config file.
const userConfigMarker = ".jobctl_config"

// ConfigPath determines the path to the jobs.toml configuration file.
// It prioritizes a custom path defined in the userConfigMarker file
// before falling back to a default location in the user's config directory.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err == nil {
		marker := filepath.Join(home, userConfigMarker)
		if data, err := os.ReadFile(marker); err == nil {
			return strings.TrimSpace(string(data)), nil
		}
	}

	// Fallback to the default config directory.
	cfgBase, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot find user config dir: %w", err)
	}

	cfgDir := filepath.Join(cfgBase, "jobctl")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		return "", fmt.Errorf("creating config dir %q: %w", cfgDir, err)
	}

	// Create a default config file if one doesn't exist.
	cfgPath := filepath.Join(cfgDir, "jobs.toml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		defaultContents := "# jobctl config\n\n# Add your [[job]] entries below\n"
		if err := os.WriteFile(cfgPath, []byte(defaultContents), 0644); err != nil {
			return "", fmt.Errorf("writing default config %q: %w", cfgPath, err)
		}
	}

	return cfgPath, nil
}

// SetConfigPath sets a custom path for the configuration file
// by writing the path to the userConfigMarker file.
func SetConfigPath(path string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.New("cannot determine home directory")
	}
	marker := filepath.Join(home, userConfigMarker)
	return os.WriteFile(marker, []byte(path), 0644)
}

// GetUserConfigMarker returns the name of the marker file.
func GetUserConfigMarker() string {
	return userConfigMarker
}
