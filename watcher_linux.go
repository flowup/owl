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
func (job *WatcherJob) Start() error {
	job.Cmd = exec.Command("bash", "-c", job.Command)
	job.Cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	stderr, err := job.Cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := job.Cmd.StdoutPipe()
	if err != nil {
		return err
	}

	errscanner := bufio.NewScanner(stderr)
	outscanner := bufio.NewScanner(stdout)

	go func() {
		for {
			if outscanner.Scan() {
				fmt.Fprintln(job.OutPipe, outscanner.Text())
			} else if errscanner.Scan() {
				fmt.Fprintln(job.OutPipe, errscanner.Text())
			} else {
				break
			}
		}
	}()

	return job.Cmd.Run()
}

// Stop is function that kills job
func (job *WatcherJob) Stop() error {
	if job.Cmd != nil && job.Cmd.Process != nil {
		return syscall.Kill(-job.Cmd.Process.Pid, syscall.SIGKILL)
	}
	return nil
}