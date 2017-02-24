package owl

import (
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"errors"
	"strconv"
	"go.uber.org/zap"
	"regexp"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"fmt"
)

var (
	errRunFlagMissing = errors.New("flag --run or -r is required ")
	defaultIgnored    = [...] string{"vendor", "node_modules", "bower_components", ".glide", ".git"}
)

// Owl is tool for advanced watching of changes in directories and files
type Owl struct {
	App    *cli.App
	Logger *zap.Logger

	config       *viper.Viper
	watcher      *fsnotify.Watcher
	filterRegexp *regexp.Regexp
	ignoreList   map[string]bool
}

// NewOwl is function that creates new instance of the Owl
func NewOwl() *Owl {
	return &Owl{
		config:     viper.New(),
		App:        GetApp(),
		ignoreList: make(map[string]bool),
	}
}

// ReadConfigAndInit reads config from configuration file and if configuration file is not present then configuration
// is read from program arguments.
func (o *Owl) ReadConfigAndInit(c *cli.Context) error {

	o.config.SetConfigType("yaml")
	o.config.SetConfigName("owl")
	o.config.AddConfigPath(".")

	err := o.config.ReadInConfig()
	o.config.SetDefault("debounce", 500)
	o.config.SetDefault("verbose", false)
	o.config.SetDefault("ignore", make([]string, 0))

	// If no config is present in current folder, read options from args
	if err != nil {
		o.config.Set("run", c.String("run"))
		if c.Bool("v") {
			o.config.Set("verbose", true)
		}
		if c.String("d") != "" {
			debounce, err := strconv.ParseInt(c.String("d"), 10, 64)
			if err != nil {
				panic(err)
			}
			o.config.Set("debounce", debounce)
		}
		o.config.Set("ignore", c.StringSlice("ignore"))
		o.config.Set("filter", c.StringSlice("filter"))
	}

	if o.config.GetString("run") == "" {
		return errRunFlagMissing
	}

	if o.config.GetBool("verbose") {
		o.Logger = GetLogger(zap.InfoLevel)
	} else {
		o.Logger = GetLogger(zap.WarnLevel)
	}

	// job dir are ignored by default
	o.ignoreList = make(map[string]bool)

	for _, dir := range defaultIgnored {
		o.ignoreList[dir] = true
	}

	for _, dir := range viper.GetStringSlice("ignore") {
		o.ignoreList[dir] = true
	}

	if o.config.GetString("filter") != "" {
		o.filterRegexp = regexp.MustCompile(o.config.GetString("filter"))
	}

	// set new watcher
	o.watcher, err = fsnotify.NewWatcher()

	return err
}

// SetupWatcher is function for setting watcher's directories
func (o *Owl) SetupWatcher() error {

	// get path to job dir
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// check if dir is not in ignorelist
		if o.ignoreList[info.Name()] && info.IsDir() {
			return filepath.SkipDir
		}

		if info.IsDir() {
			o.watcher.Add(path)
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
		Command: o.config.GetString("run"),
		OutPipe: os.Stdout}

	// start the command at the start of the owl
	jobs <- watcherJob

	go func() {
		for {
			select {
			case ev := <-o.watcher.Events:

				// check if is set filter
				if o.filterRegexp != nil {
					if !o.filterRegexp.MatchString(ev.Name) {
						break
					}
				}

				// Write is running only once
				if ev.Op == fsnotify.Create {
					o.watcher.Add(ev.Name)
				}

				// Log event
				o.Logger.Info(ev.Name)

				// Add Job to jobs
				jobs <- watcherJob

			case err := <-o.watcher.Errors:
				o.Logger.Fatal(err.Error())
			}
		}
	}()

	dJobs := Debounce(jobs, o.config.GetInt64("debounce"))
	results := Scheduler(dJobs)

	for {
		err := <-results
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}
