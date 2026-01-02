package errors

// IsRetryableDeleteError determines if an error is retryable for delete operations based on the response http status code and error information
func IsRetryableDeleteError(errorInfo *GraphErrorInfo) bool {
	if errorInfo == nil {
		return false
	}

	switch errorInfo.StatusCode {
	case 429, 503, 502, 504: // Rate limiting and service unavailable errors
		return true
	case 500: // Internal server errors might be retryable
		return true
	case 5001: // Resource assignment propagation - resource is currently assigned, waiting for propagation
		return true
	default:
		// Check specific error codes that might be retryable
		retryableErrorCodes := map[string]bool{
			"ServiceUnavailable":  true,
			"RequestThrottled":    true,
			"RequestTimeout":      true,
			"InternalServerError": true,
			"BadGateway":          true,
			"GatewayTimeout":      true,
			"UnknownError":        true, // Generic errors often related to resource locking or propagation
		}
		return retryableErrorCodes[errorInfo.ErrorCode]
	}
}

// IsNonRetryableDeleteError determines if an error should never be retried for delete operations based on the response http status code and error information
func IsNonRetryableDeleteError(errorInfo *GraphErrorInfo) bool {
	if errorInfo == nil {
		return false
	}

	switch errorInfo.StatusCode {
	case 200, 204: // Success cases
		return true
	case 400, 401, 403, 404, 409, 410, 422: // Client errors that won't change on retry
		return true
	default:
		// Check specific error codes that are permanent failures
		nonRetryableErrorCodes := map[string]bool{
			"BadRequest":          true,
			"Unauthorized":        true,
			"Forbidden":           true,
			"NotFound":            true,
			"Conflict":            true,
			"Gone":                true,
			"UnprocessableEntity": true,
		}
		return nonRetryableErrorCodes[errorInfo.ErrorCode]
	}
}

// IsRetryableReadError determines if an error is retryable for read operations (especially after create/update) based on the response http status code and error information
func IsRetryableReadError(errorInfo *GraphErrorInfo) bool {
	if errorInfo == nil {
		return false
	}

	switch errorInfo.StatusCode {
	case 404, 409, 423, 429: // Rate limiting and service unavailable errors
		return true
	case 500, 503, 502, 504: // Internal server errors might be retryable
		return true
	default:
		// Check specific error codes that might be retryable
		retryableErrorCodes := map[string]bool{
			"ServiceUnavailable":  true,
			"RequestThrottled":    true,
			"RequestTimeout":      true,
			"InternalServerError": true,
			"BadGateway":          true,
			"GatewayTimeout":      true,
			"NotFound":            true, // Resource propagation
			"ResourceNotFound":    true, // Resource propagation
			"NetworkError":        true, // Network connectivity issues
		}
		return retryableErrorCodes[errorInfo.ErrorCode]
	}
}

// IsNonRetryableReadError determines if an error should never be retried for read operations based on the response http status code and error information
func IsNonRetryableReadError(errorInfo *GraphErrorInfo) bool {
	if errorInfo == nil {
		return false
	}

	switch errorInfo.StatusCode {
	case 200, 204: // Success cases
		return true
	case 400, 401, 403, 405, 406, 410, 422: // Client errors that won't change on retry
		return true
	// Note: 409 Conflict is retryable for reads (removed from here)
	default:
		// Check specific error codes that are permanent failures
		nonRetryableErrorCodes := map[string]bool{
			"BadRequest":          true,
			"Unauthorized":        true,
			"Forbidden":           true,
			"Conflict":            true,
			"Gone":                true,
			"UnprocessableEntity": true,
			"ValidationError":     true,
			// Note: "NotFound" is NOT here because it's retryable for reads (propagation)
		}
		return nonRetryableErrorCodes[errorInfo.ErrorCode]
	}
}
