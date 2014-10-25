package jobs

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type entries []*Entry

// The name format string for jobs that are scheduled to run at a specific time.
// Please avoid using names that could clash with this
const RunOnceJobFormat = "Run once job (#%d)"

// Type for the JobRunner
type JobRunner struct {
	executors        uint
	runOnceJobNumber uint32
	entries          entries
	stopExecutor     chan struct{}
	schedulerStop    chan struct{}
	starting         chan struct{}
	add              chan entryAndErrorChannel
	remove           chan jobNameAndErrorChannel
	mutex            sync.Mutex
	snapshot         chan entries
	execute          chan *Entry
	running          bool
	destroyed        bool
}

// An Entry is a type with a Schedule and a Job, along with other information regarding
// When to run the job next, when it was previously run and whether or not it should only run once
type Entry struct {
	Schedule Schedule
	Next     time.Time
	Prev     time.Time
	Job      Job
	Name     string
	Once     bool
}

type byTime entries

func (b byTime) Len() int {
	return len(b)
}

func (b byTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byTime) Less(i, j int) bool {
	// We schedule time.Zero as the last element
	// we need this as we set the Next of the entry to time.Zero when it is still executing
	if b[i].Next.IsZero() {
		return false
	}

	if b[j].Next.IsZero() {
		return true
	}
	return b[i].Next.Before(b[j].Next)
}

type entryAndErrorChannel struct {
	entry *Entry
	error chan<- error
}

type jobNameAndErrorChannel struct {
	jobName string
	error   chan<- error
}

// Constructs a new JobRunner with support for 1 concurrent execution at a time
// The runner must be started before it starts running jobs
func NewRunner() *JobRunner {
	return NewRunnerWithConcurrentExecutors(1)
}

// Constructs a new JobRunner with support for the supplied number of concurrent execution at a time
// The runner must be started before it starts running jobs
// Panics if 0 is provided
func NewRunnerWithConcurrentExecutors(concurrentExecutions uint) *JobRunner {
	if concurrentExecutions == 0 {
		panic("Must provide a non zero value for concurrentExecutions")
	}
	r := &JobRunner{
		destroyed:        false,
		running:          false,
		executors:        concurrentExecutions,
		execute:          make(chan *Entry),
		runOnceJobNumber: 0,
		entries:          nil,
		add:              make(chan entryAndErrorChannel),
		remove:           make(chan jobNameAndErrorChannel),
		stopExecutor:     make(chan struct{}),
		schedulerStop:    make(chan struct{}),
		snapshot:         make(chan entries),
	}
	return r
}

func (e entries) pos(name string) int {
	for p, e := range e {
		if e.Name == name {
			return p
		}
	}
	return -1
}

// Adds a job to the JobRunner. This is a very low level API, prefer to use one of the other proxy methods
func (r *JobRunner) AddJob(name string, s Schedule, j Job, once bool) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.destroyed {
		return ErrJobRunnerDestroyed{}
	}
	entry := &Entry{
		Name:     name,
		Schedule: s,
		Job:      j,
		Once:     once,
	}

	if r.running {
		ch := make(chan error)
		ee := entryAndErrorChannel{
			entry: entry,
			error: ch,
		}

		r.add <- ee

		return <-ch
	} else {
		return r.appendEntry(entry)
	}
}

// Adds a job that will run at time t - and only be run once
func (r *JobRunner) RunJobAt(t time.Time, j Job) error {
	n := atomic.AddUint32(&r.runOnceJobNumber, 1)
	name := fmt.Sprintf(RunOnceJobFormat, n)

	return r.AddJob(name, scheduledAt(t), j, true)
}

// Adds a function that will run at time t - and only be run once
func (r *JobRunner) RunFuncAt(t time.Time, f func()) error {
	return r.RunJobAt(t, NewFuncJob(f))
}

