package customValidators

import (
	"net/url"
	"strconv"
	"strings"
)

var validGender = map[string]bool{
	"male":   true,
	"female": true,
}

func ParseLimitOffset(val string, defaultVal int) int {
	parsedVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	if parsedVal < 0 {
		return defaultVal
	}

	return parsedVal
}

func ParseGender(genderStr string) (string, bool) {
	if genderStr == "" {
		return "", true
	}

	_, isValid := validGender[genderStr]
	return genderStr, isValid
}

func ParseDepartmentID(id string) (int, bool) {
	if id == "" {
		return 0, true
	}
	departmentID, err := strconv.Atoi(id)
	if err != nil {
		return 0, false
	}

	if departmentID < 1 {
		return 0, false
	}

	return departmentID, true
}

func ParseURI(uri string) (string, bool) {
	// Check for empty string
	if strings.TrimSpace(uri) == "" {
		return "", false
	}

	// Parse the URI
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return "", false
	}

	// Check if scheme is present and valid
	if parsedURL.Scheme == "" {
		return "", false
	}

	// Validate host presence for network-based URIs
	if parsedURL.Scheme != "file" && parsedURL.Host == "" {
		return "", false
	}

	// Additional validation for specific schemes
	switch parsedURL.Scheme {
	case "http", "https":
		// For HTTP(S), ensure there's a valid host
		if !strings.Contains(parsedURL.Host, ".") && parsedURL.Host != "localhost" {
			return "", false
		}
	case "file":
		// For file scheme, ensure there's a path
		if parsedURL.Path == "" {
			return "", false
		}
	}

	return uri, true
}
