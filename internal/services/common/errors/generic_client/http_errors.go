package generic_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// HTTPGraphError represents a comprehensive Graph API error extracted from HTTP responses
// This mirrors the errors.GraphErrorInfo structure but for raw HTTP calls
type HTTPGraphError struct {
	StatusCode      int
	ErrorCode       string
	ErrorMessage    string
	Target          string
	IsODataError    bool
	AdditionalData  map[string]any
	Headers         map[string][]string
	RequestDetails  string
	RetryAfter      string
	RequestID       string
	ClientRequestID string
	ErrorDate       string
	InnerErrors     []errors.InnerErrorInfo
	ErrorDetails    []errors.ErrorDetailInfo
	CorrelationID   string
	ThrottledReason string
	Category        errors.ErrorCategory
	DiagnosticInfo  string
	ResponseBody    string
}

// ExtractHTTPGraphError extracts and analyzes error information from HTTP Graph API responses
// This is an exact replica of errors.GraphError but for raw HTTP calls
func ExtractHTTPGraphError(ctx context.Context, httpResp *http.Response) *HTTPGraphError {
	errorInfo := &HTTPGraphError{
		AdditionalData: make(map[string]any),
		InnerErrors:    []errors.InnerErrorInfo{},
		ErrorDetails:   []errors.ErrorDetailInfo{},
		Headers:        make(map[string][]string),
	}

	if httpResp == nil {
		return errorInfo
	}

	errorInfo.StatusCode = httpResp.StatusCode
	errorInfo.ErrorMessage = httpResp.Status

	tflog.Debug(ctx, "Extracting HTTP error information", map[string]any{
		"status_code": httpResp.StatusCode,
		"status":      httpResp.Status,
	})

	// Extract headers with comprehensive information (mirrors extractHeaders in SDK)
	extractHTTPHeaders(httpResp, errorInfo)

	// Read response body for OData error parsing
	if httpResp.Body != nil {
		body, err := io.ReadAll(httpResp.Body)
		if err == nil {
			errorInfo.ResponseBody = string(body)
			// Try to parse as OData error format
			extractODataErrorFromHTTP(ctx, body, errorInfo)
		}
	}

	// Categorize the error (same logic as SDK)
	errorInfo.Category = categorizeHTTPError(errorInfo)

	logHTTPErrorDetails(ctx, errorInfo)
	return errorInfo
}

// extractHTTPHeaders extracts comprehensive header information from HTTP response
// This mirrors extractHeaders from the SDK error handling
func extractHTTPHeaders(httpResp *http.Response, errorInfo *HTTPGraphError) {
	for key, values := range httpResp.Header {
		errorInfo.Headers[key] = values
		if len(values) > 0 {
			// Extract specific headers that are useful for debugging (same as SDK)
			switch strings.ToLower(key) {
			case "request-id":
				errorInfo.RequestID = values[0]
			case "client-request-id":
				errorInfo.ClientRequestID = values[0]
			case "ms-correlation-id":
				errorInfo.CorrelationID = values[0]
			case "retry-after":
				errorInfo.RetryAfter = values[0]
			case "x-throttled-reason":
				errorInfo.ThrottledReason = values[0]
			case "date":
				errorInfo.ErrorDate = values[0]
			}
			errorInfo.RequestDetails += fmt.Sprintf("%s: %v\n", key, values)
		}
	}
}

