package logic

import (
	"fmt"
	"os"

	"github.com/Ege-Okyay/jobctl/internal/config"
	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/types"
)

// SyncDBWithConfig synchronizes the database with the configuration file.
// It ensures that the database reflects the state of the config file by
// adding, updating, and deleting jobs as necessary.
func SyncDBWithConfig(configPath string) error {
	conf, err := config.LoadConfig(configPath)
	if err != nil {
		// If the config file doesn't exist, create it from the database.
		if os.IsNotExist(err) {
			dbJobs, dbErr := db.GetAllJobs()
			if dbErr != nil {
				return fmt.Errorf("reading DB for bootstrap: %w", dbErr)
			}

			if saveErr := config.SaveConfig(configPath, &types.Config{Jobs: dbJobs}); saveErr != nil {
				return fmt.Errorf("bootstrap writing %s: %w", configPath, saveErr)
			}

			return nil
		}

		return fmt.Errorf("loading config %s: %w", configPath, err)
	}

	dbJobs, err := db.GetAllJobs()
	if err != nil {
		return fmt.Errorf("loading jobs from DB: %w", err)
	}

	cfgMap := make(map[string]types.JobConfig, len(conf.Jobs))
	for _, j := range conf.Jobs {
		cfgMap[j.Name] = j
	}

	dbMap := make(map[string]types.JobConfig, len(dbJobs))
	for _, j := range dbJobs {
		dbMap[j.Name] = j
	}

	// Delete jobs from the database that are no longer in the config file.
	for name := range dbMap {
		if _, ok := cfgMap[name]; !ok {
			if err := db.DeleteJob(name); err != nil {
				return fmt.Errorf("deleting %q: %w", name, err)
			}
		}
	}

	// Add new jobs from the config file to the database.
	for name, cfgJobs := range cfgMap {
		if _, ok := dbMap[name]; !ok {
			if err := AddJob(cfgJobs, configPath); err != nil {
				return fmt.Errorf("adding %q: %w", name, err)
			}
		}
	}

	// Update existing jobs in the database if they have changed in the config file.
	for name, cfgJob := range cfgMap {
		if dbJobs, ok := dbMap[name]; ok && needsUpdate(dbJobs, cfgJob) {
			if err := db.UpdateJob(cfgJob); err != nil {
				return fmt.Errorf("updating %q: %w", name, err)
			}
		}
	}

	return nil
}

// needsUpdate checks if a job's configuration has changed.
func needsUpdate(dbJob, cfgJob types.JobConfig) bool {
	return dbJob.Interval != cfgJob.Interval ||
		dbJob.Command != cfgJob.Command ||
		dbJob.Retries != cfgJob.Retries ||
		dbJob.Timeout != cfgJob.Timeout ||
		dbJob.Enabled != cfgJob.Enabled
}
