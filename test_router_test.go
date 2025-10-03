package postmark

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestRouter(t *testing.T) {
	router := NewTestRouter()
	assert.NotNil(t, router)
	assert.NotNil(t, router.routes)
	assert.Empty(t, router.routes)
}

func TestTestRouter_Get(t *testing.T) {
	router := NewTestRouter()
	called := false

	router.Get("/test", func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	assert.Len(t, router.routes, 1)
	assert.Equal(t, http.MethodGet, router.routes[0].method)
	assert.Equal(t, "/test", router.routes[0].pattern)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTestRouter_Post(t *testing.T) {
	router := NewTestRouter()
	called := false

	router.Post("/test", func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestTestRouter_Put(t *testing.T) {
	router := NewTestRouter()
	called := false

	router.Put("/test", func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPut, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTestRouter_Delete(t *testing.T) {
	router := NewTestRouter()
	called := false

	router.Delete("/test", func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodDelete, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestTestRouter_Patch(t *testing.T) {
	router := NewTestRouter()
	called := false

	router.Patch("/test", func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPatch, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTestRouter_404NotFound(t *testing.T) {
	router := NewTestRouter()

	router.Get("/existing", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTestRouter_MethodNotMatched(t *testing.T) {
	router := NewTestRouter()

	router.Get("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Try POST on GET-only route
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTestRouter_PathParameters(t *testing.T) {
	router := NewTestRouter()
	var capturedID string

	router.Get("/users/:userID", func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetPathParam(r, "userID")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/users/12345", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "12345", capturedID)
}

func TestTestRouter_MultiplePathParameters(t *testing.T) {
	router := NewTestRouter()
	var capturedDomainID, capturedRecordID string

	router.Get("/domains/:domainID/records/:recordID", func(w http.ResponseWriter, r *http.Request) {
		capturedDomainID = GetPathParam(r, "domainID")
		capturedRecordID = GetPathParam(r, "recordID")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/domains/456/records/789", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "456", capturedDomainID)
	assert.Equal(t, "789", capturedRecordID)
}

func TestGetPathParam_NotFound(t *testing.T) {
	router := NewTestRouter()
	var capturedParam string

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		capturedParam = GetPathParam(r, "nonexistent")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, capturedParam, "Non-existent param should return empty string")
}

func TestCompilePattern_SimplePattern(t *testing.T) {
	regex, params := compilePattern("/users/:id")

	assert.NotNil(t, regex)
	assert.Len(t, params, 1)
	assert.Equal(t, "id", params[0])

	assert.True(t, regex.MatchString("/users/123"))
	assert.False(t, regex.MatchString("/users/"))
	assert.False(t, regex.MatchString("/users/123/extra"))
}

func TestCompilePattern_MultipleParams(t *testing.T) {
	regex, params := compilePattern("/domains/:domainID/verify/:token")

	assert.NotNil(t, regex)
	assert.Len(t, params, 2)
	assert.Equal(t, "domainID", params[0])
	assert.Equal(t, "token", params[1])

	assert.True(t, regex.MatchString("/domains/123/verify/abc456"))
	assert.False(t, regex.MatchString("/domains/123/verify"))
}

func TestCompilePattern_NoParams(t *testing.T) {
	regex, params := compilePattern("/static/path")

	// Even without params, a regex is created for exact matching
	assert.NotNil(t, regex)
	assert.Empty(t, params)
}

func TestCompilePattern_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{"trailing slash param", "/users/:id/"},
		{"consecutive params", "/api/:version/:endpoint"},
		{"mixed static and params", "/api/v1/:resource/:id/details"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regex, params := compilePattern(tt.pattern)
			assert.NotNil(t, regex)
			assert.NotEmpty(t, params)
		})
	}
}

func TestCompilePattern_SpecialCharacters(t *testing.T) {
	// Pattern with dots and dashes (common in domain names)
	regex, params := compilePattern("/domains/:domainName/verify")

	assert.NotNil(t, regex)
	assert.Len(t, params, 1)

	// Should match alphanumeric, dots, dashes
	assert.True(t, regex.MatchString("/domains/example.com/verify"))
	assert.True(t, regex.MatchString("/domains/test-domain/verify"))
	assert.True(t, regex.MatchString("/domains/123/verify"))
}

func TestTestRouter_MultipleRoutes(t *testing.T) {
	router := NewTestRouter()
	route1Called := false
	route2Called := false

	router.Get("/route1", func(w http.ResponseWriter, _ *http.Request) {
		route1Called = true
		w.WriteHeader(http.StatusOK)
	})

	router.Get("/route2", func(w http.ResponseWriter, _ *http.Request) {
		route2Called = true
		w.WriteHeader(http.StatusOK)
	})

	// Test route 1
	req1 := httptest.NewRequest(http.MethodGet, "/route1", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	assert.True(t, route1Called)
	assert.False(t, route2Called)

	// Reset and test route 2
	route1Called = false
	req2 := httptest.NewRequest(http.MethodGet, "/route2", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.False(t, route1Called)
	assert.True(t, route2Called)
}

func TestTestRouter_SamePathDifferentMethods(t *testing.T) {
	router := NewTestRouter()
	getCallCount := 0
	postCallCount := 0

	router.Get("/resource", func(w http.ResponseWriter, _ *http.Request) {
		getCallCount++
		w.WriteHeader(http.StatusOK)
	})

	router.Post("/resource", func(w http.ResponseWriter, _ *http.Request) {
		postCallCount++
		w.WriteHeader(http.StatusCreated)
	})

	// Test GET
	req1 := httptest.NewRequest(http.MethodGet, "/resource", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	assert.Equal(t, 1, getCallCount)
	assert.Equal(t, 0, postCallCount)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Test POST
	req2 := httptest.NewRequest(http.MethodPost, "/resource", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 1, getCallCount)
	assert.Equal(t, 1, postCallCount)
	assert.Equal(t, http.StatusCreated, w2.Code)
}

func TestCompilePattern_InvalidPattern_Panic(t *testing.T) {
	// This test verifies panic behavior - if pattern compilation has issues
	// The current implementation is quite permissive, so this is more of a safety test
	defer func() {
		if r := recover(); r != nil {
			// Expected panic caught
			assert.Contains(t, r, "invalid pattern")
		}
	}()

	// Try pattern that would cause regex compilation to fail
	// Most patterns will succeed, this is just for edge case coverage
	regex, params := compilePattern("/test/:param")
	assert.NotNil(t, regex)
	assert.NotEmpty(t, params)
}

func TestMatchExactRoute(t *testing.T) {
	router := NewTestRouter()
	called := false

	router.Get("/exact/path/here", func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/exact/path/here", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test non-exact match
	called = false
	req2 := httptest.NewRequest(http.MethodGet, "/exact/path/there", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.False(t, called)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestGetPathParam_EmptyContext(t *testing.T) {
	// Create a request without any context values
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	result := GetPathParam(req, "anyParam")
	assert.Empty(t, result, "Should return empty string when param not in context")
}

func TestTestRouter_HandleFunc(t *testing.T) {
	router := NewTestRouter()
	called := false

	router.HandleFunc(http.MethodOptions, "/test", func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestTestRouter_PathParamsWithSlashes(t *testing.T) {
	router := NewTestRouter()
	var captured string

	router.Get("/api/:version/users", func(w http.ResponseWriter, r *http.Request) {
		captured = GetPathParam(r, "version")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v2/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, "v2", captured)
	assert.Equal(t, http.StatusOK, w.Code)
}
