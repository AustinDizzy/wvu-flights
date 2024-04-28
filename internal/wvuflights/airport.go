package wvuflights

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const (
	AIRPORTS_CSV_URL = "https://davidmegginson.github.io/ourairports-data/airports.csv"
)

var (
	ErrRemainOvernight = fmt.Errorf("overnight stay")
	ErrUnknownAirport  = fmt.Errorf("unknown airport")
	ErrCSVFormat       = fmt.Errorf("CSV format not as expected")
)

type AirportLookup struct {
	data              [][]string
	mapboxAccessToken string
}

func NewAirportLookup(opt ...string) (*AirportLookup, error) {
	u := AIRPORTS_CSV_URL
	var mbt string

	if len(opt) > 0 {
		mbt = opt[0]

		if len(opt) > 1 {
			u = opt[1]
		}
	}

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return &AirportLookup{
		data:              data,
		mapboxAccessToken: mbt,
	}, nil
}

func (a *AirportLookup) ToNames(route string) (map[string]string, error) {
	route = strings.ToUpper(route)
	segments := strings.Split(route, "-")

	names := make(map[string]string)

	for _, s := range segments {
		if _, ok := names[s]; ok {
			continue
		}

		airport, err := a.CodeToData(s)
		if err == ErrRemainOvernight {
			continue
		} else if err != nil {
			return nil, err
		} else {
			switch s {
			case "BKW":
				names[s] = "Beckley, WV"
			default:
				names[s] = fmt.Sprintf("%s, %s", airport.City, airport.State)
			}
		}
	}

	return names, nil
}

func (a *AirportLookup) RouteToImage(route string) ([]byte, error) {
	route = CleanRouteStr(route)
	segments := strings.Split(route, "-")
	if len(segments) < 2 {
		return nil, fmt.Errorf("invalid route")
	}
	locMap := make(map[string][]float64)
	for _, s := range segments {
		if _, ok := locMap[s]; ok {
			continue
		}
		airport, err := a.CodeToData(s)
		if err == ErrRemainOvernight {
			continue
		} else if err != nil {
			return nil, err
		} else {
			locMap[s] = []float64{airport.Lon, airport.Lat}
		}
	}

	if len(locMap) < 2 {
		return nil, fmt.Errorf("invalid route")
	}

	baseURL := "https://api.mapbox.com/styles/v1/mapbox/streets-v12/static"

	var path string
	for _, s := range segments {
		path += fmt.Sprintf("[%f,%f],", locMap[s][0], locMap[s][1])
	}
	path = strings.TrimSuffix(path, ",")

	geojsonRoute := fmt.Sprintf(`{"type":"LineString","coordinates":[%s]}`, path)

	url := fmt.Sprintf(`%s/geojson(%s)/auto/768x320@2x?padding=45&access_token=%s`, baseURL, geojsonRoute, a.mapboxAccessToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Read the response body
	return io.ReadAll(resp.Body)
}

func (a *AirportLookup) RouteToDistance(route string) (float64, error) {
	route = CleanRouteStr(route)

	segments := strings.Split(route, "-")

	if len(segments) < 2 {
		return 0, fmt.Errorf("invalid route")
	}

	var distance = 0.0

	for i, s := range segments {
		if i == 0 {
			continue
		}

		a1, err := a.CodeToData(segments[i-1])
		if err != nil {
			return 0, err
		}

		a2, err := a.CodeToData(s)
		if err != nil {
			return 0, err
		}

		distance += a1.DistanceTo(a2)
	}

	return distance, nil
}

func (a *AirportLookup) CodeToData(code string) (*Airport, error) {
	if strings.HasPrefix(code, "RON(") {
		return nil, ErrRemainOvernight
	}

	codeCol := -1
	municipalityCol := -1
	regionCol := -1
	latCol := -1
	lonCol := -1

	for i, header := range a.data[0] {
		switch header {
		case "ident":
			codeCol = i
		case "municipality":
			municipalityCol = i
		case "iso_region":
			regionCol = i
		case "latitude_deg":
			latCol = i
		case "longitude_deg":
			lonCol = i
		}
	}

	if codeCol == -1 || municipalityCol == -1 || regionCol == -1 || latCol == -1 || lonCol == -1 {
		return nil, ErrCSVFormat
	}

	for _, record := range a.data[1:] {
		if record[codeCol] == fmt.Sprintf("K%s", strings.ToUpper(code)) {
			city := record[municipalityCol]
			state := strings.Split(record[regionCol], "-")[1]

			lat, err := strconv.ParseFloat(record[latCol], 64)
			if err != nil {
				return nil, err
			}
			lon, err := strconv.ParseFloat(record[lonCol], 64)
			if err != nil {
				return nil, err
			}

			return &Airport{
				Code:  code,
				City:  city,
				State: state,
				Lat:   lat,
				Lon:   lon,
			}, nil
		}
	}

	fmt.Printf("unk code: %s\n", code)

	return nil, ErrUnknownAirport
}

func (a *Airport) DistanceTo(dest *Airport) float64 {
	const earthRadiusNm = 3440.065 // Earth radius in nautical miles

	// convert degrees to radians
	lat1 := degToRad(a.Lat)
	long1 := degToRad(a.Lon)
	lat2 := degToRad(dest.Lat)
	long2 := degToRad(dest.Lon)

	dlat := lat2 - lat1
	dlong := long2 - long1

	x := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlong/2)*math.Sin(dlong/2)

	calcDistance := 2 * math.Atan2(math.Sqrt(x), math.Sqrt(1-x))

	return earthRadiusNm * calcDistance
}

// degToRad converts decimal degrees to radians.
func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
