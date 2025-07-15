package runner

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/logger"
	"github.com/Ege-Okyay/jobctl/internal/types"
)

// RunScheduler is the main loop for the job scheduler.
// It ticks every 5 seconds and runs any jobs that are due.
func RunScheduler(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	logger.Log("Scheduler: started")

	for {
		select {
		case <-ctx.Done():
			logger.Log("Scheduler: stopped")
			return
		case <-ticker.C:
			paused, _ := db.IsPaused()
			if paused {
				continue
			}
			runDueJobs()
		}
	}
}

// runDueJobs fetches and executes all jobs that are scheduled to run.
func runDueJobs() {
	jobs, err := db.GetDueJobs(time.Now())
	if err != nil {
		fmt.Println("Error fetching due jobs:", err)
		return
	}

	for _, job := range jobs {
		j := job
		go func() {
			RunJobCmd(j, time.Now(), true)
		}()
	}
}

// RunJobCmd executes a single job command, captures its output, and records the result.
// It supports timeouts and can optionally bump the job's next run time.
func RunJobCmd(job types.JobConfig, started time.Time, bumpNext bool) error {
	var cmd *exec.Cmd
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell, flag = "powershell.exe", "-Command"
	} else {
		shell, flag = "sh", "-c"
	}

	// If a timeout is set, create a context with a timeout.
	if job.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(job.Timeout)*time.Second)
		defer cancel()

		cmd = exec.CommandContext(ctx, shell, flag, job.Command)
	} else {
		cmd = exec.Command(shell, flag, job.Command)
	}

	outBytes, err := cmd.CombinedOutput()
	success := err == nil

	var sb strings.Builder
	ts := started.Format("2006-01-02 15:04:05")
	status := "SUCCESS"
	if !success {
		status = "FAILURE"
	}

	sb.WriteString(fmt.Sprintf("=== JOB %q @ %s ===\n", job.Name, ts))
	sb.WriteString(fmt.Sprintf("Command: %s\n", job.Command))
	sb.WriteString(fmt.Sprintf("Status : %s\n", status))

	out := strings.TrimRight(string(outBytes), "\n")
	if out == "" {
		sb.WriteString("Output : <none>\n")
	} else {
		sb.WriteString("Output :\n")
		for _, line := range strings.Split(out, "\n") {
			sb.WriteString("\t" + line + "\n")
		}
	}

	if err != nil {
		sb.WriteString(fmt.Sprintf("Error\t: %v\n", err))
	}

	msg := sb.String()

	logger.Log(msg)

	if dbErr := db.InsertRunRecord(job.Name, started, success, out); dbErr != nil {
		logger.Log("! Failed to record run: %v", dbErr)
	}

	// If bumpNext is true, update the job's next run time.
	if bumpNext {
		next := started.Add(time.Duration(job.Interval) * time.Second)
		if updErr := db.UpdateNextRun(job.Name, next); updErr != nil {
			logger.Log("! Failed to update next_run: %v", updErr)
		}
	}

	return err
}
