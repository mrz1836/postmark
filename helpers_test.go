package postmark

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testPath = "/messages/outbound"

func TestBuildURLWithQuery_EmptyParams(t *testing.T) {
	path := testPath
	params := url.Values{}

	result := buildURLWithQuery(path, params)
	assert.Equal(t, "/messages/outbound", result, "Empty params should return path unchanged")
}

func TestBuildURLWithQuery_SingleParam(t *testing.T) {
	path := testPath
	params := url.Values{}
	params.Add("count", "50")

	result := buildURLWithQuery(path, params)
	assert.Equal(t, "/messages/outbound?count=50", result)
}

func TestBuildURLWithQuery_MultipleParams(t *testing.T) {
	path := testPath
	params := url.Values{}
	params.Add("count", "50")
	params.Add("offset", "100")

	result := buildURLWithQuery(path, params)
	// Order might vary, so check both possibilities
	assert.Contains(t, result, "count=50")
	assert.Contains(t, result, "offset=100")
	assert.Contains(t, result, "/messages/outbound?")
}

func TestBuildURLWithQuery_SpecialCharacters(t *testing.T) {
	path := testPath
	params := url.Values{}
	params.Add("recipient", "test@example.com")
	params.Add("subject", "Hello World!")

	result := buildURLWithQuery(path, params)
	// URL encoding should happen
	assert.Contains(t, result, "/messages/outbound?")
	// @ should be encoded as %40
	assert.Contains(t, result, "%40")
	// Space should be encoded
	assert.Contains(t, result, "Hello+World")
}

func TestBuildURL_NilOptions(t *testing.T) {
	path := testPath
	result := buildURL(path, nil)
	assert.Equal(t, "/messages/outbound", result, "Nil options should return path unchanged")
}

func TestBuildURL_EmptyOptions(t *testing.T) {
	path := testPath
	options := map[string]interface{}{}

	result := buildURL(path, options)
	// Empty map should return path (buildURLWithQuery returns path when no params)
	assert.Equal(t, "/messages/outbound", result)
}

func TestBuildURL_StringValue(t *testing.T) {
	path := testPath
	options := map[string]interface{}{
		"messagestream": "outbound",
	}

	result := buildURL(path, options)
	assert.Contains(t, result, "messagestream=outbound")
}

func TestBuildURL_IntValue(t *testing.T) {
	path := testPath
	options := map[string]interface{}{
		"count": 50,
	}

	result := buildURL(path, options)
	assert.Contains(t, result, "count=50")
}

func TestBuildURL_BoolValue(t *testing.T) {
	path := testPath
	options := map[string]interface{}{
		"trackOpens": true,
	}

	result := buildURL(path, options)
	assert.Contains(t, result, "trackOpens=true")
}

func TestBuildURL_MultipleTypes(t *testing.T) {
	path := testPath
	options := map[string]interface{}{
		"count":         50,
		"offset":        100,
		"messagestream": "outbound",
		"trackOpens":    true,
	}

	result := buildURL(path, options)
	assert.Contains(t, result, "count=50")
	assert.Contains(t, result, "offset=100")
	assert.Contains(t, result, "messagestream=outbound")
	assert.Contains(t, result, "trackOpens=true")
	assert.Contains(t, result, "/messages/outbound?")
}

func TestBuildURL_FloatValue(t *testing.T) {
	path := "/stats"
	options := map[string]interface{}{
		"rate": 3.14,
	}

	result := buildURL(path, options)
	assert.Contains(t, result, "rate=3.14")
}

func TestBuildURL_SpecialCharactersInValue(t *testing.T) {
	path := "/search"
	options := map[string]interface{}{
		"query": "test@example.com",
		"name":  "John Doe",
	}

	result := buildURL(path, options)
	assert.Contains(t, result, "/search?")
	// Should be URL encoded
	assert.Contains(t, result, "%40") // @ encoded
}

func TestBuildURL_ZeroValues(t *testing.T) {
	path := "/messages"
	options := map[string]interface{}{
		"count":  0,
		"offset": 0,
		"flag":   false,
	}

	result := buildURL(path, options)
	assert.Contains(t, result, "count=0")
	assert.Contains(t, result, "offset=0")
	assert.Contains(t, result, "flag=false")
}
