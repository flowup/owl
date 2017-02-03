package owl

import (
	"time"
)

// Job interface implements method Start, and Stop.
// Start runs command and Stop sends kill signal
type Job interface {
	Start() error
	Stop() error
}

// Debounce debounces the channel of jobs by the given amount of
// milliseconds
// ToDO: the for cycle will never stop until the application stops. May be potential leak in a future
func Debounce(jobs <-chan Job, amount int64) <-chan Job {
	debouncedJobs := make(chan Job)

	// cache for Job to be debounced
	var cache Job

	go func() {
		for {
			cache = <-jobs
			time.Sleep(time.Duration(amount) * time.Millisecond)
			// draining the channel
			draining := true
			for draining {
				select {
				case cache = <-jobs:
					continue
				default:
					draining = false
				}
			}
			debouncedJobs <- cache
		}
	}()
	return debouncedJobs
}

// Scheduler continually runs jobs read from the jobs channel. If any job is running within
// the scheduler, it will be killed and replaced by the next job
func Scheduler(jobs <- chan Job) <-chan error {
	schedulerJobs := make(chan error)

	var runningJob Job

	go func() {
		for {
			job := <-jobs
			// check if another command is running
			if runningJob != nil {
				runningJob.Stop()
			}

			// swap jobs
			runningJob = job

			go func() {
				// executing
				schedulerJobs <- runningJob.Start()
			}()
		}
	}()

	return schedulerJobs
}
