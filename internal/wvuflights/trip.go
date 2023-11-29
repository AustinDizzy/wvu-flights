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
