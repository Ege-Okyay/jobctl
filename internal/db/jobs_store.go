package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/types"
)

func InsertJob(j types.JobConfig) error {
	_, err := DB.Exec(
		`INSERT INTO jobs(name, interval_sec, command, next_run, retries, timeout_sec, enabled)
		VALUES(?,?,?, datetime('now','+'||?||' seconds'),?,?,?)`,
		j.Name, j.Interval, j.Command, j.Interval,
		j.Retries, j.Timeout, j.Enabled,
	)

	return err
}

func DeleteJob(name string) error {
	res, err := DB.Exec(`DELETE FROM jobs WHERE name = ?`, name)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no job named %q", name)
	}

	return nil
}

func UpdateJob(j types.JobConfig) error {
	enabledInt := 0
	if j.Enabled {
		enabledInt = 1
	}

	res, err := DB.Exec(
		`UPDATE jobs SET interval_sec = ?, command = ?, retries = ?, timeout_sec = ?, enabled = ? WHERE name = ?`,
		j.Interval, j.Command, j.Retries, j.Timeout, enabledInt, j.Name,
	)
	if err != nil {
		return fmt.Errorf("updating job %q: %w", j.Name, err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no job named %q to update", j.Name)
	}

	return nil
}

func GetJob(name string) (*types.JobConfig, error) {
	row := DB.QueryRow(`
		SELECT name, interval_sec, command, retries, timeout_sec, enabled, next_run
		FROM jobs WHERE name = ?
	`, name)

	var j types.JobConfig
	var enabledInt int
	var nextRunStr string

	err := row.Scan(
		&j.Name,
		&j.Interval,
		&j.Command,
		&j.Retries,
		&j.Timeout,
		&enabledInt,
		&nextRunStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	j.Enabled = enabledInt != 0
	if t, err := time.Parse(time.RFC3339, nextRunStr); err == nil {
		j.NextRun = t
	}

	return &j, nil
}

func GetAllJobs() ([]types.JobConfig, error) {
	rows, err := DB.Query(
		`SELECT name, interval_sec, command, retries, timeout_sec, enabled, next_run FROM jobs`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []types.JobConfig
	for rows.Next() {
		var j types.JobConfig
		var enabledInt int

		if err := rows.Scan(&j.Name, &j.Interval, &j.Command,
			&j.Retries, &j.Timeout, &enabledInt, &j.NextRun,
		); err != nil {
			return nil, err
		}

		j.Enabled = enabledInt != 0
		out = append(out, j)
	}

	return out, nil
}

func IsEmpty() (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM jobs`).Scan(&count)

	return (err == nil && count == 0), err
}

func GetDueJobs(cutoff time.Time) ([]types.JobConfig, error) {
	rows, err := DB.Query(
		`SELECT name, interval_sec, command, retries, timeout_sec, enabled, next_run
		FROM jobs
		WHERE next_run <= ? AND enabled = 1`,
		cutoff.Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []types.JobConfig
	for rows.Next() {
		var j types.JobConfig
		var enabledInt int

		if err := rows.Scan(
			&j.Name, &j.Interval, &j.Command,
			&j.Retries, &j.Timeout, &enabledInt, &j.NextRun,
		); err != nil {
			return nil, err
		}

		j.Enabled = enabledInt != 0
		jobs = append(jobs, j)
	}

	return jobs, nil
}

func ToggleJobEnabled(name string, en bool) error {
	res, err := DB.Exec(
		`UPDATE jobs SET enabled = ? WHERE name = ?`,
		en, name,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no job named %q", name)
	}

	return nil
}

func GetNextJob() (*types.JobConfig, error) {
	row := DB.QueryRow(`
		SELECT name, interval_sec, command, next_run, retries, timeout_sec, enabled
		FROM jobs WHERE enabled = 1 ORDER BY next_run ASC LIMIT 1
	`)

	var j types.JobConfig
	var nextRunStr string
	var enabledInt int

	if err := row.Scan(
		&j.Name, &j.Interval, &j.Command,
		&nextRunStr, &j.Retries, &j.Timeout, &enabledInt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if t, err := time.Parse(time.RFC3339, nextRunStr); err == nil {
		j.NextRun = t
	}

	j.Enabled = enabledInt != 0

	return &j, nil
}
