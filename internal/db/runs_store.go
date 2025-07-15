package db

import (
	"database/sql"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/types"
)

func InsertRunRecord(jobName string, ts time.Time, success bool, output string) error {
	_, err := DB.Exec(
		`INSERT INTO job_runs (job_id, ts, success, output)
		VALUES(
			(SELECT id FROM jobs WHERE name = ?),
			?, ?, ?
		)`,
		jobName,
		ts.Format(time.RFC3339),
		success,
		output,
	)

	return err
}

func UpdateNextRun(jobName string, next time.Time) error {
	_, err := DB.Exec(
		`UPDATE jobs SET next_run = ? WHERE name = ?`,
		next.Format(time.RFC3339),
		jobName,
	)

	return err
}

func GetLastRun(jobName string) (*types.RunRecord, error) {
	row := DB.QueryRow(`
		SELECT jr.id, j.name, jr.ts, jr.success, jr.output
		FROM job_runs jr
		JOIN jobs j ON jr.job_id = j.id
		WHERE j.name = ?
		ORDER BY jr.ts DESC
		LIMIT 1
	`, jobName)

	var rec types.RunRecord
	var ts string

	if err := row.Scan(&rec.ID, &rec.JobName, &ts, &rec.Success, &rec.Output); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	rec.Timestamp, _ = time.Parse(time.RFC3339, ts)

	return &rec, nil
}

func GetRunHistory(jobName string, limit int) ([]types.RunRecord, error) {
	rows, err := DB.Query(`
		SELECT jr.id, j.name, jr.ts, jr.success, jr.output
		FROM job_runs jr
		JOIN jobs j ON jr.job_id = j.id
		WHERE j.name = ?
		ORDER BY jr.ts DESC
		LIMIT ?
	`, jobName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []types.RunRecord
	for rows.Next() {
		var rec types.RunRecord
		var ts string

		if err := rows.Scan(&rec.ID, &rec.JobName, &ts, &rec.Success, &rec.Output); err != nil {
			return nil, err
		}

		rec.Timestamp, _ = time.Parse(time.RFC3339, ts)
		history = append(history, rec)
	}

	return history, nil
}
