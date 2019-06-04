package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

//todo check,backup,restore,upload
func main() {
	app := cli.NewApp()
	app.Name = "gctf-tool"
	app.Version = "0.01"
	app.Usage = "to check backup restore upload problem"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf,c",
			Value: "conf.json",
			Usage: "config file path",
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("666")
		fmt.Println(c.String("c"))
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
