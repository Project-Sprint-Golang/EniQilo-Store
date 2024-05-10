package helper

import "net/url"

func ValidateURL(str string) bool {
	// Check if the URL is not empty
	if str == "" {
		return false
	}

	// Check if the URL has a scheme and host
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
