package main

import (
	"os"
	"github.com/urfave/cli"
	"github.com/flowup/owl"
	"github.com/spf13/viper"
	"regexp"
	"strconv"
	"errors"
	"go.uber.org/zap"
)

var (
	errRunFlagMissing = errors.New("flag --run or -r is required ")
	defaultIgnored    = [...] string{"vendor", "node_modules", "bower_components", ".glide", ".git"}
)

func main() {

	o := owl.NewOwl()

	o.App.Action = func(c *cli.Context) error {

		o.Config.SetConfigType("yaml")
		o.Config.SetConfigName("owl")
		o.Config.AddConfigPath(".")

		err := o.Config.ReadInConfig()
		o.Config.SetDefault("debounce", 500)
		o.Config.SetDefault("verbose", false)
		o.Config.SetDefault("ignore", make([]string, 0))

		// If no config is present in current folder, read options from args
		if err != nil {
			o.Config.Set("run", c.String("run"))
			if c.Bool("v") {
				o.Config.Set("verbose", true)
			}
			if c.String("d") != "" {
				debounce, err := strconv.ParseInt(c.String("d"), 10, 64)
				if err != nil {
					panic(err)
				}
				o.Config.Set("debounce", debounce)
			}
			o.Config.Set("ignore", c.StringSlice("ignore"))
			o.Config.Set("filter", c.StringSlice("filter"))
		}

		if o.Config.GetString("run") == "" {
			return errRunFlagMissing
		}

		if o.Config.GetBool("verbose") {
			o.Logger = owl.GetLogger(zap.InfoLevel)
		} else {
			o.Logger = owl.GetLogger(zap.WarnLevel)
		}

		// job dir are ignored by default
		o.IgnoreList = make(map[string]bool)

		for _, dir := range defaultIgnored {
			o.IgnoreList[dir] = true
		}

		for _, dir := range viper.GetStringSlice("ignore") {
			o.IgnoreList[dir] = true
		}

		if o.Config.GetString("filter") != "" {
			o.FilterRegexp = regexp.MustCompile(o.Config.GetString("filter"))
		}

		err = o.SetupWatcher()
		if err != nil {
			return err
		}

		err = o.Watch()

		return err
	}

	o.App.Run(os.Args)
}
