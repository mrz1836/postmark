package postmark

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

// TestRouter provides simple HTTP routing for testing purposes.
// It supports method-based routing and path parameter extraction.
type TestRouter struct {
	routes []route
}

type route struct {
	method  string
	pattern string
	handler http.HandlerFunc
	regex   *regexp.Regexp
	params  []string
}

// NewTestRouter creates a new test router.
func NewTestRouter() *TestRouter {
	return &TestRouter{
		routes: make([]route, 0),
	}
}

// ServeHTTP implements http.Handler interface.
func (tr *TestRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, r := range tr.routes {
		if r.method != req.Method {
			continue
		}

		if tr.matchRoute(r, req, w) {
			return
		}
	}

	// No route found
	http.NotFound(w, req)
}

// HandleFunc registers a handler function for the given method and pattern.
func (tr *TestRouter) HandleFunc(method, pattern string, handler http.HandlerFunc) {
	r := route{
		method:  method,
		pattern: pattern,
		handler: handler,
	}

	// Check if pattern contains parameters (e.g., :domainID)
	if strings.Contains(pattern, ":") {
		r.regex, r.params = compilePattern(pattern)
	}

	tr.routes = append(tr.routes, r)
}

// Get registers a GET handler for the given pattern.
func (tr *TestRouter) Get(pattern string, handler http.HandlerFunc) {
	tr.HandleFunc(http.MethodGet, pattern, handler)
}

// Post registers a POST handler for the given pattern.
func (tr *TestRouter) Post(pattern string, handler http.HandlerFunc) {
	tr.HandleFunc(http.MethodPost, pattern, handler)
}

// Put registers a PUT handler for the given pattern.
func (tr *TestRouter) Put(pattern string, handler http.HandlerFunc) {
	tr.HandleFunc(http.MethodPut, pattern, handler)
}

// Delete registers a DELETE handler for the given pattern.
func (tr *TestRouter) Delete(pattern string, handler http.HandlerFunc) {
	tr.HandleFunc(http.MethodDelete, pattern, handler)
}

// Patch registers a PATCH handler for the given pattern.
func (tr *TestRouter) Patch(pattern string, handler http.HandlerFunc) {
	tr.HandleFunc(http.MethodPatch, pattern, handler)
}

// matchRoute attempts to match a route and execute its handler.
func (tr *TestRouter) matchRoute(route route, r *http.Request, w http.ResponseWriter) bool {
	if route.regex == nil {
		return tr.matchExactRoute(route, r, w)
	}
	return tr.matchPatternRoute(route, r, w)
}

// matchExactRoute handles routes without parameters.
func (tr *TestRouter) matchExactRoute(route route, r *http.Request, w http.ResponseWriter) bool {
	if route.pattern == r.URL.Path {
		route.handler(w, r)
		return true
	}
	return false
}

// matchPatternRoute handles routes with path parameters.
func (tr *TestRouter) matchPatternRoute(route route, r *http.Request, w http.ResponseWriter) bool {
	matches := route.regex.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		return false
	}

	// Extract path parameters and add to request context
	if len(matches) > 1 && len(route.params) > 0 {
		ctx := r.Context()
		for i, param := range route.params {
			if i+1 < len(matches) {
				ctx = context.WithValue(ctx, paramKey(param), matches[i+1])
			}
		}
		r = r.WithContext(ctx)
	}
	route.handler(w, r)
	return true
}

// compilePattern converts a pattern like "/domains/:domainID/verify"
// to a regex and extracts parameter names.
func compilePattern(pattern string) (*regexp.Regexp, []string) {
	var params []string
	regexPattern := "^"

	parts := strings.Split(pattern, "/")
	for _, part := range parts {
		if part == "" {
			continue
		}

		if strings.HasPrefix(part, ":") {
			// Parameter segment
			paramName := part[1:] // Remove the ':'
			params = append(params, paramName)
			regexPattern += "/([^/]+)"
		} else {
			// Literal segment
			regexPattern += "/" + regexp.QuoteMeta(part)
		}
	}

	regexPattern += "$"

	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		panic("invalid pattern: " + pattern)
	}

	return regex, params
}

// paramKey is used as the key type for storing path parameters in request context.
type paramKey string

// GetPathParam extracts a path parameter from the request context.
// This is only available for handlers registered with parameterized routes.
func GetPathParam(r *http.Request, key string) string {
	if value := r.Context().Value(paramKey(key)); value != nil {
		return value.(string)
	}
	return ""
}
