package main

import (
	"errors"
	"log"
	"os"
	"sort"
	"time"

	"github.com/hiendv/geojson/internal/hxxp"
	"github.com/hiendv/geojson/internal/osm"
	"github.com/hiendv/geojson/internal/shared"
	"github.com/hiendv/geojson/pkg/util"
	"github.com/urfave/cli/v2"
)

// NewSubAreaCommand constructs sub-command SubArea.
func NewSubAreaCommand() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		relation := c.Args().First()
		if relation == "" {
			return errors.New("invalid OpenStreetMap relation ID")
		}

		logger, ok := c.App.Metadata["logger"].(shared.Logger)
		if !ok || logger == nil {
			return errors.New("invalid logger")
		}

		ctx, err := osm.NewContext(
			c.Context,
			logger,
			c.Bool("raw"),
			c.Bool("separated"),
			c.String("out"),
			c.Bool("rewind"),
		)
		if err != nil {
			return err
		}

		err = osm.SubAreas(ctx, relation)
		if err != nil {
			logger.Error(err)
		}
		return nil
	}
}

// NewServeCommand constructs sub-command Serve.
func NewServeCommand() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		logger, ok := c.App.Metadata["logger"].(shared.Logger)
		if !ok || logger == nil {
			return errors.New("invalid logger")
		}

		ttl, err := util.ParseDuration(c.String("rate-ttl"))
		if err != nil {
			return errors.New("invalid duration")
		}

		ctx, err := hxxp.NewContext(
			c.Context,
			logger,
			c.String("address"),
			c.String("origin"),
			c.Float64("rate"),
			c.Int("rate-burst"),
			ttl,
			c.String("out"),
			c.String("prefix"),
		)
		if err != nil {
			return err
		}

		handler, err := hxxp.New(ctx)
		if err != nil {
			return errors.New("could not create the request handler")
		}

		return hxxp.Listen(handler)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "GeoJSON"
	app.Usage = "Utilities for OpenStreetMap GeoJSON"
	app.Authors = []*cli.Author{
		{
			Name:  "Hien Dao",
			Email: "hien.dv.neo@gmail.com",
		},
	}
	app.Copyright = "Copyright © 2020 Hien Dao. All Rights Reserved."
	app.Version = "0.1.1"
	app.Compiled = time.Now()
	app.Commands = []*cli.Command{
		{
			Name:   "subarea",
			Usage:  "list all sub-areas of an OpenStreetMap object",
			Action: NewSubAreaCommand(),
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "raw",
					Aliases: []string{"r"},
					Usage:   "leave tags in unfornalized form (UNF)",
				},
				&cli.BoolFlag{
					Name:    "separated",
					Aliases: []string{"s"},
					Usage:   "leave sub-areas unmerged",
				},
				&cli.BoolFlag{
					Name:  "rewind",
					Usage: "rewind the output - counter to RFC 7946",
				},
			},
		},
		{
			Name:   "serve",
			Usage:  "serve the web server",
			Action: NewServeCommand(),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "address",
					Aliases: []string{"addr"},
					Value:   "127.0.0.1:8181",
					Usage:   "set the serving address",
				},
				&cli.StringFlag{
					Name:  "origin",
					Value: "*",
					Usage: "set the CORS origin",
				},
				&cli.IntFlag{
					Name:  "rate",
					Value: 10,
					Usage: "set request-per-second for rate-limiting",
				},
				&cli.IntFlag{
					Name:  "rate-burst",
					Value: 5,
					Usage: "set burst size (concurrent requests) for rate-limiting",
				},
				&cli.StringFlag{
					Name:  "rate-ttl",
					Value: "2m",
					Usage: "set the rate limit TTL for inactive sessions",
				},
				&cli.StringFlag{
					Name:  "prefix",
					Value: "/static",
					Usage: "set static fs handler base path",
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "enable verbose logging with DEBUG level",
		},
		&cli.StringFlag{
			Name:    "out",
			Aliases: []string{"o"},
			Usage:   "specify the directory of outputs",
			Value:   "./geo",
		},
	}
	app.Before = func(c *cli.Context) error {
		logger, err := shared.NewLoggerZap(c.Bool("verbose"))
		if err != nil {
			return err
		}

		app.Metadata["logger"] = logger
		return nil
	}
	app.After = func(c *cli.Context) error {
		logger, ok := app.Metadata["logger"].(shared.Logger)
		if !ok || logger == nil {
			return nil
		}

		// we ignore the error because of unknown syscall errors or os path errors when syncing /dev/stderr
		// nolint:errcheck
		logger.Sync()
		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err == nil {
		return
	}

	log.Println(err)
	os.Exit(1)
}
