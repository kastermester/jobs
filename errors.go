package jobs

import "fmt"

// ErrRemoveNameNotFound is an error returned when trying to remove a job with a name that does not exist in the JobRunners entry list.
type ErrRemoveNameNotFound struct {
	name string
}

// Error formats the error.
func (e ErrRemoveNameNotFound) Error() string {
	return fmt.Sprintf("Could find job to remove, with name: %s", e.name)
}

// Name returns the name that could not be removed.
func (e ErrRemoveNameNotFound) Name() string {
	return e.name
}

// ErrAddNameAlreadyExists i an error returned when adding a job with a name that already exists in the JobRunners entry list.
type ErrAddNameAlreadyExists struct {
	name string
}

// Error Formats the error.
func (e ErrAddNameAlreadyExists) Error() string {
	return fmt.Sprintf("Could not add job with name %s, as it already exists", e.name)
}

// Name returns the name that could not be added.
func (e ErrAddNameAlreadyExists) Name() string {
	return e.name
}

// ErrJobRunnerDestroyed is an error returned when calling methods on a JobRunner that has been destroyed.
type ErrJobRunnerDestroyed struct{}

// Error formats the error.
func (e ErrJobRunnerDestroyed) Error() string {
	return "The job runner is destroyed"
}
