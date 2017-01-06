package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"log"
	"io/ioutil"
	"os/exec"
	"regexp"
)

// all paths for watching
var Paths []string

func getPaths(folder string) {
	files, _ := ioutil.ReadDir(folder)
	folder = folder + "/"

	for _, f := range files {
		// is not dir
		if !f.IsDir() {
			continue
		}

		//ignore this (.git, vendor)
		if string([]rune(f.Name())[0]) == "." || f.Name() == "vendor" {
			continue
		}

		// connect name of folder with his path
		path := folder + f.Name()
		Paths = append(Paths, path)
		getPaths(path)
	}
}

func main() {

	//set new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	//get path to this dir
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// append path to global paths
	Paths = append(Paths, path)

	// find all paths in this dir (recursively)
	getPaths(path)

	// add all paths for watching
	for _, path := range Paths {
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
