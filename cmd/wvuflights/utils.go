package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	flights "github.com/austindizzy/wvu-flights/internal/wvuflights"
	"github.com/urfave/cli/v3"
	"gorm.io/gorm/clause"
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
			{
				Name:  "count-trips",
				Usage: "Count the number of trips including only specific destinations (defaults to destinations in West Virginia)",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:  "dest",
						Usage: "a destination from route",
						Value: strings.Split("LBE MGW CRW CKB LWB BLF EKN HLG PGC MRB PKB HTS BKW", " "),
					},
				},
				Action: countTrips,
			},
		},
	}
)

func countTrips(c *cli.Context) error {
	dests := c.StringSlice("dest")
	if len(dests) == 0 {
		return fmt.Errorf("must specify at least one destination")
	}

	if len(dests) == 1 && (strings.Contains(dests[0], ",") || strings.Contains(dests[0], " ")) {
		dests = strings.Split(dests[0], ",")
		for i, dest := range dests {
			dests[i] = strings.TrimSpace(dest)
		}

		if len(dests) == 1 {
			dests = strings.Split(dests[0], " ")
			for i, dest := range dests {
				dests[i] = strings.TrimSpace(dest)
			}
		}
	}

	var trips []flights.Trip
	db.Preload(clause.Associations).Find(&trips)

	count := 0
	totalCost := 0.0
	totalDistance := 0.0

	al, _ := flights.NewAirportLookup()
	for _, trip := range trips {

		locations, err := al.ToNames(trip.Route)
		if err != nil {
			return err
		}

		match := true
		for dest, _ := range locations {
			if !slices.Contains(dests, dest) {
				match = false
				break
			}
		}

		if match {
			count++
			totalCost += trip.GetTotalCost()
			dist, err := al.RouteToDistance(trip.Route)
			if err != nil {
				return err
			}
			totalDistance += dist
		}
	}

	fmt.Printf("Trips only including one of %v (%d):\n\n", dests, len(dests))
	fmt.Printf("Total Trips: %d\n", count)
	fmt.Printf("Total Cost: $%.2f\n", totalCost)
	fmt.Printf("Total Distance: %.2f nmi\n", totalDistance)

	return nil
}

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
