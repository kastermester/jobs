package jobs

import (
	"time"
)

// Schedule represents the schedule of a job, can be a run-once type job, or scheduled to run on certain intervals.
type Schedule interface {
	// Next returns the next activation time - that a job on this schedule should run on, that is later than the time provided.
	Next(after time.Time) time.Time
}
