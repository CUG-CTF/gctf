package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

//todo check,backup,restore,sync
//problem's dir -> db -> docker images
var ConfigFIle ,ProblemsDir,BackupDir string
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
			Destination:&ConfigFIle,
		},
		cli.StringFlag{
			Name:"problems,p",
			Value:"problems",
			Usage:"problems dir",
			Destination:&ProblemsDir,
		},
		cli.StringFlag{
			Name:"backup-dir,b",
			Value:"gctf_backups",
			Usage:"problems backup dir",
			Destination:&BackupDir,

		},
	}
	app.Commands=cli.Commands{
		cli.Command{
			Name:"check",
			Usage:"check if problem's dir equl to db and docker images,warning! this will del problem which not in problem's dir!",
			Action: func(c *cli.Context) error{
				fmt.Println("this is check function")
				return nil
			},
		},
		//Todo: back up users,team...etc
		cli.Command{
			Name:"backup",
			Usage:"backup db,not problem",
		},
		cli.Command{
			Name:"sync",
			Usage:"sync problem's dir to db and docker images",
		},
	}
	//app.Action = func(c *cli.Context) error {
	//	fmt.Println()
	//	return nil
	//}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
