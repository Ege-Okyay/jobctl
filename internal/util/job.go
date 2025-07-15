package util

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/logger"
	"github.com/Ege-Okyay/jobctl/internal/types"
)

func RunJob(job types.JobConfig, started time.Time) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		psArgs := []string{"-NoProfile", "-Command", job.Command}
		if job.Timeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(job.Timeout)*time.Second)
			defer cancel()

			cmd = exec.CommandContext(ctx, "powershell.exe", psArgs...)
		} else {
			cmd = exec.Command("powershell.exe", psArgs...)
		}
	} else {
		shArgs := []string{"-c", job.Command}
		if job.Timeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(job.Timeout)*time.Second)
			defer cancel()

			cmd = exec.CommandContext(ctx, "sh", shArgs...)
		} else {
			cmd = exec.Command("sh", shArgs...)
		}
	}

	out, err := cmd.CombinedOutput()
	success := err == nil

	msg := fmt.Sprintf(
		"Scheduler: Job %q ran at %s Success=%v\n"+
			"Command: %s\n"+
			"Output:\n%s\n",
		job.Name,
		started.Format("2006-01-02 15:04:05"),
		success,
		job.Command,
		string(out),
	)

	if err != nil {
		msg += fmt.Sprintf("Error: %v\n", err)
	}

	logger.Log(msg)

	if dbErr := db.InsertRunRecord(job.Name, started, success, string(out)); dbErr != nil {
		logger.Log("Scheduler: Failed to insert run record for %q: %v", job.Name, dbErr)
	}

	next := started.Add(time.Duration(job.Interval) * time.Second)
	if updateErr := db.UpdateNextRun(job.Name, next); updateErr != nil {
		logger.Log("Scheduler: Failed to update next_run for %q: %v", job.Name, updateErr)
	}
}
