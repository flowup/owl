package main

import (
	"os"
	"github.com/urfave/cli"
	"github.com/flowup/owl"
)


func main() {

	o := owl.NewOwl()

	o.App.Action = func(c *cli.Context) error {

		err := o.ReadConfigAndInit(c)
		if err != nil {
			return err
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
