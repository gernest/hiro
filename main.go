package main

import (
	"log"
	"os"

	"github.com/gernest/hiro/commands"
	"github.com/urfave/cli"
)


var (
	commit  = "<>"
	date    = "<>"
	version = "dev"
)

func main() {
	a := cli.NewApp()
	a.Name = "hiro"
	a.Usage = "modern qrcode service"
	a.EnableBashCompletion = true
	a.Description = "zero hype, high performance qrcode service"
	a.Version = version
	a.Commands = []cli.Command{
		commands.ServeCMD(version, commit, date),
	}
	if err := a.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