// Removes a named job from the list - preventing it from being scheduled for further execution.
// However if the entry has been scheduled to run, but not yet done so - it will still complete that execution.
func (r *JobRunner) RemoveJob(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.destroyed {
		return ErrJobRunnerDestroyed{}
	}

	if r.running {
		ch := make(chan error)
		je := jobNameAndErrorChannel{
			jobName: name,
			error:   ch,
		}

		r.remove <- je
		return <-ch
	} else {
		return r.removeEntryWithName(name)
	}
}

func (r *JobRunner) appendEntry(e *Entry) error {
	i := r.entries.pos(e.Name)
	if i != -1 {
		return ErrAddNameAlreadyExists{
			name: e.Name,
		}
	}
	e.Next = e.Schedule.Next(time.Now().Local())
	r.entries = append(r.entries, e)
	return nil
}

func (r *JobRunner) removeEntryWithName(name string) error {
	i := r.entries.pos(name)
	if i == -1 {
		return ErrRemoveNameNotFound{
			name: name,
		}
	}

	r.entries = r.entries[:i+copy(r.entries[i:], r.entries[i+1:])]
	return nil
}

// Starts the JobRunner. This is an illegal operation if the JobRunner has been Destroyed with a call to Destroy()
// If the JobRunner is already running, this is a no-op
func (r *JobRunner) Start() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.destroyed {
		return ErrJobRunnerDestroyed{}
	}

	if r.running {
		return nil
	}

	r.running = true
	go r.scheduleJobs()

	// Start all the executors
	for i := 0; i < int(r.executors); i++ {
		go r.executeJobs()
	}
	return nil
}

// Returns whether or not the JobRunner is running
func (r *JobRunner) IsRunning() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.running
}

func (r *JobRunner) stop() {
	if !r.running {
		return
	}

	// Stop the scheduler - ensure no new jobs are being scheduled
	r.schedulerStop <- struct{}{}

	// Stop all the executors
	for i := 0; i < int(r.executors); i++ {
		r.stopExecutor <- struct{}{}
	}

	r.running = false
}

// Stops the JobRunner. Once stopped, the JobRunner can be started again
func (r *JobRunner) Stop() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.destroyed {
		return ErrJobRunnerDestroyed{}
	}

	r.stop()

	return nil
}

// Destroys the JobRunner, after this is done, no further methods may be called on this instance of the JobRunner
func (r *JobRunner) Destroy() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.destroyed {
		return
	}
	// Make sure we're in the stopped state
	r.stop()

	r.entries = nil
	r.destroyed = true
	close(r.execute)
	close(r.add)
	close(r.remove)
	close(r.stopExecutor)
	close(r.schedulerStop)
	close(r.snapshot)
}

func (r *JobRunner) executeJobs() {
	for {
		select {
		case e := <-r.execute:
			e.Job.Run()
		case <-r.stopExecutor:
			return
		}
	}
}

func (r *JobRunner) scheduleJobs() {
	now := time.Now().Local()
	for _, entry := range r.entries {
		entry.Next = entry.Schedule.Next(now)
	}

	for {
		// Sort the entries by time
		sort.Sort(byTime(r.entries))
		// Figure out for how long we should wait
		var effective time.Time
		if len(r.entries) == 0 {
			// Wait for a long time, we will get interrupted by other channel operations
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = r.entries[0].Next
		}

		waitFor := effective.Sub(time.Now().Local())
		if waitFor < 0 {
			waitFor = time.Duration(0)
		}

		select {
		case now = <-time.After(waitFor):
			// There might be several entries that we should execute now
			for _, e := range r.entries {
				if e.Next != effective {
					break
				}
				// Send job off to be executed
				e.Prev = now
				e.Next = e.Schedule.Next(now)
				if e.Once {
					if err := r.removeEntryWithName(e.Name); err != nil {
						panic(err)
					}
				}
				r.execute <- e
			}
		case e := <-r.add:
			e.error <- r.appendEntry(e.entry)
			continue
		case j := <-r.remove:
			j.error <- r.removeEntryWithName(j.jobName)
			continue
		case <-r.schedulerStop:
			return
		}
	}
}