package wvuflights

func (t *Trip) GetReservation() []byte {
	return t.Reservation
}

func (t *Trip) GetItinerary() []byte {
	return t.Itinerary
}

func (t *Trip) GetTotalCost() float64 {
	return t.BillingAmount + t.Fuel + t.DomTax + t.Landing + t.CrewExpense
}

// GetCarbonFootprint returns the carbon footprint of the trip in total grams of CO2 emitted
func (t *Trip) GetCarbonFootprint() int {
	// 8,939.50 gCO2 emitted per gallon of jet fuel burned according to
	// the US Energy Information Administration (EIA)
	// https://www.eia.gov/environment/emissions/co2_vol_mass.php
	avgEmissions := 8939.50

	// flight hours (h) * fuel burn (gal/h) * avg emissions (gCO2/gal) = gCO2
	// rounded to 2 decimal places
	if fb := GetAircraftFuelBurn(t.Aircraft); fb != -1 {
		// return math.Round(((float64(fb)*t.FlightHours)*avgEmissions)*math.Pow(10, 2)) / math.Pow(10, 2)
		return int((float64(fb) * t.FlightHours) * avgEmissions)
	}

	return -1
}
