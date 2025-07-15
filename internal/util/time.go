package util

import (
	"time"

	"github.com/Ege-Okyay/jobctl/internal/db"
)

// AnchorTime returns the current time, unless the system is paused, in which case
// it returns the time the system was paused. This is used to ensure that
// job schedules are calculated relative to a consistent point in time.
func AnchorTime() time.Time {
	if paused, _ := db.IsPaused(); paused {
		if ts, err := db.GetPauseTimestamp(); err == nil && !ts.IsZero() {
			return ts
		}
	}
	return time.Now()
}
