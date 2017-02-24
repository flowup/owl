package owl

import (
	"os/exec"
	"io"
	"bufio"
	"fmt"
	"syscall"
)

// WatcherJob is structure that represents job - command and state of executed command
type WatcherJob struct {
	Command string
	Cmd     *exec.Cmd
	OutPipe io.Writer
}

// Start is function that runs job's bash command
func (job*WatcherJob) Start() error {
	job.cmd = exec.Command("bash", "-c", job.command)

	stderr, err := job.cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := job.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	errscanner := bufio.NewScanner(stderr)
	outscanner := bufio.NewScanner(stdout)

	go func() {
		for {
			if outscanner.Scan() {
				fmt.Println(outscanner.Text())
			} else if errscanner.Scan() {
				fmt.Println(errscanner.Text())
			} else {
				break
			}
		}
	}()

	return job.cmd.Run()
}

// Stop is function that kills job
func (job *WatcherJob) Stop() error {
	if job.cmd != nil && job.cmd.Process != nil {
		return job.cmd.Process.Kill()
	}
	return nil
}
