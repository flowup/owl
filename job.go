package owl

import (
	"time"
	"sync"
	"os/exec"
	"fmt"
)

type Job interface {
}

type JobResult struct {
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

func Scheduler(jobs <- chan Job) <-chan JobResult {
	schedulerJobs := make(chan JobResult)
	var mutex = &sync.Mutex{}

	var command *exec.Cmd
	go func() {
		for {
			select {
			case <-jobs:
				if command != nil {
					command.Process.Kill()
				}
				mutex.Lock()
				command = exec.Command("bash", "-c", "echo \"Good job, Sir\"")
				out, _ := command.CombinedOutput()
				fmt.Print(string(out))
				command = nil
				schedulerJobs <- JobResult{}
				mutex.Unlock()
			}
		}
	}()
	return schedulerJobs

}
