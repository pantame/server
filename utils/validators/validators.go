package validators

import (
	"net"
	"regexp"
	"strings"
)

func IsUsername(username string) bool {
	if len(username) < 1 || len(username) > 120 {
		return false
	}

	return regexp.MustCompile("^[a-zA-Z0-9_]*$").MatchString(username)
}

func IsStatus(status string) bool {
	switch status {
	case "pub", "pri", "del":
		return true
	default:
		return false
	}
}

func IsValidString(value string, min, max int) bool {
	if len(value) < min || (max > -1 && len(value) > max) {
		return false
	}
	return true
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

func IsValidEmailByMX(email string) bool {
	if !IsValidEmail(email) {
		return false
	}
	parts := strings.Split(email, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}
