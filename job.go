package jobs

// A Job is the basic interface any jobs should implement
type Job interface {
	// Runs the job. Any kind of error handling should be taken care of inside the job
	Run()
}

type funcJob func()

func (f funcJob) Run() {
	f()
}

// Takes a function and returns a job that runs that function
func NewFuncJob(f func()) Job {
	return funcJob(f)
}
