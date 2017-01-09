package main

import (
	"github.com/flowup/owl/.glide/cache/src/https-github.com-urfave-cli"
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"log"
	"regexp"
	"os/exec"
	"github.com/fsnotify/fsnotify"
)

func main() {
	app := cli.NewApp()
	app.Name = "owl"
	app.Usage = "owl watching all files in directory and when are changed, run the command"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "ignore, i",
			Usage:"All directories with name `IGNORE` are ignored",
		},
		cli.StringFlag{
			Name: "run, r",
			Usage:"If is any file changed, run `RUN`",
		},
	}

	app.Action = func(c *cli.Context) error {
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
		for _, dir := range (strings.Split(c.String("ignore"), ";")) {
			ignoreList[dir] = true
		}

		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Print(err)
				return nil
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

		// seperate argument into two parts
		fmt.Println(c.String("run"))

		re := regexp.MustCompile("^([A-Za-z0-9]+)\\s?(.*)")

		var name []string
		var arq []string

		for _, i := range strings.Split(c.String("run"), ";") {
			groups := re.FindStringSubmatch(i)

			// name of function
			name = append(name, groups[1])
			// his arguments
			arq = append(arq, groups[2])
		}

		for {
			select {
			case ev := <-watcher.Events:
				// Write is running only once
				if ev.Op == fsnotify.Write {
					// execute of function with arguments
					for x, _ := range name {
						output, err := exec.Command(name[x], arq[x]).CombinedOutput()
						if err != nil {
							os.Stderr.WriteString(err.Error())
						}
						fmt.Print(string(output))
					}

				}


			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}

		return nil
	}

	app.Run(os.Args)

}
