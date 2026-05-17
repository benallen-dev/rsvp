package handlers

import (
	"net/url"
	"regexp"
)

var modeRegex = regexp.MustCompile(`^/(day|evening)(?:/|$)`)

func getRouteType(u *url.URL) string {
	matches := modeRegex.FindStringSubmatch(u.Path)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
