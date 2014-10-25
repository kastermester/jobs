package jobs

import (
	"testing"
	"time"
)

var testTime = time.Date(2000, 10, 2, 10, 10, 20, 0, time.UTC)

func TestScheduledAtNextWhenTimeIsAfterArgument(t *testing.T) {
	tt := testTime.Add(time.Second * 10)
	s := scheduledAt(tt)

	tn := s.Next(testTime)

	if tn != tt {
		t.Errorf("Expected result to be %s, was %s", tt, tn)
	}
}

func TestScheduledAtNextWhenTimeIsEqual(t *testing.T) {
	s := scheduledAt(testTime)

	tn := s.Next(testTime)

	if tn != testTime {
		t.Errorf("Expected result to be %s, was %s", testTime, tn)
	}
}

func TestScheduledAtNextWhenTimeIsBeforeArgument(t *testing.T) {
	tt := testTime.Add(time.Second * -10)
	s := scheduledAt(tt)

	tn := s.Next(testTime)

	if tn != testTime {
		t.Errorf("Expected result to be %s, was %s", testTime, tn)
	}
}
