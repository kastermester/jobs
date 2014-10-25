package jobs

import "time"

type scheduledAt time.Time

// Next will only return the scheduled time if it is not before "after"
// Otherwise after time is reported
func (s scheduledAt) Next(after time.Time) time.Time {
	if time.Time(s).Before(after) {
		return after
	}

	return time.Time(s)
}
