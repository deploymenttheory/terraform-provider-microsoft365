package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRetryableWriteError(t *testing.T) {
	tests := []struct {
		name           string
		errorInfo      *GraphErrorInfo
		expectedResult bool
	}{
		// Test retryable status codes
		{
			name: "404 Not Found - should retry (referenced resource propagation)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 404,
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
		// Test retryable error codes
		{
			name: "NotFound error code - should retry (referenced resource propagation)",
			errorInfo: &GraphErrorInfo{
				ErrorCode: "NotFound",
			},
			expectedResult: true,
		},
		{
			name: "ResourceNotFound error code - should retry (referenced resource propagation)",
			errorInfo: &GraphErrorInfo{
				ErrorCode: "ResourceNotFound",
			},
			expectedResult: true,
		},
		{
			name: "RequestThrottled error code - should retry",
			errorInfo: &GraphErrorInfo{
				ErrorCode: "RequestThrottled",
			},
			expectedResult: true,
		},
		// Test non-retryable status codes
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
			name: "409 Conflict - should not retry",
			errorInfo: &GraphErrorInfo{
				StatusCode: 409,
			},
			expectedResult: false,
		},
		{
			name:           "nil error info - should not retry",
			errorInfo:      nil,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryableWriteError(tt.errorInfo)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestIsNonRetryableWriteError(t *testing.T) {
	tests := []struct {
		name           string
		errorInfo      *GraphErrorInfo
		expectedResult bool
	}{
		// Test non-retryable status codes
		{
			name: "200 Success - non-retryable",
			errorInfo: &GraphErrorInfo{
				StatusCode: 200,
			},
			expectedResult: true,
		},
		{
			name: "201 Created - non-retryable",
			errorInfo: &GraphErrorInfo{
				StatusCode: 201,
			},
			expectedResult: true,
		},
		{
			name: "400 Bad Request - non-retryable",
			errorInfo: &GraphErrorInfo{
				StatusCode: 400,
			},
			expectedResult: true,
		},
		{
			name: "401 Unauthorized - non-retryable",
			errorInfo: &GraphErrorInfo{
				StatusCode: 401,
			},
			expectedResult: true,
		},
		{
			name: "403 Forbidden - non-retryable",
			errorInfo: &GraphErrorInfo{
				StatusCode: 403,
			},
			expectedResult: true,
		},
		{
			name: "409 Conflict - non-retryable",
			errorInfo: &GraphErrorInfo{
				StatusCode: 409,
			},
			expectedResult: true,
		},
		{
			name: "422 Unprocessable Entity - non-retryable",
			errorInfo: &GraphErrorInfo{
				StatusCode: 422,
			},
			expectedResult: true,
		},
		// 404 must NOT be classified as non-retryable for writes
		{
			name: "404 Not Found - not classified as non-retryable (referenced resource propagation)",
			errorInfo: &GraphErrorInfo{
				StatusCode: 404,
			},
			expectedResult: false,
		},
		{
			name: "429 Too Many Requests - not classified as non-retryable",
			errorInfo: &GraphErrorInfo{
				StatusCode: 429,
			},
			expectedResult: false,
		},
		{
			name: "503 Service Unavailable - not classified as non-retryable",
			errorInfo: &GraphErrorInfo{
				StatusCode: 503,
			},
			expectedResult: false,
		},
		{
			name:           "nil error info - not classified as non-retryable",
			errorInfo:      nil,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNonRetryableWriteError(tt.errorInfo)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