// extractODataErrorFromHTTP handles OData specific errors from HTTP response body
// This mirrors extractODataError from the SDK
func extractODataErrorFromHTTP(ctx context.Context, responseBody []byte, errorInfo *HTTPGraphError) {
	var odataResponse struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Target  string `json:"target"`
			Details []struct {
				Code    string `json:"code"`
				Message string `json:"message"`
				Target  string `json:"target"`
			} `json:"details"`
			InnerError struct {
				RequestID       string    `json:"request-id"`
				ClientRequestID string    `json:"client-request-id"`
				Date            time.Time `json:"date"`
				ODataType       string    `json:"@odata.type"`
				Code            string    `json:"code"`
				Message         string    `json:"message"`
			} `json:"innerError"`
		} `json:"error"`
	}

	if json.Unmarshal(responseBody, &odataResponse) == nil && odataResponse.Error.Code != "" {
		errorInfo.IsODataError = true

		// Extract main error information (mirrors extractMainError)
		errorInfo.ErrorCode = odataResponse.Error.Code
		errorInfo.ErrorMessage = odataResponse.Error.Message
		errorInfo.Target = odataResponse.Error.Target

		tflog.Debug(ctx, "Found OData error", map[string]any{
			"code":    errorInfo.ErrorCode,
			"message": errorInfo.ErrorMessage,
			"target":  errorInfo.Target,
		})

		// Add description for known error codes (same as SDK)
		if description, exists := getODataErrorDescription(errorInfo.ErrorCode); exists {
			errorInfo.AdditionalData["error_description"] = description
		}

		// Process error details array (mirrors extractErrorDetails)
		for i, detail := range odataResponse.Error.Details {
			detailInfo := errors.ErrorDetailInfo{
				Code:    detail.Code,
				Message: detail.Message,
				Target:  detail.Target,
			}
			errorInfo.ErrorDetails = append(errorInfo.ErrorDetails, detailInfo)

			tflog.Debug(ctx, "Extracted error detail", map[string]any{
				"index":  i,
				"code":   detailInfo.Code,
				"target": detailInfo.Target,
			})
		}

		// Process inner error (mirrors extractInnerError)
		if odataResponse.Error.InnerError.RequestID != "" || odataResponse.Error.InnerError.ODataType != "" {
			innerInfo := errors.InnerErrorInfo{
				RequestID:   odataResponse.Error.InnerError.RequestID,
				ClientReqID: odataResponse.Error.InnerError.ClientRequestID,
				Date:        odataResponse.Error.InnerError.Date.String(),
				ODataType:   odataResponse.Error.InnerError.ODataType,
				Code:        odataResponse.Error.InnerError.Code,
				Message:     odataResponse.Error.InnerError.Message,
			}

			// Update main error info with inner error data (same logic as SDK)
			if errorInfo.RequestID == "" && innerInfo.RequestID != "" {
				errorInfo.RequestID = innerInfo.RequestID
			}
			if errorInfo.ClientRequestID == "" && innerInfo.ClientReqID != "" {
				errorInfo.ClientRequestID = innerInfo.ClientReqID
			}
			if errorInfo.ErrorDate == "" && innerInfo.Date != "" {
				errorInfo.ErrorDate = innerInfo.Date
			}
			if errorInfo.ErrorCode == "" && innerInfo.ODataType != "" {
				errorInfo.ErrorCode = innerInfo.ODataType
			}

			errorInfo.InnerErrors = append(errorInfo.InnerErrors, innerInfo)

			tflog.Debug(ctx, "Extracted inner error", map[string]any{
				"odata_type": innerInfo.ODataType,
				"request_id": innerInfo.RequestID,
			})
		}
	}
}

// getODataErrorDescription returns description for known OData error codes
func getODataErrorDescription(errorCode string) (string, bool) {
	// This would reference the comprehensive error codes from errors package
	descriptions := map[string]string{
		"InvalidAuthenticationToken": "The access token is invalid or expired.",
		"Forbidden":                  "The caller does not have permission to perform the operation.",
		"InvalidRequest":             "The request is invalid.",
		"ResourceNotFound":           "The specified resource does not exist.",
		"Unauthorized":               "Authentication is required to access this resource.",
		"BadRequest":                 "The request could not be understood by the server due to malformed syntax.",
		"ItemNotFound":               "The requested item was not found.",
		"TooManyRequests":            "The user has sent too many requests in a given amount of time.",
		"ServiceUnavailable":         "The server is currently unable to handle the request.",
		// Add more as needed from errors/odata_error_codes.go
	}
	desc, exists := descriptions[errorCode]
	return desc, exists
}

// categorizeHTTPError categorizes errors based on status code and error information
// This mirrors categorizeError from the SDK
func categorizeHTTPError(errorInfo *HTTPGraphError) errors.ErrorCategory {
	switch errorInfo.StatusCode {
	case 401:
		return errors.CategoryAuthentication
	case 403:
		return errors.CategoryAuthorization
	case 400, 422:
		return errors.CategoryValidation
	case 429:
		return errors.CategoryThrottling
	case 503, 502, 504, 500:
		return errors.CategoryService
	case 0: // Network errors typically have status code 0
		return errors.CategoryNetwork
	default:
		// Check error codes for more specific categorization (same as SDK)
		if errorInfo.ErrorCode != "" {
			switch {
			case strings.Contains(strings.ToLower(errorInfo.ErrorCode), "auth"):
				return errors.CategoryAuthentication
			case strings.Contains(strings.ToLower(errorInfo.ErrorCode), "forbidden"):
				return errors.CategoryAuthorization
			case strings.Contains(strings.ToLower(errorInfo.ErrorCode), "throttle"):
				return errors.CategoryThrottling
			case strings.Contains(strings.ToLower(errorInfo.ErrorCode), "network"):
				return errors.CategoryNetwork
			}
		}
		return errors.CategoryService
	}
}

