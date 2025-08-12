package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRetryableDeleteError(t *testing.T) {
	tests := []struct {
		name           string
		errorInfo      *GraphErrorInfo
		expectedResult bool
	}{
		// Test retryable status codes
		{
			name: "429 Too Many Requests - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 429,
			},
			expectedResult: true,
		},
		{
			name: "500 Internal Server Error - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 500,
			},
			expectedResult: true,
		},
		{
			name: "502 Bad Gateway - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 502,
			},
			expectedResult: true,
		},
		{
			name: "503 Service Unavailable - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 503,
			},
			expectedResult: true,
		},
		{
			name: "504 Gateway Timeout - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 504,
			},
			expectedResult: true,
		},
		{
			name: "5001 Assignment Propagation - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 5001,
			},
			expectedResult: true,
		},
		// Test non-retryable status codes
		{
			name: "200 Success - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 200,
			},
			expectedResult: false,
		},
		{
			name: "400 Bad Request - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 400,
			},
			expectedResult: false,
		},
		{
			name: "401 Unauthorized - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 401,
			},
			expectedResult: false,
		},
		{
			name: "403 Forbidden - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 403,
			},
			expectedResult: false,
		},
		{
			name: "404 Not Found - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 404,
			},
			expectedResult: false,
		},
		// Test retryable error codes
		{
			name: "ServiceUnavailable error code - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999, // Unknown status code
				ErrorCode:  "ServiceUnavailable",
			},
			expectedResult: true,
		},
		{
			name: "RequestThrottled error code - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "RequestThrottled",
			},
			expectedResult: true,
		},
		{
			name: "RequestTimeout error code - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "RequestTimeout",
			},
			expectedResult: true,
		},
		{
			name: "InternalServerError error code - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "InternalServerError",
			},
			expectedResult: true,
		},
		{
			name: "BadGateway error code - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "BadGateway",
			},
			expectedResult: true,
		},
		{
			name: "GatewayTimeout error code - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "GatewayTimeout",
			},
			expectedResult: true,
		},
		// Test unknown errors
		{
			name: "Unknown status code and error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "UnknownError",
			},
			expectedResult: false,
		},
		{
			name:           "Empty error info - should not retry",
			errorInfo:      &GraphErrorInfo{},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryableDeleteError(tt.errorInfo)
			assert.Equal(t, tt.expectedResult, result, "IsRetryableDeleteError result should match expected")
		})
	}
}

func TestIsNonRetryableDeleteError(t *testing.T) {
	tests := []struct {
		name           string
		errorInfo      *GraphErrorInfo
		expectedResult bool
	}{
		// Test success cases
		{
			name: "200 OK - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 200,
			},
			expectedResult: true,
		},
		{
			name: "204 No Content - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 204,
			},
			expectedResult: true,
		},
		// Test client error cases
		{
			name: "400 Bad Request - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 400,
			},
			expectedResult: true,
		},
		{
			name: "401 Unauthorized - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 401,
			},
			expectedResult: true,
		},
		{
			name: "403 Forbidden - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 403,
			},
			expectedResult: true,
		},
		{
			name: "404 Not Found - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 404,
			},
			expectedResult: true,
		},
		{
			name: "409 Conflict - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 409,
			},
			expectedResult: true,
		},
		{
			name: "410 Gone - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 410,
			},
			expectedResult: true,
		},
		{
			name: "422 Unprocessable Entity - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 422,
			},
			expectedResult: true,
		},
		// Test retryable status codes (should return false)
		{
			name: "429 Too Many Requests - should retry (return false)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 429,
			},
			expectedResult: false,
		},
		{
			name: "500 Internal Server Error - should retry (return false)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 500,
			},
			expectedResult: false,
		},
		{
			name: "5001 Assignment Propagation - should retry (return false)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 5001,
			},
			expectedResult: false,
		},
		// Test non-retryable error codes
		{
			name: "BadRequest error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999, // Unknown status code
				ErrorCode:  "BadRequest",
			},
			expectedResult: true,
		},
		{
			name: "Unauthorized error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "Unauthorized",
			},
			expectedResult: true,
		},
		{
			name: "Forbidden error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "Forbidden",
			},
			expectedResult: true,
		},
		{
			name: "NotFound error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "NotFound",
			},
			expectedResult: true,
		},
		{
			name: "Conflict error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "Conflict",
			},
			expectedResult: true,
		},
		{
			name: "Gone error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "Gone",
			},
			expectedResult: true,
		},
		{
			name: "UnprocessableEntity error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "UnprocessableEntity",
			},
			expectedResult: true,
		},
		// Test unknown errors (should return false)
		{
			name: "Unknown status code and error code - should potentially retry (return false)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "UnknownError",
			},
			expectedResult: false,
		},
		{
			name:           "Empty error info - should potentially retry (return false)",
			errorInfo:      &GraphErrorInfo{},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNonRetryableDeleteError(tt.errorInfo)
			assert.Equal(t, tt.expectedResult, result, "IsNonRetryableDeleteError result should match expected")
		})
	}
}

