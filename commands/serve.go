package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gernest/hiro/config"
	"github.com/gernest/hiro/query"
	"github.com/gernest/hiro/server"
	"github.com/gernest/hiro/templates"
	"github.com/urfave/cli"
)

// ServeCMD command for starting webserver.
func ServeCMD(version, commit, buildDate string) cli.Command {
	return cli.Command{
		Name:  "serve",
		Usage: "starts the hiro server",
		Flags: []cli.Flag{
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
			cli.IntFlag{
				Name:   "port",
				Usage:  "port to bind the server",
				EnvVar: "HIRO_PORT",
				Value:  8080,
			},
		},
		Action: Serve(version, commit, buildDate),
	}
}

// Serve starts the server.
func Serve(version, commit, buildDate string) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		start := time.Now()
		cfg, err := config.FromCtx(ctx)
		if err != nil {
			return err
		}
		db, err := query.New(cfg.DBDriver, cfg.DBConn)
		if err != nil {
			return err
		}
		//perform migrations
		if err = db.Up(context.Background()); err != nil {
			return err
		}
		err = templates.Write(os.Stdout, templates.HomeBanner, map[string]interface{}{
			"Duration":    time.Now().Sub(start),
			"Version":     version,
			"Author":      "Geofrey Ernest",
			"Email":       "geofreyernest@live.com",
			"Website":     "https://bq.co.tz",
			"DocsWebsite": "https://docs.bq.co.tz",
			"Repository":  "https://github.com/gernest/bq",
			"Twitter":     "@gernesti",
			"Host":        fmt.Sprintf(":%d", cfg.Port),
			"Commit":      commit,
			"BuildDate":   buildDate,
		})
		if err != nil {
			return err
		}
		return server.ServeAPI(context.Background(), db, cfg)
	}
}
