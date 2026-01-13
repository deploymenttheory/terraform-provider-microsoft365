package errors

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupHTTPMock initializes httpmock and registers common error responses
func setupHTTPMock() {
	httpmock.Activate()

	// Register standard error responses
	registerErrorResponse(400, "BadRequest", "The request is invalid.")
	registerErrorResponse(401, "Unauthorized", "Authentication failed.")
	registerErrorResponse(403, "Forbidden", "Access is denied.")
	registerErrorResponse(404, "NotFound", "The resource was not found.")
	registerErrorResponse(429, "TooManyRequests", "Rate limit exceeded.")
	registerErrorResponse(500, "InternalServerError", "An internal server error occurred.")
	registerErrorResponse(503, "ServiceUnavailable", "The service is temporarily unavailable.")
}

// registerErrorResponse is a helper to register a generic graph error response
func registerErrorResponse(statusCode int, errorCode, errorMessage string) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/errors/%d", statusCode)
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(statusCode, map[string]any{
				"error": map[string]any{
					"code":    errorCode,
					"message": errorMessage,
				},
			})
		},
	)
}

// teardownHTTPMock deactivates httpmock
func teardownHTTPMock() {
	httpmock.DeactivateAndReset()
}

// TestGraphError_URLError tests handling of URL errors
func TestGraphError_URLError(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   string
		expectedCat    ErrorCategory
	}{
		{
			name:           "Context deadline exceeded",
			err:            &url.Error{Op: "Get", URL: "https://graph.microsoft.com/v1.0/users", Err: fmt.Errorf("context deadline exceeded")},
			expectedStatus: 504,
			expectedCode:   "RequestTimeout",
			expectedCat:    CategoryService,
		},
		{
			name:           "Connection refused",
			err:            &url.Error{Op: "Get", URL: "https://graph.microsoft.com/v1.0/users", Err: fmt.Errorf("connection refused")},
			expectedStatus: 503,
			expectedCode:   "ConnectionRefused",
			expectedCat:    CategoryService,
		},
		{
			name:           "No such host",
			err:            &url.Error{Op: "Get", URL: "https://graph.microsoft.com/v1.0/users", Err: fmt.Errorf("no such host")},
			expectedStatus: 503,
			expectedCode:   "HostNotFound",
			expectedCat:    CategoryService,
		},
		{
			name:           "Network is unreachable",
			err:            &url.Error{Op: "Get", URL: "https://graph.microsoft.com/v1.0/users", Err: fmt.Errorf("network is unreachable")},
			expectedStatus: 503,
			expectedCode:   "NetworkUnreachable",
			expectedCat:    CategoryService,
		},
		{
			name:           "Certificate error",
			err:            &url.Error{Op: "Get", URL: "https://graph.microsoft.com/v1.0/users", Err: fmt.Errorf("certificate has expired")},
			expectedStatus: 503,
			expectedCode:   "CertificateError",
			expectedCat:    CategoryService,
		},
		{
			name:           "Generic URL error",
			err:            &url.Error{Op: "Get", URL: "https://graph.microsoft.com/v1.0/users", Err: fmt.Errorf("generic error")},
			expectedStatus: 400,
			expectedCode:   "URLError",
			expectedCat:    CategoryValidation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errorInfo := GraphError(ctx, tc.err)

			assert.Equal(t, tc.expectedStatus, errorInfo.StatusCode, "Status code should match")
			assert.Equal(t, tc.expectedCode, errorInfo.ErrorCode, "Error code should match")
			assert.Equal(t, tc.expectedCat, errorInfo.Category, "Category should match")
			assert.NotEmpty(t, errorInfo.ErrorMessage, "Error message should not be empty")
			assert.Contains(t, errorInfo.AdditionalData["url"], "graph.microsoft.com", "URL should be in additional data")
		})
	}
}

// TestCategorizeError tests the categorizeError function
func TestCategorizeError(t *testing.T) {
	testCases := []struct {
		name        string
		errorInfo   GraphErrorInfo
		expectedCat ErrorCategory
	}{
		{
			name: "401 Unauthorized",
			errorInfo: GraphErrorInfo{
				StatusCode: 401,
			},
			expectedCat: CategoryAuthentication,
		},
		{
			name: "403 Forbidden",
			errorInfo: GraphErrorInfo{
				StatusCode: 403,
			},
			expectedCat: CategoryAuthorization,
		},
		{
			name: "400 Bad Request",
			errorInfo: GraphErrorInfo{
				StatusCode: 400,
			},
			expectedCat: CategoryValidation,
		},
		{
			name: "429 Too Many Requests",
			errorInfo: GraphErrorInfo{
				StatusCode: 429,
			},
			expectedCat: CategoryThrottling,
		},
		{
			name: "503 Service Unavailable",
			errorInfo: GraphErrorInfo{
				StatusCode: 503,
			},
			expectedCat: CategoryService,
		},
		{
			name: "Network Error",
			errorInfo: GraphErrorInfo{
				StatusCode: 0,
			},
			expectedCat: CategoryNetwork,
		},
		{
			name: "Auth in error code",
			errorInfo: GraphErrorInfo{
				StatusCode: 0,
				ErrorCode:  "AuthenticationFailed",
			},
			expectedCat: CategoryNetwork,
		},
		{
			name: "Forbidden in error code",
			errorInfo: GraphErrorInfo{
				StatusCode: 0,
				ErrorCode:  "AccessForbidden",
			},
			expectedCat: CategoryNetwork,
		},
		{
			name: "Throttle in error code",
			errorInfo: GraphErrorInfo{
				StatusCode: 0,
				ErrorCode:  "RequestThrottled",
			},
			expectedCat: CategoryNetwork,
		},
		{
			name: "Network in error code",
			errorInfo: GraphErrorInfo{
				StatusCode: 0,
				ErrorCode:  "NetworkError",
			},
			expectedCat: CategoryNetwork,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			category := categorizeError(&tc.errorInfo)
			assert.Equal(t, tc.expectedCat, category, "Error category should match expectation")
		})
	}
}

