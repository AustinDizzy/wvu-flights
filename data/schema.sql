CREATE TABLE `people` (
	`name` text, -- the person's name
	`type` text, -- the person's type, either "wvu" or "lja"
	PRIMARY KEY (`name`) -- there can't be two people with the same name, full names are used from the trip itineraries where possible
);
CREATE TABLE `trip_crew` (`trip_id` text,`person_name` text,PRIMARY KEY (`trip_id`,`person_name`),CONSTRAINT `fk_trip_crew_trip` FOREIGN KEY (`trip_id`) REFERENCES `trips`(`id`),CONSTRAINT `fk_trip_crew_person` FOREIGN KEY (`person_name`) REFERENCES `people`(`name`));
CREATE TABLE `invoices` (
	`name` text, -- the name of the invoice, in the format YYYY-MM for a month invoice, or in the date format for a trip invoice (YYYY-MM-DD;YYYY-MM-DD, or YYYY-MM-DD)
	`type` text, -- the type of invoice ("month" or "trip")
	`pdf` blob,  -- a blob containing the PDF page data for the invoice
	PRIMARY KEY (`name`)
);
CREATE TABLE IF NOT EXISTS "trip_passengers" (
	"trip_id"	text, -- the trip ID the passenger was on
	"person_name"	text, -- the name of the passenger
	"department"	text, -- the Department column value for this passenger on the reservation form
	"vp_div"	text, -- the VP/Division column value for this passenger on the reservation form
	"justification"	text, -- the Justification column value for this passenger on the reservation form, if this is blank then the trip's justification should be used
	"code"	text, -- the Code column value for this passenger on the reservation form
	"pax_no"	integer, -- the number this passenger appears in the list of passengers on the reservation form
	PRIMARY KEY("trip_id","person_name"), -- the same person can't appear on the same trip more than once
	CONSTRAINT "fk_trip_passengers_person" FOREIGN KEY("person_name") REFERENCES "people"("name"), -- the person must exist in the people table
	CONSTRAINT "fk_trips_passengers" FOREIGN KEY("trip_id") REFERENCES "trips"("id") -- the trip must exist in the trips table
);
CREATE TABLE IF NOT EXISTS "trips" (
	"id"	text, -- the trip ID, set by the vendor found on the trip itinerary
	"date"	text, -- the date of the trip, in the format YYYY-MM-DD for a single day trip or YYYY-MM-DD;YYYY-MM-DD for a multi-day trip, found on the invoice
	"route"	text, -- the route of the trip, in a list of IATA airport codes separated by hyphens (ex. LBE-MGW-IAD, with RON meaning 'remain overnight'), found on the invoice
	"aircraft"	text, -- the aircraft used for the trip, found on the invoice or trip itinerary
	"reg_no"	text, -- the registration number of the aircraft used for the trip, found on the invoice or trip itinerary
	"num_pax"	integer, -- the number of passengers on the trip, found on the invoice or trip itinerary
	"justification"	text, -- the justification for the trip written by WVU administration, found on the reservation form
	"distance"	real, -- the distance of the trip in nautical miles, calculated from the route
	"fuel"	real, -- the fuel surcharge for the trip, found on the invoice
	"landing"	real, -- the landing / deicing fees for the trip, found on the invoice
	"crew_expense"	real, -- the crew expenses for the trip, found on the invoice
	"dom_tax"	real, -- the domestic segment tax for the trip, found on the invoice
	"flight_hours"	real, -- the flight hours for the trip, found on the invoice
	"hourly_rate"	real, -- the hourly rate for the trip, found on the invoice
	"billing_amount"	real, -- the total amount billed for the trip, found on the invoice
	"signed_by_name"	text, -- the name of the person who signed the reservation form
	"approved_by_name"	text, -- the name of the person in the president's office who is marked as approving the trip, found on the reservation form
	"carbon_footprint"	integer, -- the carbon footprint for the trip in total grams of CO2, calculated from the distance, aircraft type, and EIA data
	"notes"	text, -- any notes about the trip, found on the invoice or compiled during processing
	"itinerary"	blob, -- a blob containing the PDF page data for the trip's itinerary
	"reservation"	blob, -- a blob containing the PDF page data for the trip's reservation form
	PRIMARY KEY("id"), -- the trip ID is unique
	CONSTRAINT "fk_trips_signed_by" FOREIGN KEY("signed_by_name") REFERENCES "people"("name"), -- the person who signed the reservation form must exist in the people table
	CONSTRAINT "fk_trips_approved_by" FOREIGN KEY("approved_by_name") REFERENCES "people"("name") -- the person who approved the trip must exist in the people table
);