package factories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jarcoal/httpmock"
)

// ErrorResponse creates a standard Graph API error response
func ErrorResponse(statusCode int, errorCode, errorMessage string) httpmock.Responder {
	return httpmock.NewJsonResponderOrPanic(statusCode, map[string]any{
		"error": map[string]any{
			"code":    errorCode,
			"message": errorMessage,
		},
	})
}

// SuccessResponse creates a standard success response with the given body
func SuccessResponse(statusCode int, body interface{}) httpmock.Responder {
	return httpmock.NewJsonResponderOrPanic(statusCode, body)
}

// EmptySuccessResponse creates a standard empty success response
func EmptySuccessResponse(statusCode int) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(statusCode, ""), nil
	}
}

// RequestValidatorResponder creates a responder that validates the request before returning a response
func RequestValidatorResponder(validator func(req *http.Request) error, successResponder httpmock.Responder) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		if err := validator(req); err != nil {
			return httpmock.NewStringResponse(400, fmt.Sprintf(`{"error":{"code":"BadRequest","message":"%s"}}`, err.Error())), nil
		}
		return successResponder(req)
	}
}

// JSONBodyValidator validates that the request body is valid JSON and matches the expected structure
func JSONBodyValidator(expectedKeys []string) func(req *http.Request) error {
	return func(req *http.Request) error {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return fmt.Errorf("invalid JSON body: %v", err)
		}

		for _, key := range expectedKeys {
			if _, ok := body[key]; !ok {
				return fmt.Errorf("missing required field: %s", key)
			}
		}

		return nil
	}
}

// ExtractIDFromURL extracts the resource ID from a URL path given a base path
func ExtractIDFromURL(urlPath, basePath string) string {
	// Remove base path prefix
	urlPath = strings.TrimPrefix(urlPath, basePath)

	// Split by "/" and get the last non-empty part
	urlParts := strings.Split(urlPath, "/")
	for i := len(urlParts) - 1; i >= 0; i-- {
		if urlParts[i] != "" {
			return urlParts[i]
		}
	}

	// Fallback to simple split method
	parts := strings.Split(urlPath, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}