// TestHTTPMockErrorResponses tests the error handling with httpmock
func TestHTTPMockErrorResponses(t *testing.T) {
	setupHTTPMock()
	defer teardownHTTPMock()

	testCases := []struct {
		name       string
		statusCode int
		errorCode  string
	}{
		{
			name:       "Bad Request",
			statusCode: 400,
			errorCode:  "BadRequest",
		},
		{
			name:       "Unauthorized",
			statusCode: 401,
			errorCode:  "Unauthorized",
		},
		{
			name:       "Not Found",
			statusCode: 404,
			errorCode:  "NotFound",
		},
		{
			name:       "Rate Limit",
			statusCode: 429,
			errorCode:  "TooManyRequests",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("https://graph.microsoft.com/v1.0/errors/%d", tc.statusCode)
			resp, err := http.Get(url)
			require.NoError(t, err, "HTTP request should succeed")
			defer resp.Body.Close()

			assert.Equal(t, tc.statusCode, resp.StatusCode, "Response status code should match")

			// Verify the call was made
			info := httpmock.GetCallCountInfo()
			count := info["GET "+url]
			assert.Equal(t, 1, count, "Expected 1 call to the URL")
		})
	}
}

// TestHandleKiotaGraphError_WithSimpleErrors tests HandleKiotaGraphError with simple error types
func TestHandleKiotaGraphError_WithSimpleErrors(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name              string
		statusCode        int
		operation         string
		expectDiagnostics bool
	}{
		{
			name:              "Bad Request Error",
			statusCode:        400,
			operation:         constants.TfOperationCreate,
			expectDiagnostics: true,
		},
		{
			name:              "Not Found Error",
			statusCode:        404,
			operation:         constants.TfOperationUpdate,
			expectDiagnostics: true,
		},
		{
			name:              "Internal Server Error",
			statusCode:        500,
			operation:         constants.TfOperationDelete,
			expectDiagnostics: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a response object
			resp := &resource.CreateResponse{
				Diagnostics: diag.Diagnostics{},
			}

			// Create a simple error
			err := fmt.Errorf("error with status code %d", tc.statusCode)

			// Process the error
			HandleKiotaGraphError(ctx, err, resp, tc.operation, []string{"User.Read"})

			// Check if diagnostics were added
			assert.Equal(t, tc.expectDiagnostics, resp.Diagnostics.HasError(), "Diagnostics error state should match expectation")
		})
	}
}

// TestAddErrorToDiagnostics tests the addErrorToDiagnostics function
func TestAddErrorToDiagnostics(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name     string
		respType string
		summary  string
		detail   string
	}{
		{
			name:     "Create Response",
			respType: constants.TfOperationCreate,
			summary:  "Test Summary",
			detail:   "Test Detail",
		},
		{
			name:     "Read Response",
			respType: constants.TfOperationRead,
			summary:  "Read Error",
			detail:   "Read Error Detail",
		},
		{
			name:     "Update Response",
			respType: constants.TfOperationUpdate,
			summary:  "Update Error",
			detail:   "Update Error Detail",
		},
		{
			name:     "Delete Response",
			respType: constants.TfOperationDelete,
			summary:  "Delete Error",
			detail:   "Delete Error Detail",
		},
		{
			name:     "DataSource Response",
			respType: "datasource",
			summary:  "DataSource Error",
			detail:   "DataSource Error Detail",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var resp any

			// Create the appropriate response type
			switch tc.respType {
			case constants.TfOperationCreate:
				resp = &resource.CreateResponse{
					Diagnostics: diag.Diagnostics{},
				}
			case constants.TfOperationRead:
				resp = &resource.ReadResponse{
					Diagnostics: diag.Diagnostics{},
				}
			case constants.TfOperationUpdate:
				resp = &resource.UpdateResponse{
					Diagnostics: diag.Diagnostics{},
				}
			case constants.TfOperationDelete:
				resp = &resource.DeleteResponse{
					Diagnostics: diag.Diagnostics{},
				}
			case "datasource":
				resp = &resource.CreateResponse{
					Diagnostics: diag.Diagnostics{},
				}
			}

			// Add error to diagnostics
			addErrorToDiagnostics(ctx, resp, tc.summary, tc.detail)

			// Check if diagnostics were added
			switch r := resp.(type) {
			case *resource.CreateResponse:
				assert.True(t, r.Diagnostics.HasError(), "Diagnostics should have error")
			case *resource.ReadResponse:
				assert.True(t, r.Diagnostics.HasError(), "Diagnostics should have error")
			case *resource.UpdateResponse:
				assert.True(t, r.Diagnostics.HasError(), "Diagnostics should have error")
			case *resource.DeleteResponse:
				assert.True(t, r.Diagnostics.HasError(), "Diagnostics should have error")
			}
		})
	}
}
