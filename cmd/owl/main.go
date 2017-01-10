package main

import (
	"os"
	"path/filepath"
	"log"
	"os/exec"
	"github.com/fsnotify/fsnotify"
	"errors"
	"github.com/urfave/cli"
	"fmt"
)

var (
	errFlagRunIsPresent = errors.New("flag --run or -r is required ")
)

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
		cli.StringSliceFlag{
			Name: "run, r",
			Usage:"If is any file changed, run `RUN`",
		},
		cli.BoolFlag{
			Name: "verbose, v",
			Usage:"verbose mode",
		},

	}

	app.Action = func(c *cli.Context) error {

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

		for {
			select {
			case ev := <-watcher.Events:

				// Write is running only once
				if ev.Op == fsnotify.Chmod {
					// execute of function with arguments
					if c.Bool("verbose") {
						log.Println(ev.Name)
					}

					for _, r := range c.StringSlice("run") {
						output, err := exec.Command("bash", "-c", r).CombinedOutput()
						if err != nil {
							os.Stderr.WriteString(err.Error())
						}
						fmt.Print(string(output))
					}

				}


			case err := <-watcher.Errors:
				return err
			}
		}

		return nil
	}

	app.Run(os.Args)

}
