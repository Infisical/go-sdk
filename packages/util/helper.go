package util

import "strings"

func AppendAPIEndpoint(siteUrl string) string {
	// Ensure the address does not already end with "/api"
	if strings.HasSuffix(siteUrl, "/api") {
		return siteUrl
	}

	// Check if the address ends with a slash and append accordingly
	if siteUrl[len(siteUrl)-1] == '/' {
		return siteUrl + "api"
	}
	return siteUrl + "/api"
}
