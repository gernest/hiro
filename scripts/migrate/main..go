package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gernest/hiro/query"
	"github.com/urfave/cli"
)

func main() {
	a := cli.NewApp()
	a.Name = "hiro-migrate"
	a.Usage = "runs hiro migrations"
	a.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db-conn",
			Usage:  "connection string to postgres database",
			EnvVar: "HIRO_DB_CONN",
		},
		cli.StringFlag{
			Name:   "driver",
			Usage:  "database driver to use",
			EnvVar: "HIRO_DB_DRIVER",
			Value:  "postgres",
		},
		cli.StringFlag{
			Name:   "secret",
			Usage:  "hmac jwt secret",
			EnvVar: "HIRO_JWT_SECRET",
			Value:  "secret",
		},
	}
	a.Commands = []cli.Command{
		{
			Name:   "up",
			Usage:  "creates tables if they don't exist yet",
			Action: up,
		},
		{
			Name:   "down",
			Usage:  "clears the database by dropping the tables",
			Action: down,
		},
	}
	if err := a.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func up(ctx *cli.Context) error {
	driver := ctx.GlobalString("driver")
	conn := ctx.GlobalString("db-conn")
	db, err := query.New(driver, conn)
	if err != nil {
		return err
	}
	return db.Up(context.Background())
}

func down(ctx *cli.Context) error {
	driver := ctx.GlobalString("driver")
	conn := ctx.GlobalString("db-conn")
	db, err := query.New(driver, conn)
	if err != nil {
		return err
	}
	return db.Down(context.Background())
}
