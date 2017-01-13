package owl

import (
	"os/exec"
	"fmt"
	"sync"
	"time"
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
			cache = <-jobs
			time.Sleep(time.Duration(amount))
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

// Scheduler run the newest job and killed old one
func Scheduler(jobs <- chan Job, run string) <-chan JobResult {
	schedulerJobs := make(chan JobResult)
	var mutex = &sync.Mutex{}
	var command *exec.Cmd
	go func() {
		for {
			select {
			case <-jobs:
				mutex.Lock()

				// check if another command is running
				// if is, kill it
				if command != nil {
					if command.Process != nil {
						command.Process.Kill()
					}
				}
				go func() {
					mutex.Unlock()
					// executing
					command = exec.Command("bash", "-c", run)
					out, _ := command.CombinedOutput()
					fmt.Print(string(out))
					schedulerJobs <- JobResult{}
				}()
			}
		}
	}()
	return schedulerJobs
}
