package jobs

import (
	"testing"
	"time"
)

type testJob struct {
	run *bool
}

func (t testJob) Run() {
	*t.run = true
}

func TestCanAddSingleJobOnStoppedJobRunner(t *testing.T) {
	r := NewRunner()
	defer r.Destroy()

	err := r.AddJob("My job", scheduledAt(time.Now()), testJob{}, false)

	if err != nil {
		t.Errorf("Could not add job, got error %s", err)
	}
}

func TestCannotAddTwoJobsWithSameNameOnStoppedJobRunner(t *testing.T) {
	r := NewRunner()
	defer r.Destroy()

	err := r.AddJob("My job", scheduledAt(time.Now()), testJob{}, false)
	if err != nil {
		t.Errorf("Could not add job, got error %s", err)
	}
	err = r.AddJob("My job", scheduledAt(time.Now().Add(time.Duration(time.Second*1))), testJob{}, false)

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

func TestStuff(t *testing.T) {
	r := NewRunnerWithConcurrentExecutors(2)
	r.Start()
	defer r.Destroy()
	now := time.Now()
	job1run := false
	job2run := false
	job1 := testJob{run: &job1run}
	job2 := testJob{run: &job2run}
	r.RunJobAt(now.Add(time.Millisecond*10), job1)
	r.RunJobAt(now.Add(time.Millisecond*15), job2)

	<-time.After(time.Duration(time.Millisecond * 20))
	r.Stop()

	if !job1run {
		t.Error("Expected job 1 to have run")
	}

	if !job2run {
		t.Error("Expected job 2 to have run")
	}
}
