package wvuflights

import (
	"regexp"
	"strings"
)

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
