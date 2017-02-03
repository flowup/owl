package main

import (
	"os"
	"path/filepath"
	"github.com/fsnotify/fsnotify"
	"errors"
	"github.com/urfave/cli"
	"github.com/flowup/owl"
	"github.com/uber-go/zap"
	"strconv"
	"os/exec"
	"github.com/spf13/viper"
	"fmt"
	"bufio"
	"io"
	"syscall"
	"regexp"
)

var (
	errFlagRunIsPresent = errors.New("flag --run or -r is required ")
)

type WatcherJob struct {
	command string
	cmd     *exec.Cmd
	outpipe io.Writer
}

func NewWatcherJob(command string) *WatcherJob {
	return &WatcherJob{
		command: command,
		cmd:     nil,
		outpipe: os.Stdout,
	}
}

func (job *WatcherJob) Start() error {
	job.cmd = exec.Command("bash", "-c", job.command)
	job.cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

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
				fmt.Fprintln(job.outpipe, outscanner.Text())
			} else if errscanner.Scan() {
				fmt.Fprintln(job.outpipe, errscanner.Text())
			} else {
				break
			}
		}
	}()

	return job.cmd.Run()
}

func (job *WatcherJob) Stop() error {
	if job.cmd != nil && job.cmd.Process != nil {
		return syscall.Kill(-job.cmd.Process.Pid, syscall.SIGKILL)
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
			Name: "debounce, d",
			Usage:"Waiting time for executing in miliseconds",
		},
		cli.StringFlag{
			Name: "filter, f",
			Usage:"Files are filtered by expression",
		},
	}

	app.Action = func(c *cli.Context) error {

		viper.SetConfigType("yaml")
		viper.SetConfigName("owl")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		viper.SetDefault("debounce", 500)
		viper.SetDefault("verbose", false)
		viper.SetDefault("ignore", make([]string, 0))

		// If no config is present in current folder, read options from args
		if err != nil {
			viper.Set("run", c.String("run"))
			if c.Bool("v") {
				viper.Set("verbose", true)
			}
			if c.String("d") != "" {
				debounce, err := strconv.ParseInt(c.String("d"), 10, 64)
				if err != nil {
					panic(err)
				}
				viper.Set("debounce", debounce)
			}
			viper.Set("ignore", c.StringSlice("ignore"))
			viper.Set("filter", c.StringSlice("filter"))
		}

		if viper.GetString("run") == "" {
			return errFlagRunIsPresent
		}

		var rgxp *regexp.Regexp
		if viper.GetString("filter") != "" {
			rgxp = regexp.MustCompile(viper.GetString("filter"))
		}

		err = errors.New("")

		loglevel := zap.WarnLevel
		if viper.GetBool("verbose") {
			loglevel = zap.InfoLevel
		}

		logger := zap.New(zap.NewTextEncoder(), loglevel)

		// set new watcher
		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			panic(err)
		}

		// get path to job dir
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dirList := []string{}

		// job dir are ignored by default
		ignoreList := make(map[string]bool)
		ignoreList["vendor"] = true
		ignoreList["node_modules"] = true
		ignoreList["bower_components"] = true
		ignoreList[".glide"] = true
		ignoreList[".git"] = true

		for _, dir := range (viper.GetStringSlice("ignore")) {
			ignoreList[dir] = true
		}

		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// check if dir is not in ignorelist
			if ignoreList[info.Name()] && info.IsDir() {
				return filepath.SkipDir
			}

			if info.IsDir() {
				dirList = append(dirList, path)
			}
			return nil
		})

		if err != nil {
			logger.Fatal(err.Error())
		}

		// add all paths for watching
		for _, path := range dirList {
			watcher.Add(path)
		}

		jobs := make(chan owl.Job, 10)

		// init job
		watcherJob := &WatcherJob{
			command: viper.GetString("run"),
			outpipe: os.Stdout}

		// to start the command at the start of the owl
		jobs <- watcherJob

		go func() {
			for {
				select {
				case ev := <-watcher.Events:

					// Write is running only once
					if ev.Op == fsnotify.Chmod {

						// check if is set filter
						if viper.GetString("filter") != "" {
							if !rgxp.MatchString(ev.Name) {
								break
							}
						}

						// log event
						logger.Info(ev.Name)

						// add Job to jobs
						jobs <- watcherJob
					}
				case err := <-watcher.Errors:
					logger.Fatal(err.Error())
				}
			}
		}()

		debounced := owl.Debounce(jobs, viper.GetInt64("debounce"))
		results := owl.Scheduler(debounced)

		for {
			err := <-results
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		return nil
	}
	app.Run(os.Args)
}
