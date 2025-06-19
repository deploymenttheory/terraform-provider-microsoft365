package mocks

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
)

// ErrorMocks provides mock HTTP responses for graph API errors.
type ErrorMocks struct{}

// NewErrorMocks creates a new instance of ErrorMocks.
func NewErrorMocks() *ErrorMocks {
	return &ErrorMocks{}
}

// RegisterMocks registers all error-related mock responders.
func (m *ErrorMocks) RegisterMocks() {
	// Generic error responders
	m.registerErrorResponder(400, "BadRequest", "The request is invalid.")
	m.registerErrorResponder(401, "Unauthorized", "Authentication failed.")
	m.registerErrorResponder(403, "Forbidden", "Access is denied.")
	m.registerErrorResponder(404, "NotFound", "The resource was not found.")
	m.registerErrorResponder(429, "TooManyRequests", "Rate limit exceeded.")
	m.registerErrorResponder(500, "InternalServerError", "An internal server error occurred.")
	m.registerErrorResponder(503, "ServiceUnavailable", "The service is temporarily unavailable.")
}

// registerErrorResponder is a helper to register a generic graph error response.
func (m *ErrorMocks) registerErrorResponder(statusCode int, errorCode, errorMessage string) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/errors/%d", statusCode)
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(statusCode, map[string]interface{}{
				"error": map[string]interface{}{
					"code":    errorCode,
					"message": errorMessage,
				},
			})
		},
	)
}
