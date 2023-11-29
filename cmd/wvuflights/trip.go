package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"

	flights "github.com/austindizzy/wvu-flights/internal/wvuflights"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v3"
)

var (
	tripAddCmd = &cli.Command{
		Name:  "add",
		Usage: "Add a trip",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "id",
				Usage: "the ID of the trip",
			},
			&cli.StringFlag{
				Name:  "date",
				Usage: "the date of the trip",
			},
			&cli.StringFlag{
				Name:  "route",
				Usage: "the route of the trip",
			},
			&cli.StringFlag{
				Name:  "aircraft",
				Usage: "the aircraft used on the trip",
			},
			&cli.StringFlag{
				Name:  "reg",
				Usage: "the registration number of the aircraft used on the trip",
			},
			&cli.BoolFlag{
				Name:  "missing-res",
				Usage: "add trip with missing a reservation form (no justification, no approval, no signature, no passenger titles)",
				Value: false,
			},
		},
		Action: tripAdd,
	}
)

func tripAdd(c *cli.Context) error {
	var (
		tripID       = c.String("id")
		tripDate     = c.String("date")
		tripRoute    = c.String("route")
		tripAircraft = c.String("aircraft")
		tripReg      = c.String("reg")
		err          error
		_count       int64
	)

	if tripID == "" {
		pmpt := promptui.Prompt{
			Label: "Trip ID",
			// use Validate to check if input is valid ID
			Validate: func(input string) error {
				re := regexp.MustCompile(`^X?\d+$`)
				if !re.MatchString(input) {
					return errors.New("invalid ID format")
				}
				return nil
			},
		}

		tripID, err = pmpt.Run()
		if err != nil {
			return err
		}
	}

	// return error if trip with tripID exists
	db.Model(&flights.Trip{}).Where("id = ?", tripID).Count(&_count)
	if _count > 0 {
		return errors.New("trip with ID already exists")
	}

	if tripDate == "" {
		pmpt := promptui.Prompt{
			Label: "Trip Date",
			// use Validate to check if input is YYYY-MM-DD or YYYY-MM-DD;YYYY-MM-DD
			Validate: func(input string) error {
				re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2};)?\d{4}-\d{2}-\d{2}$`)
				if !re.MatchString(input) {
					return errors.New("invalid date format")
				}
				return nil
			},
		}

		tripDate, err = pmpt.Run()
		if err != nil {
			return err
		}
	}

	// log warning if a trip with tripDate exists
	_count = 0
	db.Model(&flights.Trip{}).Where("date = ?", tripDate).Count(&_count)
	if _count > 0 {
		log.Printf("warning: a trip with date \"%s\" already exists", tripDate)
	}

	if tripRoute == "" {
		pmpt := promptui.Prompt{
			Label: "Trip Route",
			// use Validate to check if input is valid route
			Validate: func(input string) error {
				re := regexp.MustCompile(`^([A-Z]{3}(\(\d\))?(-|\s))+[A-Z]{3}$`)
				valid := []string{"charleston", "crw", "lewisburg", "lwb", "beckley", "bkw"}
				if !re.MatchString(input) && !slices.Contains(valid, strings.ToLower(input)) {
					return errors.New("invalid route format")
				}
				return nil
			},
		}

		tripRoute, err = pmpt.Run()
		if err != nil {
			return err
		}

		switch strings.ToLower(tripRoute) {
		case "charleston":
		case "crw":
			tripRoute = "LBE-MGW-CRW-MGW-LBE"
		case "lewisburg":
		case "lwb":
			tripRoute = "LBE-MGW-LWB-MGW-LBE"
		case "beckley":
		case "bkw":
			tripRoute = "LBE-MGW-BKW-MGW-LBE"
		default:
			tripRoute = strings.ReplaceAll(tripRoute, " ", "-")
		}
	}

	if tripAircraft == "" {
		var availAircraft []string
		db.Model(&flights.Trip{}).Distinct("aircraft").Order("aircraft").Pluck("aircraft", &availAircraft)

		tripAircraft, _, err = promptWithNew(availAircraft, "Aircraft")
		if err != nil {
			return err
		}
	}

	if tripReg == "" {
		var availReg []string
		db.Model(&flights.Trip{}).Distinct("reg_no").Order("reg_no DESC").Pluck("reg_no", &availReg)

		tripReg, _, err = promptWithNew(availReg, "Registration Number")
		if err != nil {
			return err
		}

		tripReg = strings.ToUpper(tripReg)
	}

	addPaxPmpt := promptui.Select{
		Label: "Add Passengers?",
		Items: []string{"Yes", "No"},
	}

	passengers := make([]flights.TripPassenger, 0)
	paxNo := 1

	var availPass []string
	pplSql := "SELECT p.name FROM people"
	pplSql += " p LEFT JOIN trip_passengers tp ON p.name = tp.person_name"
	pplSql += " WHERE p.type = 'wvu'"
	pplSql += " GROUP BY p.name ORDER BY COUNT(tp.person_name) DESC;"
	db.Raw(pplSql).Pluck("name", &availPass)

	for {
		_, addPax, err := addPaxPmpt.Run()
		if err != nil {
			return err
		}

		if addPax == "Yes" {
			paxName, isNew, err := promptWithNew(availPass, "Passenger")
			if err != nil {
				return err
			}

			if isNew {
				db.Create(&flights.Person{
					Name: paxName,
					Type: flights.PersonTypeWVU,
				})
			}

			var dept, vpDiv string
			if !c.Bool("missing-res") {
				var depts []string
				db.Raw(`SELECT department FROM trip_passengers GROUP BY department ORDER BY count(*) DESC;`).Pluck("department", &depts)

				dept, _, err = promptWithNew(depts, "Department")
				if err != nil {
					return err
				}

				var vpDivs []string
				db.Raw(`SELECT vp_div FROM trip_passengers GROUP BY vp_div ORDER BY count(*) DESC;`).Pluck("vp_div", &vpDivs)

				vpDiv, _, err = promptWithNew(vpDivs, "VP Div")
				if err != nil {
					return err
				}
			}

			var p flights.Person
			db.Where("name = ?", paxName).First(&p)

			passengers = append(passengers, flights.TripPassenger{
				Person:     p,
				PaxNo:      paxNo,
				Department: dept,
				VPDiv:      vpDiv,
			})
			paxNo++
		}

		if addPax == "No" {
			break
		}
	}

	var availCrew []string
	db.Model(&flights.Person{}).Where("type = ?", flights.PersonTypeLJA).Order("name ASC").Distinct("name").Pluck("name", &availCrew)

	crew := make([]flights.Person, 0)

	addCrewPmpt := promptui.Select{
		Label: "Add Crew?",
		Items: []string{"Yes", "No"},
	}

	firstRun := true

	for {
		var (
			addCrew string
			err     error
		)

		if !firstRun {
			_, addCrew, err = addCrewPmpt.Run()
			if err != nil {
				return err
			}
		}

		if addCrew == "Yes" || firstRun {
			firstRun = false
			crewName, isNew, err := promptWithNew(availCrew, "Crew Member")
			if err != nil {
				return err
			}

			if isNew {
				db.Create(&flights.Person{
					Name: crewName,
					Type: flights.PersonTypeLJA,
				})
				db.Model(&flights.Person{}).Distinct("name").Order("name ASC").Where("type = ?", flights.PersonTypeLJA).Pluck("name", &availCrew)
			}

			var c flights.Person
			db.Where("name = ?", crewName).First(&c)

			crew = append(crew, c)
		}

		if addCrew == "No" {
			break
		}
	}

	numPax := -1
	pmpt := promptui.Prompt{
		Label: "Num. of Passengers",
		Validate: func(input string) error {
			re := regexp.MustCompile(`^\d+$`)
			if !re.MatchString(input) {
				return errors.New("invalid number format")
			}
			return nil
		},
	}
	numPaxStr, err := pmpt.Run()
	if err != nil {
		return err
	}

	numPax, _ = strconv.Atoi(numPaxStr)

	var validNum = func(input string) error {
		re := regexp.MustCompile(`^(\d{1,3}(,\d{3})*|(\d+))(\.\d{2})?$`)
		if !re.MatchString(input) {
			return errors.New("invalid number format")
		}
		return nil
	}

	var fuelSurcharge float64
	pmpt = promptui.Prompt{
		Label:    "Fuel Surcharge",
		Validate: validNum,
	}

	fuelSurchargeStr, err := pmpt.Run()
	if err != nil {
		return err
	}

	fuelSurchargeStr = strings.ReplaceAll(fuelSurchargeStr, ",", "")
	fuelSurcharge, _ = strconv.ParseFloat(fuelSurchargeStr, 64)

	var landingFees float64
	pmpt = promptui.Prompt{
		Label:    "Landing Fees",
		Validate: validNum,
	}
	landingFeesStr, err := pmpt.Run()
	if err != nil {
		return err
	}

	landingFeesStr = strings.ReplaceAll(landingFeesStr, ",", "")
	landingFees, _ = strconv.ParseFloat(landingFeesStr, 64)

	var crewExpense float64
	pmpt = promptui.Prompt{
		Label:    "Crew Expenses",
		Validate: validNum,
	}

	crewExpenseStr, err := pmpt.Run()
	if err != nil {
		return err
	}

	crewExpenseStr = strings.ReplaceAll(crewExpenseStr, ",", "")
	crewExpense, _ = strconv.ParseFloat(crewExpenseStr, 64)

	var domTax float64 = 0
	pmpt = promptui.Prompt{
		Label: "Domestic Tax",
	}

	domTaxStr, err := pmpt.Run()
	if err != nil {
		return err
	}

	if domTaxStr != "" {
		domTaxStr = strings.ReplaceAll(domTaxStr, ",", "")
		domTax, _ = strconv.ParseFloat(domTaxStr, 64)
	}

	var flightHours float64
	pmpt = promptui.Prompt{
		Label: "Flight Hours",
		// use Validate to see if input is in valid number format 2.0
		Validate: func(s string) error {
			re := regexp.MustCompile(`^\d+(\.\d+)?$`)
			if !re.MatchString(s) {
				return errors.New("invalid number format")
			}
			return nil
		},
	}

	flightHoursStr, err := pmpt.Run()
	if err != nil {
		return err
	}

	flightHours, _ = strconv.ParseFloat(flightHoursStr, 64)

	var hourlies []string
	db.Model(&flights.Trip{}).Distinct("hourly_rate").Order("hourly_rate ASC").Pluck("hourly_rate", &hourlies)
	hrlyPmpt := promptui.SelectWithAdd{
		Label:    "Hourly Rate",
		Items:    hourlies,
		AddLabel: "+ New Hourly Rate",
	}

	_, hourlyRateStr, err := hrlyPmpt.Run()
	if err != nil {
		return err
	}

	var hourlyRate float64
	hourlyRateStr = strings.ReplaceAll(hourlyRateStr, ",", "")
	hourlyRate, _ = strconv.ParseFloat(hourlyRateStr, 64)

	var billingAmt = hourlyRate * flightHours
	amtPmpt := promptui.Prompt{
		Label:    "Billing Amount",
		Validate: validNum,
		Default:  strconv.FormatFloat(billingAmt, 'f', 2, 64),
	}

	billingAmtStr, err := amtPmpt.Run()
	if err != nil {
		return err
	}

	billingAmtStr = strings.ReplaceAll(billingAmtStr, ",", "")
	billingAmt, _ = strconv.ParseFloat(billingAmtStr, 64)

	var signer, approver flights.Person

	if !c.Bool("missing-res") {
		var signedBy string
		sigPmpt := promptui.SelectWithAdd{
			Label:    "Signed By",
			Items:    prioritize(availPass, "Amy Garbrick", "Maryanne Reed", "Melissa A. Patterson"),
			AddLabel: "+ New WVU Member",
		}

		i, signedBy, err := sigPmpt.Run()
		if err != nil {
			return err
		}

		if i == -1 {
			db.Create(&flights.Person{
				Name: signedBy,
				Type: flights.PersonTypeWVU,
			})
		}

		db.Where("name = ?", signedBy).First(&signer)

		var approvedBy string
		appPmpt := promptui.SelectWithAdd{
			Label:    "Approved By",
			Items:    prioritize(availPass, "Amy Garbrick"),
			AddLabel: "+ New WVU Member",
		}

		i, approvedBy, err = appPmpt.Run()
		if err != nil {
			return err
		}

		if i == -1 {
			db.Create(&flights.Person{
				Name: approvedBy,
				Type: flights.PersonTypeWVU,
			})
		}

		db.Where("name = ?", approvedBy).First(&approver)
	}

	trip := flights.Trip{
		ID:            tripID,
		Date:          tripDate,
		Route:         tripRoute,
		Aircraft:      tripAircraft,
		RegNo:         tripReg,
		NumPax:        numPax,
		Passengers:    passengers,
		Crew:          crew,
		Fuel:          fuelSurcharge,
		Landing:       landingFees,
		CrewExpense:   crewExpense,
		DomTax:        domTax,
		FlightHours:   flightHours,
		HourlyRate:    hourlyRate,
		BillingAmount: billingAmt,
		ApprovedBy:    approver,
		SignedBy:      signer,
	}

	if err := db.Create(&trip).Error; err != nil {
		return err
	}

	log.Printf("trip added: %#v", trip)

	return nil
}

func prioritize(list []string, target ...string) []string {
	var newList []string
	for _, item := range list {
		if !slices.Contains(target, item) {
			newList = append(newList, item)
		}
	}
	newList = append(target, newList...)
	return newList
}

func promptWithNew(items []string, typeStr string) (choice string, isNew bool, err error) {
	addStr := fmt.Sprintf("+ Add New %s", typeStr)
	items = prioritize(items, addStr)

	searcher := func(arr []string) func(input string, index int) bool {
		return func(input string, index int) bool {
			input = strings.ReplaceAll(strings.ToLower(input), " ", "")
			return strings.Contains(strings.ReplaceAll(strings.ToLower(arr[index]), " ", ""), input)
		}
	}

	prompt := promptui.Select{
		Label:    strings.ToTitle(typeStr),
		Items:    items,
		Searcher: searcher(items),
	}

	_, choice, err = prompt.Run()
	if err != nil {
		return
	}

	isNew = choice == addStr
	if isNew {
		addPrompt := promptui.Prompt{
			Label: fmt.Sprintf("New %s", typeStr),
		}

		choice, err = addPrompt.Run()
		if err != nil {
			return
		}

		return
	}

	return
}
