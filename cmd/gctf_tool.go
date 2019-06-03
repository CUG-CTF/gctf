package main

import (
	"fmt"
	"log"
	"github.com/urfave/cli"
	"os"
)

//todo check,backup,restore,upload
func main() {
	app:=cli.NewApp()
	app.Name = "gctf-tool"
	app.Usage="to check backup restore upload problem"
	app.Action= func(context *cli.Context) {
		fmt.Println("666")
	}
	err:=app.Run(os.Args)
	if err!=nil	{
		log.Fatal(err)
	}

}
