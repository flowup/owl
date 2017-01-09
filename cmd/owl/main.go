package main

import (
	"fmt"
	"path/filepath"
	"os"
	"log"
	"regexp"
	"os/exec"
	"github.com/fsnotify/fsnotify"
)

func main() {

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

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		// this files are ignore
		ignoreList := []string{"vendor", ".glide", ".git"}

		// check if file is not in ignorelist
		for _, i := range ignoreList {
			if info.Name() == i {
				return filepath.SkipDir
			}
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
	re := regexp.MustCompile("^([A-Za-z0-9]+)\\s?(.*)")
	groups := re.FindStringSubmatch(os.Args[1:][0])

	// name of function
	name := groups[1]
	// his arguments
	arq := groups[2]

	for {
		select {
		case ev := <-watcher.Events:
			// Write is running only once
			if ev.Op == fsnotify.Write {
				// execute of function with arguments
				output, err := exec.Command(name, arq).CombinedOutput()
				if err != nil {
					os.Stderr.WriteString(err.Error())
				}
				fmt.Print(string(output))

			}


		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}

}
