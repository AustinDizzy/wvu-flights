package wvuflights

import (
	"regexp"
	"strings"
)

func GetAircraftFuelBurn(aircraft string) int {
	// gallons per hour of fuel burn by aircraft type
	// from https://jetadvisors.com/
	fuelBurn := map[string]int{
		"Citation Excel":     248,
		"Citation V Ultra":   195,
		"Citation Encore":    180,
		"Citation Bravo":     171,
		"Citation CJ3":       170,
		"Citation CJ4":       218,
		"Citation Sovereign": 283,
		"Challenger 300":     304,
		"Challenger 350":     302,
		"Hawker 900XP":       261,
		"Gulfstream 280":     274,
		"Citation Latitude":  295,
		"King Air B200":      102,
		"Citation XLS+":      241,
	}

	if fb, ok := fuelBurn[aircraft]; ok {
		return fb
	}

	return -1
}

func ToSlug(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9']+`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`^-+|-+$`).ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "'", "")
	return s
}

func CleanRouteStr(route string) string {
	route = strings.ToUpper(route)
	route = regexp.MustCompile(`-RON\([0-9]+\)`).ReplaceAllString(route, "")
	route = regexp.MustCompile(`^RON\([0-9]+\)-`).ReplaceAllString(route, "")

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
