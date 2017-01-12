package owl

import (
	"time"
)

type Job interface {

}


// Debounce debounces the channel of jobs by the given amount of
// milliseconds
func Debounce(jobs <-chan Job, amount int64) <-chan Job {
	debouncedJobs := make(chan Job)

	// cache for Job to be debounced
	var cache Job

	go func() {
		for {
			select {
			case cache = <-jobs:
				continue
			case <-time.After(time.Duration(amount) * time.Millisecond):
				// if we don't have cache, just continue
				if cache == nil {
					continue
				}

				debouncedJobs <- cache
				cache = nil
			}
		}
	}()

	return debouncedJobs
}