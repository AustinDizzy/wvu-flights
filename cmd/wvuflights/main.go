package main

import (
	"context"
	"log"
	"os"

	flights "github.com/austindizzy/wvu-flights/internal/wvuflights"

	"github.com/urfave/cli/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func main() {
	app := &cli.Command{
		Name:  "wvuflights",
		Usage: "A command line tool for aggregating and visualizing private charter flight data from West Virginia University",
		Commands: []*cli.Command{
			{
				Name:   "server",
				Usage:  "Start the webserver",
				Action: server,
			},
			{
				Name:  "trip",
				Usage: "Manage trips",
				Commands: []*cli.Command{
					tripAddCmd,
				},
			},
			{
				Name:  "person",
				Usage: "Manage people",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "the person's name",
					},
				},
				Commands: []*cli.Command{
					personAddCmd,
					personLookupCmd,
				},
			},
			syncCmd,
			utilsCmd,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "db",
				Usage:    "the SQLite database file (.db | .sqlite | .sqlite3)",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "enable verbose logging",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "dryrun",
				Usage: "enable dry run mode",
				Value: false,
			},
		},
		Before: func(c *cli.Context) error {
			var err error
			db, err = gorm.Open(sqlite.Open(c.String("db")), &gorm.Config{
				DryRun: c.Bool("dryrun"),
			})
			if err != nil {
				return err
			}

			db.AutoMigrate(&flights.Trip{}, &flights.TripPassenger{}, &flights.Person{}, &flights.Invoice{})

			return nil
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
