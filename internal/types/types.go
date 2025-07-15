package types

import "time"

// Command defines the structure for a CLI command.
type Command struct {
	Name        string
	Description string
	Usage       string
	Flags       []Flag
	Execute     func(args []string)
}

// Flag defines the structure for a command-line flag.
type Flag struct {
	Name        string
	Description string
	Required    bool
}

// CommandDistance is used for finding the closest command name.
type CommandDistance struct {
	Name  string
	Score int
}

// JobConfig represents the configuration for a single job.
type JobConfig struct {
	Name     string `toml:"name"`
	Interval int    `toml:"interval"`
	Command  string `toml:"command"`
	Retries  int    `toml:"retries"`
	Timeout  int    `toml:"timeout"`
	Enabled  bool   `toml:"enabled"`

	// NextRun is the next scheduled run time for the job.
	// It is not stored in the config file, hence the `toml:"-"` tag.
	NextRun time.Time `toml:"-"`
}

// RunRecord represents a single execution of a job.
type RunRecord struct {
	ID        int
	JobName   string
	Timestamp time.Time
	Success   bool
	Output    string
}

// Config represents the top-level structure of the TOML config file.
type Config struct {
	Jobs []JobConfig `toml:"job"`
}

// AppPaths holds the absolute paths to key application files and directories.
type AppPaths struct {
	DBPath     string
	ConfigPath string
	LogDir     string
	MigDir     string
}
