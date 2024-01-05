package utils

import "net/url"

func ValidateURL(urlString string) (bool, string) {
	parsedURL, err := url.Parse(urlString)

	if err != nil || parsedURL == nil {
		return false, "Trouble Parsing URL"
	}

	if parsedURL.Host == "" {
		return false, "No Host"
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false, "Invalid URL Scheme"
	}

	return true, ""
}
