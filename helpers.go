package postmark

import (
	"fmt"
	"net/url"
)

// buildURLWithQuery constructs a URL with the given path and query parameters.
func buildURLWithQuery(path string, queryParams url.Values) string {
	if len(queryParams) == 0 {
		return path
	}
	return fmt.Sprintf("%s?%s", path, queryParams.Encode())
}

// buildURL constructs a URL with the given path and options.
func buildURL(path string, options map[string]interface{}) string {
	if options == nil {
		return path
	}

	values := &url.Values{}
	for k, v := range options {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	return buildURLWithQuery(path, *values)
}