// HandleHTTPGraphError processes HTTP Graph API errors - exact replica of HandleKiotaGraphError
func HandleHTTPGraphError(ctx context.Context, httpResp *http.Response, resp interface{}, operation string, requiredPermissions []string) {
	errorInfo := ExtractHTTPGraphError(ctx, httpResp)
	errorDesc := getHTTPErrorDescription(errorInfo.StatusCode)

	tflog.Debug(ctx, "Handling HTTP Graph error:", map[string]any{
		"status_code":    errorInfo.StatusCode,
		"operation":      operation,
		"error_code":     errorInfo.ErrorCode,
		"error_message":  errorInfo.ErrorMessage,
		"target":         errorInfo.Target,
		"is_odata_error": errorInfo.IsODataError,
		"request_id":     errorInfo.RequestID,
		"correlation_id": errorInfo.CorrelationID,
		"inner_errors":   len(errorInfo.InnerErrors),
		"error_details":  len(errorInfo.ErrorDetails),
		"category":       errorInfo.Category,
	})

	// Record error metrics for observability (same as SDK)
	recordHTTPErrorMetrics(ctx, errorInfo, operation)

	// Handle special cases first (exact same logic as HandleKiotaGraphError)
	switch errorInfo.StatusCode {
	case 400:
		if operation == "Read" {
			tflog.Warn(ctx, "Resource appears to no longer exist (400 Response), removing from state")
			removeHTTPResourceFromState(ctx, resp)
			return
		}
		addHTTPErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructHTTPDetailedErrorMessage(errorDesc.Detail, errorInfo))

	case 401, 403:
		tflog.Warn(ctx, fmt.Sprintf("Permission error on %s operation, check required Graph permissions", operation))
		handleHTTPPermissionError(ctx, *errorInfo, resp, operation, requiredPermissions)
		return

	case 404:
		if operation == "Read" {
			tflog.Warn(ctx, "Resource not found (404 Response), removing from state")
			removeHTTPResourceFromState(ctx, resp)
			return
		}
		addHTTPErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructHTTPDetailedErrorMessage(errorDesc.Detail, errorInfo))

	case 429:
		if operation == "Read" {
			tflog.Warn(ctx, "Rate limit exceeded on read operation")
			handleHTTPRateLimitError(ctx, *errorInfo, resp)
			return
		}
		addHTTPErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructHTTPDetailedErrorMessage(errorDesc.Detail, errorInfo))

	case 503:
		if operation == "Read" {
			tflog.Warn(ctx, "Service Unavailable (503 Response), service is temporarily unavailable")
			handleHTTPServiceUnavailableError(ctx, *errorInfo, resp)
			return
		}
		addHTTPErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructHTTPDetailedErrorMessage(errorDesc.Detail, errorInfo))

	default:
		addHTTPErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructHTTPDetailedErrorMessage(errorDesc.Detail, errorInfo))
	}
}

// toGraphErrorInfo converts HTTPGraphError to errors.GraphErrorInfo for reuse of existing retry logic
func toGraphErrorInfo(httpErr *HTTPGraphError) *errors.GraphErrorInfo {
	return &errors.GraphErrorInfo{
		StatusCode:      httpErr.StatusCode,
		ErrorCode:       httpErr.ErrorCode,
		ErrorMessage:    httpErr.ErrorMessage,
		Target:          httpErr.Target,
		IsODataError:    httpErr.IsODataError,
		AdditionalData:  httpErr.AdditionalData,
		RequestDetails:  httpErr.RequestDetails,
		RetryAfter:      httpErr.RetryAfter,
		RequestID:       httpErr.RequestID,
		ClientRequestID: httpErr.ClientRequestID,
		ErrorDate:       httpErr.ErrorDate,
		InnerErrors:     httpErr.InnerErrors,
		ErrorDetails:    httpErr.ErrorDetails,
		CorrelationID:   httpErr.CorrelationID,
		ThrottledReason: httpErr.ThrottledReason,
		Category:        httpErr.Category,
		DiagnosticInfo:  httpErr.DiagnosticInfo,
	}
}