func TestRetryableAndNonRetryableConsistency(t *testing.T) {
	// Test that the two functions are consistent with each other
	// An error should not be both retryable and non-retryable

	testCases := []GraphErrorInfo{
		{StatusCode: 200},
		{StatusCode: 204},
		{StatusCode: 400},
		{StatusCode: 401},
		{StatusCode: 403},
		{StatusCode: 404},
		{StatusCode: 409},
		{StatusCode: 410},
		{StatusCode: 422},
		{StatusCode: 429},
		{StatusCode: 500},
		{StatusCode: 502},
		{StatusCode: 503},
		{StatusCode: 504},
		{StatusCode: 5001},
		{StatusCode: 999, ErrorCode: "BadRequest"},
		{StatusCode: 999, ErrorCode: "ServiceUnavailable"},
		{StatusCode: 999, ErrorCode: "UnknownError"},
	}

	for _, errorInfo := range testCases {
		t.Run(fmt.Sprintf("StatusCode_%d_ErrorCode_%s", errorInfo.StatusCode, errorInfo.ErrorCode), func(t *testing.T) {
			isRetryable := IsRetryableDeleteError(&errorInfo)
			isNonRetryable := IsNonRetryableDeleteError(&errorInfo)

			// For definitive cases, exactly one should be true
			if isRetryable && isNonRetryable {
				t.Errorf("Error cannot be both retryable and non-retryable: StatusCode=%d, ErrorCode=%s",
					errorInfo.StatusCode, errorInfo.ErrorCode)
			}

			// Log the results for visibility
			t.Logf("StatusCode=%d, ErrorCode=%s: Retryable=%t, NonRetryable=%t",
				errorInfo.StatusCode, errorInfo.ErrorCode, isRetryable, isNonRetryable)
		})
	}
}

func TestIsRetryableReadError(t *testing.T) {
	tests := []struct {
		name           string
		errorInfo      *GraphErrorInfo
		expectedResult bool
	}{
		// Test retryable status codes for reads
		{
			name: "404 Not Found - should retry (propagation)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 404,
			},
			expectedResult: true,
		},
		{
			name: "409 Conflict - should retry (propagation)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 409,
			},
			expectedResult: true,
		},
		{
			name: "423 Locked - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 423,
			},
			expectedResult: true,
		},
		{
			name: "429 Too Many Requests - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 429,
			},
			expectedResult: true,
		},
		{
			name: "500 Internal Server Error - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 500,
			},
			expectedResult: true,
		},
		{
			name: "502 Bad Gateway - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 502,
			},
			expectedResult: true,
		},
		{
			name: "503 Service Unavailable - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 503,
			},
			expectedResult: true,
		},
		{
			name: "504 Gateway Timeout - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 504,
			},
			expectedResult: true,
		},
		// Test non-retryable status codes for reads
		{
			name: "200 Success - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 200,
			},
			expectedResult: false,
		},
		{
			name: "400 Bad Request - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 400,
			},
			expectedResult: false,
		},
		{
			name: "401 Unauthorized - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 401,
			},
			expectedResult: false,
		},
		{
			name: "403 Forbidden - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 403,
			},
			expectedResult: false,
		},
		{
			name: "405 Method Not Allowed - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 405,
			},
			expectedResult: false,
		},
		{
			name: "406 Not Acceptable - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 406,
			},
			expectedResult: false,
		},
		{
			name: "410 Gone - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 410,
			},
			expectedResult: false,
		},
		{
			name: "422 Unprocessable Entity - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 422,
			},
			expectedResult: false,
		},
		// Test retryable error codes for reads
		{
			name: "NotFound error code - should retry (propagation)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "NotFound",
			},
			expectedResult: true,
		},
		{
			name: "ResourceNotFound error code - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "ResourceNotFound",
			},
			expectedResult: true,
		},
		{
			name: "RequestThrottled error code - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "RequestThrottled",
			},
			expectedResult: true,
		},
		{
			name: "ServiceUnavailable error code - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "ServiceUnavailable",
			},
			expectedResult: true,
		},
		{
			name: "NetworkError error code - should retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "NetworkError",
			},
			expectedResult: true,
		},
		// Test unknown errors
		{
			name: "Unknown status code and error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "UnknownError",
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryableReadError(tt.errorInfo)
			assert.Equal(t, tt.expectedResult, result, "IsRetryableReadError result should match expected")
		})
	}
}

