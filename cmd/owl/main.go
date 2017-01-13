package main

import (
	"os"
	"path/filepath"
	"log"
	"github.com/fsnotify/fsnotify"
	"errors"
	"github.com/urfave/cli"
	"github.com/flowup/owl"
	"strconv"
	"os/exec"
	"fmt"
)

var (
	errFlagRunIsPresent = errors.New("flag --run or -r is required ")
)

type WatcherJob struct {
	command string
	cmd     *exec.Cmd
}

func NewWatcherJob(command string) *WatcherJob {
	return &WatcherJob{
		command: command,
		cmd:     nil,
	}
}

func (this*WatcherJob) Start() owl.JobResult {
	this.cmd = exec.Command("bash", "-c", this.command)
	out, err := this.cmd.CombinedOutput()
	return owl.JobResult{
		Output: string(out),
		Error:  err,
	}
}

func (this *WatcherJob) Stop() error {
	if this.cmd != nil && this.cmd.Process != nil {
		return this.cmd.Process.Kill()
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "owl"
	app.Usage = "owl watching all files in directory and when are changed, run the command"

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "print only the version",
	}

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name: "ignore, i",
			Usage:"All directories with name `IGNORE` are ignored",
		},
		cli.StringFlag{
			Name: "run, r",
			Usage:"If is any file changed, run `RUN`",
		},
		cli.BoolFlag{
			Name: "verbose, v",
			Usage:"verbose mode",
		},
		cli.StringFlag{
			Name: "time, t",
			Usage:"Waiting time for executing in miliseconds",
		},

	}

	app.Action = func(c *cli.Context) error {
		// default time is 5s
		amount := int64(5000)
		err := errors.New("")

		if c.String("t") != "" {
			amount, err = strconv.ParseInt(c.String("t"), 10, 64)
			if err != nil {
				panic(err)
			}
		}

		if c.String("run") == "" {
			return errFlagRunIsPresent
		}

		// set new watcher
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			panic(err)
		}

		// get path to this dir
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		// append path to global paths
		dirList := []string{}

		// this files are ignore
		ignoreList := make(map[string]bool)
		for _, dir := range (c.StringSlice("ignore")) {
			ignoreList[dir] = true
		}

		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// check if file is not in ignorelist
			if ignoreList[info.Name()] {
				return filepath.SkipDir
			}

			// append dir to list
			if info.IsDir() {
				dirList = append(dirList, path)
			}
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		// add all paths for watching
		for _, path := range dirList {
			watcher.Add(path)
		}

		jobs := make(chan owl.Job, 10)

		go func() {
			for {
				select {
				case ev := <-watcher.Events:

					// Write is running only once
					if ev.Op == fsnotify.Chmod {

						// execute of function with arguments
						if c.Bool("verbose") {
							log.Println(ev.Name)
						}

						// add fakeJob to jobs
						jobs <- &WatcherJob{command:c.String("run")}
					}
				case err := <-watcher.Errors:
					log.Fatal(err)
				}
			}
		}()

		debounced := owl.Debounce(jobs, amount)
		results := owl.Scheduler(debounced)

		for {
			out := <-results
			fmt.Print(out.Output)
		}

		return nil
	}
	app.Run(os.Args)
}
