package jobs

import (
	"testing"
	"time"
)

func TestConstantDelayScheduleOnOddTime(t *testing.T) {
	s := Every(time.Duration(time.Second * 2))

	testTime := time.Date(2000, 11, 1, 1, 1, 3, int(time.Millisecond)*300, time.UTC)

	delayed := s.Next(testTime)

	difference := delayed.Sub(testTime)

	expectedDifference := time.Duration(time.Millisecond * 1700)

	if difference != expectedDifference {
		t.Errorf("Expected a difference of %v, difference is %v", expectedDifference, difference)
	}
}

func TestConstantDelayScheduleOntimeToThePoint(t *testing.T) {
	s := Every(time.Duration(time.Second * 1))

	testTime := time.Date(2000, 11, 1, 1, 1, 3, 0, time.UTC)

	delayed := s.Next(testTime)

	difference := delayed.Sub(testTime)

	expectedDifference := time.Duration(time.Second)

	if difference != expectedDifference {
		t.Errorf("Expected a difference of %v, difference is %v", expectedDifference, difference)
	}
}

func TestConstantDelayScheduleWithLessThanOneSecondDelay(t *testing.T) {
	s := Every(time.Duration(time.Millisecond * 100))

	testTime := time.Date(2000, 11, 1, 1, 1, 3, 0, time.UTC)

	delayed := s.Next(testTime)

	difference := delayed.Sub(testTime)

	expectedDifference := time.Duration(time.Second)

	if difference != expectedDifference {
		t.Errorf("Expected a difference of %v, difference is %v", expectedDifference, difference)
	}
}