// constructHTTPDetailedErrorMessage creates a comprehensive error message - mirrors SDK logic
func constructHTTPDetailedErrorMessage(standardDetail string, errorInfo *HTTPGraphError) string {
	var parts []string
	parts = append(parts, standardDetail)

	if errorInfo.ErrorMessage != "" {
		parts = append(parts, fmt.Sprintf("Error: %s", errorInfo.ErrorMessage))
	}

	if errorInfo.ErrorCode != "" {
		parts = append(parts, fmt.Sprintf("Code: %s", errorInfo.ErrorCode))
		if description, exists := getODataErrorDescription(errorInfo.ErrorCode); exists {
			parts = append(parts, fmt.Sprintf("Description: %s", description))
		}
	}

	if errorInfo.Target != "" {
		parts = append(parts, fmt.Sprintf("Target: %s", errorInfo.Target))
	}

	// Add error details if present (same logic as SDK)
	if len(errorInfo.ErrorDetails) > 0 {
		var detailParts []string
		for _, detail := range errorInfo.ErrorDetails {
			detailStr := ""
			if detail.Code != "" {
				detailStr += fmt.Sprintf("Code: %s", detail.Code)
			}
			if detail.Message != "" {
				if detailStr != "" {
					detailStr += " - "
				}
				detailStr += detail.Message
			}
			if detail.Target != "" {
				if detailStr != "" {
					detailStr += " "
				}
				detailStr += fmt.Sprintf("(Target: %s)", detail.Target)
			}
			if detailStr != "" {
				detailParts = append(detailParts, detailStr)
			}
		}
		if len(detailParts) > 0 {
			parts = append(parts, fmt.Sprintf("Details: %s", strings.Join(detailParts, "; ")))
		}
	}

	// Add inner error information (same logic as SDK)
	if len(errorInfo.InnerErrors) > 0 {
		var innerParts []string
		for i, inner := range errorInfo.InnerErrors {
			innerStr := fmt.Sprintf("Level %d", i+1)
			if inner.ODataType != "" {
				innerStr += fmt.Sprintf(" - Type: %s", inner.ODataType)
			}
			if inner.Code != "" {
				innerStr += fmt.Sprintf(" - Code: %s", inner.Code)
			}
			if inner.Message != "" {
				innerStr += fmt.Sprintf(" - Message: %s", inner.Message)
			}
			innerParts = append(innerParts, innerStr)
		}
		if len(innerParts) > 0 {
			parts = append(parts, fmt.Sprintf("Inner Errors: %s", strings.Join(innerParts, "; ")))
		}
	}

	// Add tracking information (same logic as SDK)
	var trackingParts []string
	if errorInfo.RequestID != "" {
		trackingParts = append(trackingParts, fmt.Sprintf("Request ID: %s", errorInfo.RequestID))
	}
	if errorInfo.ClientRequestID != "" {
		trackingParts = append(trackingParts, fmt.Sprintf("Client Request ID: %s", errorInfo.ClientRequestID))
	}
	if errorInfo.CorrelationID != "" {
		trackingParts = append(trackingParts, fmt.Sprintf("Correlation ID: %s", errorInfo.CorrelationID))
	}
	if errorInfo.ErrorDate != "" {
		trackingParts = append(trackingParts, fmt.Sprintf("Date: %s", errorInfo.ErrorDate))
	}
	if len(trackingParts) > 0 {
		parts = append(parts, fmt.Sprintf("Tracking: %s", strings.Join(trackingParts, ", ")))
	}

	// Add category information (same as SDK)
	if errorInfo.Category != "" {
		parts = append(parts, fmt.Sprintf("Category: %s", errorInfo.Category))
	}

	return strings.Join(parts, "\n")
}

