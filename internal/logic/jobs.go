package logic

import (
	"fmt"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/config"
	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/runner"
	"github.com/Ege-Okyay/jobctl/internal/types"
)

// RunJob executes a job immediately.
func RunJob(name string, keepRemaining bool) error {
	var job *types.JobConfig

	all, err := db.GetAllJobs()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, j := range all {
		if j.Name == name {
			job = &j
			break
		}
	}

	if job == nil {
		return fmt.Errorf("no job named %q", name)
	}

	if !job.Enabled {
		return fmt.Errorf("job %q is disabled", name)
	}

	if err := runner.RunJobCmd(*job, time.Now(), !keepRemaining); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// AddJob adds a new job to the database and syncs the config file.
func AddJob(j types.JobConfig, configPath string) error {
	if err := db.InsertJob(j); err != nil {
		return fmt.Errorf("%w", err)
	}

	all, err := db.GetAllJobs()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := config.SyncJobs(configPath, all); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// EditJob updates an existing job in the database and syncs the config file.
func EditJob(name string, updates types.JobConfig, configPath string) error {
	job, err := db.GetJob(name)
	if err != nil {
		return fmt.Errorf("fetching job: %w", err)
	}

	if job == nil {
		return fmt.Errorf("no job named %q", name)
	}

	merged := MergeJob(*job, updates)

	if err := db.UpdateJob(merged); err != nil {
		return fmt.Errorf("%w", err)
	}

	all, err := db.GetAllJobs()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := config.SyncJobs(configPath, all); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// EnableJob enables a job in the database and syncs the config file.
func EnableJob(name, configPath string) error {
	if err := db.ToggleJobEnabled(name, true); err != nil {
		return fmt.Errorf("%w", err)
	}

	all, err := db.GetAllJobs()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := config.SyncJobs(configPath, all); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// DisableJob disables a job in the database and syncs the config file.
func DisableJob(name, configPath string) error {
	if err := db.ToggleJobEnabled(name, false); err != nil {
		return fmt.Errorf("%w", err)
	}

	all, err := db.GetAllJobs()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := config.SyncJobs(configPath, all); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// DeleteJob deletes a job from the database and syncs the config file.
func DeleteJob(name, configPath string) error {
	if err := db.DeleteJob(name); err != nil {
		return fmt.Errorf("%w", err)
	}

	all, err := db.GetAllJobs()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := config.SyncJobs(configPath, all); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// MergeJob merges the non-zero fields of an update into an original job config.
func MergeJob(orig, upd types.JobConfig) types.JobConfig {
	out := orig
	if upd.Interval > 0 {
		out.Interval = upd.Interval
	}
	if upd.Command != "" {
		out.Command = upd.Command
	}
	if upd.Retries >= 0 {
		out.Retries = upd.Retries
	}
	if upd.Timeout >= 0 {
		out.Timeout = upd.Timeout
	}

	return out
}
