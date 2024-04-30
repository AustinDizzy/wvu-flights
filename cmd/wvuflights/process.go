package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/austindizzy/wvu-flights/internal/wvuflights"
	flights "github.com/austindizzy/wvu-flights/internal/wvuflights"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm/clause"
)

func saveTripToFile(trip flights.Trip, dir string) error {
	dateRange := strings.Split(trip.Date, ";")
	endDate := ""
	if len(dateRange) > 1 {
		endDate = dateRange[1]
	}
	crewNames := []string{}
	for _, crew := range trip.Crew {
		crewNames = append(crewNames, crew.Name)
	}

	if trip.Distance <= 0 {
		dist, err := al.RouteToDistance(trip.Route)
		if err != nil {
			return err
		}

		trip.Distance = dist
		db.Save(&trip)
	}

	if trip.CarbonFootprint <= 0 {
		trip.CarbonFootprint = trip.GetCarbonFootprint()
		db.Save(&trip)
	}

	fb := float64(wvuflights.GetAircraftFuelBurn(trip.Aircraft))
	if fb != -1 {
		fb = fb * trip.FlightHours
	}

	var destinations []string
	locations, err := al.ToNames(trip.Route)
	if err != nil {
		return err
	}

	for _, l := range locations {
		destinations = append(destinations, l)
	}

	var (
		// route pieces
		rp = strings.Split(flights.CleanRouteStr(trip.Route), "-")
		// route str pieces
		rsp []string
		// route string (trip title)
		routeStr string
	)

	for i, c := range rp {
		if (i == 0 || i == len(rp)-1) && c == "LBE" {
			continue
		}

		rsp = append(rsp, strings.TrimSuffix(locations[c], ", WV"))
	}

	if len(rsp) > 1 && rsp[0] == rsp[len(rsp)-1] {
		routeStr = fmt.Sprintf("%s (round trip)", strings.Join(rsp[:len(rsp)-1], " to "))
	} else {
		routeStr = strings.Join(rsp, " to ")
	}

	invoiceName := regexp.MustCompile(`;?([0-9]{4}-[0-9]{2})-[0-9]{2}$`).FindAllStringSubmatch(trip.Date, 1)[0][1]
	var nInvoice int64

	db.Model(&flights.Invoice{}).Where("name = ? AND type = ?", invoiceName, "month").Count(&nInvoice)
	if nInvoice > 0 {
		invoiceName = fmt.Sprintf("%s.pdf", invoiceName)
	} else {
		invoiceName = ""
	}

	frontMatter := tripData{
		ID:           trip.ID,
		StartDate:    dateRange[0],
		EndDate:      endDate,
		Route:        trip.Route,
		RouteStr:     routeStr,
		Aircraft:     trip.Aircraft,
		Destinations: destinations,
		Distance:     trip.Distance,
		Carbon:       trip.CarbonFootprint,
		FuelBurn:     fb,
		RegNo:        trip.RegNo,
		NumPax:       trip.NumPax,
		Passengers:   trip.Passengers,
		Crew:         crewNames,
		FlightHours:  trip.FlightHours,
		TotalCost:    trip.GetTotalCost(),
		Cost: map[string]float64{
			"fuel":           trip.Fuel,
			"landing_fees":   trip.Landing,
			"crew_expense":   trip.CrewExpense,
			"domestic_tax":   trip.DomTax,
			"hourly_rate":    trip.HourlyRate,
			"billing_amount": trip.BillingAmount,
		},
		Justification: trip.Justification,
		Invoice:       invoiceName,
		Notes:         trip.Notes,
		SignedBy:      trip.SignedByName,
		ApprovedBy:    trip.ApprovedByName,
	}

	data, err := yaml.Marshal(frontMatter)
	if err != nil {
		return err
	}

	dirName := fmt.Sprintf("%s/content/trips/%s", dir, trip.ID)
	err = os.MkdirAll(dirName, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/index.md", dirName))
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = fmt.Fprintf(f, "---\n%s\n---\n", data)
	if err != nil {
		return err
	}

	r := trip.GetReservation()
	if r != nil {
		f, err := os.Create(fmt.Sprintf("%s/reservation.pdf", dirName))
		if err != nil {
			return err
		}
		_, err = f.Write(r)
		if err != nil {
			return err
		}
	}

	i := trip.GetItinerary()
	if i != nil {
		f, err := os.Create(fmt.Sprintf("%s/itinerary.pdf", dirName))
		if err != nil {
			return err
		}
		_, err = f.Write(i)
		if err != nil {
			return err
		}
	}

	cr := wvuflights.CleanRouteStr(trip.Route)
	routeImg := fmt.Sprintf("%s/static/img/routes/%s.png", dir, cr)
	if _, err := os.Stat(routeImg); err == nil {
		// file exists
		f, err := os.Open(routeImg)
		if err != nil {
			return err
		}
		defer f.Close()

		f2, err := os.Create(fmt.Sprintf("%s/route.png", dirName))
		if err != nil {
			return err
		}
		defer f2.Close()

		_, err = f2.ReadFrom(f)
		if err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		// file does not exist
		data, err := al.RouteToImage(cr)
		if err != nil {
			return err
		}

		f, err := os.Create(fmt.Sprintf("%s/route.png", dirName))
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write(data)
		if err != nil {
			return err
		}

		// save to routeImg
		f, err = os.Create(routeImg)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write(data)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func savePersonToFile(person flights.Person, dir string) error {
	dirName := fmt.Sprintf("%s/content/passengers", dir)
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.md", dirName, wvuflights.ToSlug(person.Name)))
	if err != nil {
		return err
	}

	defer f.Close()

	var (
		totalDistance     = 0.0
		totalCost         = 0.0
		totalSoloTrips    = 0
		totalSoloDistance = 0.0
		totalSoloCost     = 0.0

		personTrips     []flights.Trip
		personTripsData []tripData

		_max_date = "0000-00-00"
		lastDept  = ""
		lastVPDiv = ""
	)

	db.Preload(clause.Associations).Where("id IN (SELECT trip_id FROM trip_passengers WHERE person_name = ?)", person.Name).Find(&personTrips)

	for _, trip := range personTrips {
		dateRange := strings.Split(trip.Date, ";")
		endDate := ""
		if len(dateRange) > 1 {
			endDate = dateRange[1]
		}
		dist, err := al.RouteToDistance(trip.Route)
		if err != nil {
			return err
		}
		totalDistance += dist
		totalCost += trip.GetTotalCost()
		if len(trip.Passengers) == 1 {
			totalSoloTrips++
			totalSoloDistance += dist
			totalSoloCost += trip.GetTotalCost()
		}

		if endDate > _max_date || dateRange[0] > _max_date {
			_max_date = max(endDate, dateRange[0])
			for _, p := range trip.Passengers {
				if p.PersonName == person.Name && p.Department != "" && p.VPDiv != "" {
					lastDept = p.Department
					lastVPDiv = p.VPDiv
				}
			}
		}

		personTripsData = append(personTripsData, tripData{
			ID:          trip.ID,
			StartDate:   dateRange[0],
			EndDate:     endDate,
			Route:       trip.Route,
			Aircraft:    trip.Aircraft,
			Distance:    dist,
			Carbon:      trip.CarbonFootprint,
			RegNo:       trip.RegNo,
			NumPax:      trip.NumPax,
			Passengers:  trip.Passengers,
			FlightHours: trip.FlightHours,
			TotalCost:   trip.GetTotalCost(),
			Cost: map[string]float64{
				"fuel":           trip.Fuel,
				"landing_fees":   trip.Landing,
				"crew_expense":   trip.CrewExpense,
				"domestic_tax":   trip.DomTax,
				"hourly_rate":    trip.HourlyRate,
				"billing_amount": trip.BillingAmount,
			},
			Justification: trip.Justification,
		})
	}

	data, err := yaml.Marshal(peopleData{
		Name:              person.Name,
		PersonType:        person.Type,
		LastDept:          lastDept,
		LastVPDiv:         lastVPDiv,
		Trips:             personTripsData,
		TotalTrips:        len(personTrips),
		TotalTripDistance: totalDistance,
		TotalTripCost:     totalCost,
		TotalSoloTrips:    totalSoloTrips,
		TotalSoloDistance: totalSoloDistance,
		TotalSoloCost:     totalSoloCost,
	})
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(f, "---\n%s\n---\n", data)
	if err != nil {
		return err
	}

	return nil
}

func saveInvoiceToFile(invoice flights.Invoice, dir string) error {
	dirName := fmt.Sprintf("%s/static/invoices", dir)
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.pdf", dirName, invoice.Name))
	if err != nil {
		return err
	}

	_, err = f.Write(invoice.PDF)
	if err != nil {
		return err
	}

	f.Close()

	dirName = fmt.Sprintf("%s/content/invoices", dir)
	err = os.MkdirAll(dirName, 0755)
	if err != nil {
		return err
	}

	f, err = os.Create(fmt.Sprintf("%s/%s.md", dirName, invoice.Name))
	if err != nil {
		return err
	}

	defer f.Close()

	var (
		trips     []tripData
		_trips    []flights.Trip
		totalcost float64
	)

	db.Preload(clause.Associations).
		Where("strftime('%Y-%m', CASE WHEN instr(date, ';') > 0 THEN substr(date, instr(date, ';') + 1) ELSE date END) = ?", invoice.Name).
		Find(&_trips)

	if len(_trips) > 0 {
		db.Table("trips").
			Select("SUM(fuel + landing + crew_expense + dom_tax + billing_amount) as total_cost").
			Where("strftime('%Y-%m', CASE WHEN instr(date, ';') > 0 THEN substr(date, instr(date, ';') + 1) ELSE date END) = ?", invoice.Name).
			Scan(&totalcost)
	}

	for _, trip := range _trips {
		dateRange := strings.Split(trip.Date, ";")
		endDate := ""
		if len(dateRange) > 1 {
			endDate = dateRange[1]
		}
		dist, err := al.RouteToDistance(trip.Route)
		if err != nil {
			return err
		}

		var crewNames []string
		for _, crew := range trip.Crew {
			crewNames = append(crewNames, crew.Name)
		}

		trips = append(trips, tripData{
			ID:          trip.ID,
			StartDate:   dateRange[0],
			EndDate:     endDate,
			Route:       trip.Route,
			Aircraft:    trip.Aircraft,
			Distance:    dist,
			RegNo:       trip.RegNo,
			Passengers:  trip.Passengers,
			Crew:        crewNames,
			NumPax:      trip.NumPax,
			FlightHours: trip.FlightHours,
			TotalCost:   trip.GetTotalCost(),
			Cost: map[string]float64{
				"fuel":           trip.Fuel,
				"landing_fees":   trip.Landing,
				"crew_expense":   trip.CrewExpense,
				"domestic_tax":   trip.DomTax,
				"hourly_rate":    trip.HourlyRate,
				"billing_amount": trip.BillingAmount,
			},
			Justification: trip.Justification,
		})
	}

	data, err := yaml.Marshal(map[string]interface{}{
		"name":        invoice.Name,
		"date":        invoice.Name,
		"invoicetype": invoice.Type,
		"linkto":      fmt.Sprintf("/invoices/%s.pdf", invoice.Name),
		"totaltrips":  len(trips),
		"totalcost":   totalcost,
		"trips":       trips,
	})

	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(f, "---\n%s\n---\n", data)

	return err
}