// recordHTTPErrorMetrics records error metrics for observability - mirrors SDK
func recordHTTPErrorMetrics(ctx context.Context, errorInfo *HTTPGraphError, operation string) {
	tflog.Info(ctx, "HTTP Graph API Error Metrics", map[string]any{
		"metric_type":       "http_graph_api_error",
		"status_code":       errorInfo.StatusCode,
		"error_code":        errorInfo.ErrorCode,
		"operation":         operation,
		"category":          errorInfo.Category,
		"is_odata_error":    errorInfo.IsODataError,
		"has_inner_errors":  len(errorInfo.InnerErrors) > 0,
		"has_error_details": len(errorInfo.ErrorDetails) > 0,
		"request_id":        errorInfo.RequestID,
		"timestamp":         time.Now().Unix(),
	})
}

// handleHTTPPermissionError processes permission-related errors - mirrors SDK
func handleHTTPPermissionError(ctx context.Context, errorInfo HTTPGraphError, resp interface{}, operation string, requiredPermissions []string) {
	var permissionMsg string

	if len(requiredPermissions) == 1 {
		permissionMsg = fmt.Sprintf("%s operation requires permission: %s", operation, requiredPermissions[0])
	} else if len(requiredPermissions) > 1 {
		permissionMsg = fmt.Sprintf("%s operation requires one of the following permissions: %s", operation, strings.Join(requiredPermissions, ", "))
	} else {
		permissionMsg = fmt.Sprintf("%s operation: No specific permissions defined. Please check Microsoft documentation.", operation)
	}

	errorDesc := getHTTPErrorDescription(errorInfo.StatusCode)
	detail := fmt.Sprintf("%s\n%s\n%s",
		errorDesc.Detail,
		permissionMsg,
		constructHTTPDetailedErrorMessage("Graph API Error Details:", &errorInfo))

	addHTTPErrorToDiagnostics(ctx, resp, errorDesc.Summary, detail)
}

// handleHTTPRateLimitError processes rate limit errors - mirrors SDK
func handleHTTPRateLimitError(ctx context.Context, errorInfo HTTPGraphError, resp interface{}) {
	tflog.Warn(ctx, "Rate limit exceeded", map[string]any{
		"retry_after":      errorInfo.RetryAfter,
		"throttled_reason": errorInfo.ThrottledReason,
		"request_id":       errorInfo.RequestID,
		"details":          errorInfo.ErrorMessage,
	})

	errorDesc := getHTTPErrorDescription(429)
	detail := constructHTTPDetailedErrorMessage(errorDesc.Detail, &errorInfo)

	if errorInfo.RetryAfter != "" {
		detail += fmt.Sprintf("\nRetry-After: %s seconds", errorInfo.RetryAfter)
	}

	if errorInfo.ThrottledReason != "" {
		detail += fmt.Sprintf("\nThrottled Reason: %s", errorInfo.ThrottledReason)
	}

	addHTTPErrorToDiagnostics(ctx, resp, errorDesc.Summary, detail)
}

// handleHTTPServiceUnavailableError processes 503 Service Unavailable errors - mirrors SDK
func handleHTTPServiceUnavailableError(ctx context.Context, errorInfo HTTPGraphError, resp interface{}) {
	retryAfter := errorInfo.RetryAfter
	if retryAfter == "" {
		retryAfter = "unspecified"
	}

	tflog.Warn(ctx, "Service temporarily unavailable", map[string]any{
		"retry_after": retryAfter,
		"request_id":  errorInfo.RequestID,
		"details":     errorInfo.ErrorMessage,
	})

	errorDesc := getHTTPErrorDescription(503)
	detail := fmt.Sprintf(
		"%s\nThe service may be experiencing high load or undergoing maintenance.\nRetry-After: %s\n%s",
		errorDesc.Detail,
		retryAfter,
		constructHTTPDetailedErrorMessage("Additional Details:", &errorInfo),
	)

	addHTTPErrorToDiagnostics(ctx, resp, errorDesc.Summary, detail)
}