func TestIsNonRetryableReadError(t *testing.T) {
	tests := []struct {
		name           string
		errorInfo      *GraphErrorInfo
		expectedResult bool
	}{
		// Test success cases
		{
			name: "200 OK - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 200,
			},
			expectedResult: true,
		},
		{
			name: "204 No Content - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 204,
			},
			expectedResult: true,
		},
		// Test client error cases
		{
			name: "400 Bad Request - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 400,
			},
			expectedResult: true,
		},
		{
			name: "401 Unauthorized - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 401,
			},
			expectedResult: true,
		},
		{
			name: "403 Forbidden - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 403,
			},
			expectedResult: true,
		},
		{
			name: "405 Method Not Allowed - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 405,
			},
			expectedResult: true,
		},
		{
			name: "406 Not Acceptable - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 406,
			},
			expectedResult: true,
		},
		{
			name: "410 Gone - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 410,
			},
			expectedResult: true,
		},
		{
			name: "422 Unprocessable Entity - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 422,
			},
			expectedResult: true,
		},
		// Test retryable status codes (should return false)
		{
			name: "404 Not Found - should retry for reads (return false)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 404,
			},
			expectedResult: false,
		},
		{
			name: "409 Conflict - should retry for reads (return false)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 409,
			},
			expectedResult: false,
		},
		{
			name: "423 Locked - should retry (return false)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 423,
			},
			expectedResult: false,
		},
		{
			name: "429 Too Many Requests - should retry (return false)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 429,
			},
			expectedResult: false,
		},
		{
			name: "500 Internal Server Error - should retry (return false)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 500,
			},
			expectedResult: false,
		},
		// Test non-retryable error codes
		{
			name: "BadRequest error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "BadRequest",
			},
			expectedResult: true,
		},
		{
			name: "Unauthorized error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "Unauthorized",
			},
			expectedResult: true,
		},
		{
			name: "Forbidden error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "Forbidden",
			},
			expectedResult: true,
		},
		{
			name: "ValidationError error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "ValidationError",
			},
			expectedResult: true,
		},
		{
			name: "UnprocessableEntity error code - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "UnprocessableEntity",
			},
			expectedResult: true,
		},
		// Test unknown errors (should return false)
		{
			name: "Unknown status code and error code - should potentially retry (return false)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 999,
				ErrorCode:  "UnknownError",
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNonRetryableReadError(tt.errorInfo)
			assert.Equal(t, tt.expectedResult, result, "IsNonRetryableReadError result should match expected")
		})
	}
}

func TestReadRetryConsistency(t *testing.T) {
	// Test that the read retry functions are consistent with each other
	// An error should not be both retryable and non-retryable

	testCases := []GraphErrorInfo{
		{StatusCode: 200},
		{StatusCode: 204},
		{StatusCode: 400},
		{StatusCode: 401},
		{StatusCode: 403},
		{StatusCode: 404}, // Special case: retryable for reads
		{StatusCode: 405},
		{StatusCode: 406},
		{StatusCode: 409}, // Special case: retryable for reads
		{StatusCode: 410},
		{StatusCode: 422},
		{StatusCode: 423}, // Special case: retryable for reads
		{StatusCode: 429},
		{StatusCode: 500},
		{StatusCode: 502},
		{StatusCode: 503},
		{StatusCode: 504},
		{StatusCode: 999, ErrorCode: "BadRequest"},
		{StatusCode: 999, ErrorCode: "NotFound"}, // Special case: retryable for reads
		{StatusCode: 999, ErrorCode: "ServiceUnavailable"},
		{StatusCode: 999, ErrorCode: "UnknownError"},
	}

	for _, errorInfo := range testCases {
		t.Run(fmt.Sprintf("StatusCode_%d_ErrorCode_%s", errorInfo.StatusCode, errorInfo.ErrorCode), func(t *testing.T) {
			isRetryable := IsRetryableReadError(&errorInfo)
			isNonRetryable := IsNonRetryableReadError(&errorInfo)

			// For definitive cases, exactly one should be true
			if isRetryable && isNonRetryable {
				t.Errorf("Error cannot be both retryable and non-retryable for reads: StatusCode=%d, ErrorCode=%s",
					errorInfo.StatusCode, errorInfo.ErrorCode)
			}

			// Log the results for visibility
			t.Logf("StatusCode=%d, ErrorCode=%s: Retryable=%t, NonRetryable=%t",
				errorInfo.StatusCode, errorInfo.ErrorCode, isRetryable, isNonRetryable)
		})
	}
}

// TestEdgeCases tests edge cases and boundary conditions
func TestEdgeCases(t *testing.T) {
	t.Run("Nil GraphErrorInfo", func(t *testing.T) {
		// This should not panic
		assert.False(t, IsRetryableDeleteError(nil))
		assert.False(t, IsNonRetryableDeleteError(nil))
		assert.False(t, IsRetryableReadError(nil))
		assert.False(t, IsNonRetryableReadError(nil))
	})

	t.Run("Empty strings in ErrorCode", func(t *testing.T) {
		errorInfo := &GraphErrorInfo{
			StatusCode: 999,
			ErrorCode:  "",
		}

		// Should fall back to status code logic
		assert.False(t, IsRetryableDeleteError(errorInfo))
		assert.False(t, IsNonRetryableDeleteError(errorInfo))
	})

	t.Run("Case sensitivity in ErrorCode", func(t *testing.T) {
		// Test that error codes are case sensitive
		errorInfo := &GraphErrorInfo{
			StatusCode: 999,
			ErrorCode:  "serviceunavailable", // lowercase
		}

		// Should not match "ServiceUnavailable"
		assert.False(t, IsRetryableDeleteError(errorInfo))
	})
}
