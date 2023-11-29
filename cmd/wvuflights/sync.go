package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	flights "github.com/austindizzy/wvu-flights/internal/wvuflights"
	"github.com/urfave/cli/v3"
	"gorm.io/gorm/clause"
)

var (
	syncCmd = &cli.Command{
		Name:  "sync",
		Usage: "Sync web content with the database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "web",
				Usage:    "the web content directory",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "the type of content to sync (trips | people | invoices)",
			},
		},
		Action: sync,
	}
	al *flights.AirportLookup
)

type tripData struct {
	ID            string
	StartDate     string
	EndDate       string
	Route         string
	Aircraft      string
	Destinations  []string
	Distance      float64
	RegNo         string
	NumPax        int
	Passengers    []flights.TripPassenger
	Crew          []string
	FlightHours   float64
	TotalCost     float64
	Cost          map[string]float64
	Invoice       string
	Justification string
	Notes         string
	SignedBy      string
	ApprovedBy    string
}

type peopleData struct {
	Name              string
	PersonType        string
	LastDept          string
	LastVPDiv         string
	Trips             []tripData
	TotalTrips        int
	TotalTripDistance float64
	TotalTripCost     float64
	TotalSoloTrips    int
	TotalSoloDistance float64
	TotalSoloCost     float64
}

func sync(c *cli.Context) error {
	var (
		trips     []flights.Trip
		people    []flights.Person
		invoices  []flights.Invoice
		err       error
		syncTypes = strings.Split(c.String("type"), ",")
		n         = 0
		start     = time.Now()
	)

	al, err = flights.NewAirportLookup(os.Getenv("MAPBOX_ACCESS_TOKEN"))
	if err != nil {
		return err
	}

	if slices.Contains(syncTypes, "trips") || len(c.String("type")) == 0 {
		db.Preload(clause.Associations).Find(&trips)
		fmt.Printf("%s: Syncing %d trips...", time.Now(), len(trips))
		for _, trip := range trips {
			err := saveTripToFile(trip, c.String("web"))
			if err != nil {
				return err
			}
			n++
		}
		fmt.Printf("synced %d trips\n", n)
	}

	if slices.Contains(syncTypes, "people") || len(c.String("type")) == 0 {
		db.Find(&people)
		fmt.Printf("%s: Syncing %d people...", time.Now(), len(people))
		n = 0

		for _, person := range people {
			// skip LJA crewpeople for now
			if person.Type == flights.PersonTypeLJA {
				continue
			}

			err := savePersonToFile(person, c.String("web"))
			if err != nil {
				return err
			}
			n++
		}

		fmt.Printf("synced %d people\n", n)
	}

	if slices.Contains(syncTypes, "invoices") || len(c.String("type")) == 0 {
		db.Find(&invoices)
		fmt.Printf("%s: Syncing %d invoices...", time.Now(), len(invoices))
		n = 0

		for _, invoice := range invoices {
			err := saveInvoiceToFile(invoice, c.String("web"))
			if err != nil {
				return err
			}
			n++
		}

		fmt.Printf("synced %d invoices\n", n)
	}

	fmt.Printf("Finished sync in %s\n", time.Since(start))

	return nil
}