// getHTTPErrorDescription returns a standard error description - reuses SDK error descriptions
func getHTTPErrorDescription(statusCode int) errors.ErrorDescription {
	descriptions := map[int]errors.ErrorDescription{
		400: {Summary: "Bad Request - 400", Detail: "The request was invalid or malformed. Please check the request parameters and try again."},
		401: {Summary: "Unauthorized - 401", Detail: "Authentication failed. Please check your Entra ID credentials and permissions."},
		403: {Summary: "Forbidden - 403", Detail: "Your credentials lack sufficient authorisation to perform this operation. Grant the required Microsoft Graph permissions to your Entra ID authentication method."},
		404: {Summary: "Not Found - 404", Detail: "The requested resource was not found."},
		429: {Summary: "Too Many Requests - 429", Detail: "Request throttled by Microsoft Graph API rate limits. Please try again later."},
		500: {Summary: "Internal Server Error - 500", Detail: "Microsoft Graph API encountered an internal error. Please try again later."},
		503: {Summary: "Service Unavailable - 503", Detail: "The Microsoft Graph API service is temporarily unavailable or overloaded. This is typically a transient condition that will be automatically resolved after a short time."},
	}

	if desc, ok := descriptions[statusCode]; ok {
		return desc
	}
	return errors.ErrorDescription{
		Summary: fmt.Sprintf("HTTP Error %d", statusCode),
		Detail:  "An unexpected error occurred. Please check the request and try again.",
	}
}

// addHTTPErrorToDiagnostics adds an error to the response diagnostics - mirrors SDK
func addHTTPErrorToDiagnostics(ctx context.Context, resp interface{}, summary, detail string) {
	switch r := resp.(type) {
	case *resource.CreateResponse:
		r.Diagnostics.AddError(summary, detail)
	case *resource.ReadResponse:
		r.Diagnostics.AddError(summary, detail)
	case *resource.UpdateResponse:
		r.Diagnostics.AddError(summary, detail)
	case *resource.DeleteResponse:
		r.Diagnostics.AddError(summary, detail)
	case *datasource.ReadResponse:
		r.Diagnostics.AddError(summary, detail)
	default:
		tflog.Error(ctx, "Unknown response type in addHTTPErrorToDiagnostics", map[string]any{
			"response_type": fmt.Sprintf("%T", resp),
			"summary":       summary,
			"detail":        detail,
		})
	}
}

// removeHTTPResourceFromState removes a resource from the state - mirrors SDK
func removeHTTPResourceFromState(ctx context.Context, resp interface{}) {
	switch r := resp.(type) {
	case *resource.ReadResponse:
		r.State.RemoveResource(ctx)
		tflog.Debug(ctx, "Resource removed from state due to not found error")
	default:
		tflog.Error(ctx, "Cannot remove resource from state for this response type", map[string]any{
			"response_type": fmt.Sprintf("%T", resp),
		})
	}
}

// logHTTPErrorDetails logs comprehensive error details for debugging - mirrors SDK
func logHTTPErrorDetails(ctx context.Context, errorInfo *HTTPGraphError) {
	details := map[string]any{
		"status_code":       errorInfo.StatusCode,
		"error_code":        errorInfo.ErrorCode,
		"is_odata_error":    errorInfo.IsODataError,
		"error_message":     errorInfo.ErrorMessage,
		"target":            errorInfo.Target,
		"category":          errorInfo.Category,
		"request_id":        errorInfo.RequestID,
		"client_request_id": errorInfo.ClientRequestID,
		"correlation_id":    errorInfo.CorrelationID,
		"error_date":        errorInfo.ErrorDate,
		"retry_after":       errorInfo.RetryAfter,
		"throttled_reason":  errorInfo.ThrottledReason,
	}

	if len(errorInfo.AdditionalData) > 0 {
		details["additional_data"] = errorInfo.AdditionalData
	}

	if len(errorInfo.InnerErrors) > 0 {
		details["inner_errors_count"] = len(errorInfo.InnerErrors)
		details["inner_errors"] = errorInfo.InnerErrors
	}

	if len(errorInfo.ErrorDetails) > 0 {
		details["error_details_count"] = len(errorInfo.ErrorDetails)
		details["error_details"] = errorInfo.ErrorDetails
	}

	if errorInfo.RequestDetails != "" {
		details["request_headers"] = strings.Split(errorInfo.RequestDetails, "\n")
	}

	tflog.Debug(ctx, "Comprehensive HTTP error details", details)
}

// GetHTTPRetryDelay calculates the appropriate retry delay - mirrors SDK logic
func GetHTTPRetryDelay(errorInfo *HTTPGraphError, attempt int) time.Duration {
	graphErrorInfo := toGraphErrorInfo(errorInfo)
	return errors.GetRetryDelay(graphErrorInfo, attempt)
}
