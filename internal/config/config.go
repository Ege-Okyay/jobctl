package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/Ege-Okyay/jobctl/internal/types"
)

// LoadConfig reads and decodes a TOML configuration file into a Config struct.
func LoadConfig(path string) (*types.Config, error) {
	var conf types.Config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", path)
	}

	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

// SaveConfig encodes a Config struct into a TOML file.
func SaveConfig(path string, conf *types.Config) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to open config file for writing: %w", err)
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(conf); err != nil {
		return fmt.Errorf("failed to encode config to TOML: %w", err)
	}

	return nil
}

// SyncJobs overwrites the configuration file with a new set of jobs.
func SyncJobs(path string, jobs []types.JobConfig) error {
	conf := types.Config{
		Jobs: jobs,
	}
	return SaveConfig(path, &conf)
}
