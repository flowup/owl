package owl

import (
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"regexp"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"fmt"
)

// Owl is tool for advanced watching of changes in directories and files
type Owl struct {
	App    *cli.App
	Logger *zap.Logger

	Config       *viper.Viper
	Watcher      *fsnotify.Watcher
	FilterRegexp *regexp.Regexp
	IgnoreList   map[string]bool
}

// NewOwl is function that creates new instance of the Owl
func NewOwl() *Owl {
	return &Owl{
		Config:     viper.New(),
		App:        GetApp(),
		IgnoreList: make(map[string]bool),
	}
}

// SetupWatcher is function for setting watcher's directories
func (o *Owl) SetupWatcher() error {

	// get path to job dir
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// set new watcher
	o.Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// check if dir is not in ignorelist
		if o.IgnoreList[info.Name()] && info.IsDir() {
			return filepath.SkipDir
		}

		if info.IsDir() {
			o.Watcher.Add(path)
		}

		return nil
	})

	return err
}

// Watch is function that watches changes in files and runs jobs
func (o *Owl) Watch() error {

	jobs := make(chan Job, 24)

	// init job
	watcherJob := &WatcherJob{
		Command: o.Config.GetString("run"),
		OutPipe: os.Stdout}

	// start the command at the start of the owl
	jobs <- watcherJob

	go func() {
		for {
			select {
			case ev := <-o.Watcher.Events:

				// check if is set filter
				if o.FilterRegexp != nil {
					if !o.FilterRegexp.MatchString(ev.Name) {
						break
					}
				}

				// Write is running only once
				if ev.Op == fsnotify.Create {
					o.Watcher.Add(ev.Name)
				}

				// Log event
				o.Logger.Info(ev.Name)

				// Add Job to jobs
				jobs <- watcherJob

			case err := <-o.Watcher.Errors:
				o.Logger.Fatal(err.Error())
			}
		}
	}()

	dJobs := Debounce(jobs, o.Config.GetInt64("debounce"))
	results := Scheduler(dJobs)

	for {
		err := <-results
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}
