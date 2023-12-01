package main

import (
	"fmt"
	"os"
	"strings"

	flights "github.com/austindizzy/wvu-flights/internal/wvuflights"
	"github.com/urfave/cli/v3"
)

var (
	utilsCmd = &cli.Command{
		Name:  "utils",
		Usage: "Utilities for doing various tasks",
		Commands: []*cli.Command{
			{
				Name:  "generate-img",
				Usage: "Generate an image of a route",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "route",
						Usage:    "the route to generate an image of",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Usage:    "the output file",
						Required: true,
					},
				},
				Action: generateRouteImg,
			},
		},
	}
)

func generateRouteImg(c *cli.Context) error {
	if os.Getenv("MAPBOX_ACCESS_TOKEN") == "" {
		return fmt.Errorf("MAPBOX_ACCESS_TOKEN environment variable not set")
	}

	al, err := flights.NewAirportLookup(os.Getenv("MAPBOX_ACCESS_TOKEN"))
	if err != nil {
		return err
	}

	route := c.String("route")
	output := c.String("output")

	data, err := al.RouteToImage(route)
	if err != nil {
		return err
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote image to %s\n", output)

	return nil
}
