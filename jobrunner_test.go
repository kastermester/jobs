package jobs

import (
	"testing"
	"time"
)

type testJob struct {
	run bool
}

func (t *testJob) Run() {
	t.run = true
}

func TestCanAddSingleJobOnStoppedJobRunner(t *testing.T) {
	r := NewRunner()
	defer r.Destroy()

	err := r.AddJob("My job", scheduledAt(time.Now()), &testJob{}, false)

	if err != nil {
		t.Errorf("Could not add job, got error %s", err)
	}
}

func TestCannotAddTwoJobsWithSameNameOnStoppedJobRunner(t *testing.T) {
	r := NewRunner()
	defer r.Destroy()

	err := r.AddJob("My job", scheduledAt(time.Now()), &testJob{}, false)
	if err != nil {
		t.Errorf("Could not add job, got error %s", err)
	}
	err = r.AddJob("My job", scheduledAt(time.Now().Add(time.Duration(time.Second*1))), &testJob{}, false)

	if err == nil {
		t.Error("Expected to get an error, got nothing")
	} else if e, ok := err.(ErrAddNameAlreadyExists); !ok {
		t.Error("Expected error to be of type ErrAddNameAlreadyExists, it was not, message was %s", e)
	}
}

func TestCannotRemoveNonExistingJob(t *testing.T) {
	r := NewRunner()
	defer r.Destroy()

	err := r.RemoveJob("Non job")

	if err == nil {
		t.Error("Expected to get an error")
	} else if e, ok := err.(ErrRemoveNameNotFound); !ok {
		t.Error("Expected error to be of type ErrRemoveNameNotFound, it was not, message was %s", e)
	}
}

func TestRunsJobs(t *testing.T) {
	r := NewRunnerWithConcurrentExecutors(2)
	r.Start()
	defer r.Destroy()
	now := time.Now()
	job1 := testJob{false}
	job2 := testJob{false}
	r.RunJobAt(now.Add(time.Millisecond*10), &job1)
	r.RunJobAt(now.Add(time.Millisecond*15), &job2)

	<-time.After(time.Duration(time.Millisecond * 20))
	r.Stop()

	if !job1.run {
		t.Error("Expected job 1 to have run")
	}

	if !job2.run {
		t.Error("Expected job 2 to have run")
	}
}

type jobWithData struct {
	data int
	run  *bool
}

func (j *jobWithData) Run() {
	*j.run = true
	j.data = j.data + 1
}

const startValue = 2

func TestCanCreateJobsThatCarryData(t *testing.T) {
	r := NewRunner()
	r.Start()
	defer r.Destroy()
	run := false
	job := jobWithData{startValue, &run}
	r.RunJobAt(time.Now(), &job)

	<-time.After(10 * time.Millisecond)
	r.Stop()

	if !run {
		t.Error("Expected job to have run")
	}

	if startValue+1 != job.data {
		t.Error("Expected data to have been incremented by one")
	}
}
