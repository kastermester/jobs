package jobs

import "time"

// ConstantDelaySchedule is a simple schedule that will schedule a Job to run with a constant delay, specified by delay.
type ConstantDelaySchedule struct {
	// The delay to schedule the Job by
	delay time.Duration
}

// Next rounds after time down towards the whole second (that is, subtract after.NanoSecond())
// and adds s.Delay() to it.
func (s ConstantDelaySchedule) Next(after time.Time) time.Time {
	delay := s.delay
	return after.Add(delay - time.Duration(after.Nanosecond())*time.Nanosecond)
}

// Delay returns the delay that this ConstantDelaySchedule represents.
func (s ConstantDelaySchedule) Delay() time.Duration {
	return s.delay
}

// Every is a constructor function for ConstantDelaySchedule.
// If the delay specified is less than a second, a delay of one Second will be created.
func Every(delay time.Duration) ConstantDelaySchedule {
	if delay < time.Duration(time.Second) {
		delay = time.Duration(time.Second)
	}
	return ConstantDelaySchedule{delay}
}
