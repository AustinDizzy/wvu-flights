package main

import (
	"fmt"
	"os"
	"regexp"
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

func toSlug(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9']+`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`^-+|-+$`).ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "'", "")
	return s
}

func cleanRouteStr(route string) string {
	route = strings.ToUpper(route)
	route = regexp.MustCompile(`-RON\([0-9]+\)`).ReplaceAllString(route, "")

	segments := strings.Split(route, "-")
	res := []string{}

	for i := 0; i < len(segments); i++ {
		res = append(res, segments[i])

		for i+1 < len(segments) && segments[i] == segments[i+1] {
			i++
		}
	}

	return strings.Join(res, "-")
}
