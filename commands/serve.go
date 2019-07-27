package commands

import (
	"context"
	"os"
	"time"

	"github.com/gernest/hiro/config"
	"github.com/gernest/hiro/query"
	"github.com/gernest/hiro/server"
	"github.com/gernest/hiro/templates"
	"github.com/gernest/hiro/util"
	"github.com/urfave/cli"
)

// ServeCMD command for starting webserver.
func ServeCMD(version, commit, buildDate string) cli.Command {
	fl := []cli.Flag{
		util.ConnFlag(),
		util.DriverFlag(),
		util.SecretFlag(),
		util.HostFlag(),
		util.ImageHostFlag(),
	}
	fl = append(fl, util.MinioFlags()...)
	fl = append(fl, util.NSQFlags()...)
	return cli.Command{
		Name:   "serve",
		Usage:  "starts the bq server",
		Flags:  fl,
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
			"Host":        cfg.Host,
			"Commit":      commit,
			"BuildDate":   buildDate,
		})
		if err != nil {
			return err
		}
		return server.ServeAPI(context.Background(), db, cfg)
	}
}
