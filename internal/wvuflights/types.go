package wvuflights

const (
	PersonTypeWVU = "wvu"
	PersonTypeLJA = "lja"
)

type Trip struct {
	ID              string `gorm:"primaryKey"`
	Date            string // "YYYY-MM-DD" or "YYYY-MM-DD;YYYY-MM-DD"
	Route           string // ex. "LBE-MGW-CRW-LBE" or "LBE-MGW-SFO-RON(2)-MGW-LBE"
	Aircraft        string
	RegNo           string // ex. "N12345"
	NumPax          int
	Passengers      []TripPassenger `gorm:"foreignKey:TripID"`
	Crew            []Person        `gorm:"many2many:trip_crew;"`
	Distance        float64
	CarbonFootprint int
	Fuel            float64
	Landing         float64
	CrewExpense     float64
	DomTax          float64
	FlightHours     float64
	HourlyRate      float64
	BillingAmount   float64
	Justification   string
	Notes           string
	Itinerary       []byte `yaml:"-"`
	Reservation     []byte `yaml:"-"`
	SignedBy        Person `gorm:"foreignKey:SignedByName" yaml:"-"`
	SignedByName    string
	ApprovedBy      Person `gorm:"foreignKey:ApprovedByName" yaml:"-"`
	ApprovedByName  string
}

type TripPassenger struct {
	TripID        string `gorm:"primaryKey" yaml:"-"`
	PaxNo         int
	Person        Person `gorm:"foreignKey:PersonName;references:Name" yaml:"-"`
	PersonName    string `gorm:"primaryKey"`
	Department    string
	VPDiv         string
	Justification string
	Code          string
}

type Person struct {
	Name string `gorm:"primaryKey"`
	Type string // "wvu" or "lja"
}

type Airport struct {
	Code     string `gorm:"primaryKey"`
	Name     string
	City     string
	State    string
	Lat, Lon float64
}

type Invoice struct {
	Name string `gorm:"primaryKey"`
	Type string // "month" or "trip"
	PDF  []byte `yaml:"-"`
}
