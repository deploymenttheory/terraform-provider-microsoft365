package errors

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ErrorDescription contains standardized error messages and summaries
type ErrorDescription struct {
	Summary string
	Detail  string
}

// GraphErrorInfo contains extracted information from a Graph API error
type GraphErrorInfo struct {
	StatusCode     int
	ErrorCode      string
	ErrorMessage   string
	IsODataError   bool
	AdditionalData map[string]interface{}
	Headers        *abstractions.ResponseHeaders
	RequestDetails string
	RetryAfter     string
}

// HandleGraphError processes Graph API errors and dispatches them appropriately
func HandleGraphError(ctx context.Context, err error, resp interface{}, operation string, requiredPermissions []string) {
	errorInfo := GraphError(ctx, err)
	errorDesc := getErrorDescription(errorInfo.StatusCode)

	tflog.Debug(ctx, "Handling Graph error:", map[string]interface{}{
		"status_code":    errorInfo.StatusCode,
		"operation":      operation,
		"error_code":     errorInfo.ErrorCode,
		"error_message":  errorInfo.ErrorMessage,
		"is_odata_error": errorInfo.IsODataError,
	})

	// Handle special cases first
	switch errorInfo.StatusCode {
	case 400:
		if operation == "Read" {
			tflog.Warn(ctx, "Resource appears to no longer exist (400 Response), removing from state")
			removeResourceFromState(ctx, resp)
			return
		}
		addErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructErrorDetail(errorDesc.Detail, errorInfo.ErrorMessage))

	case 401, 403:
		tflog.Warn(ctx, fmt.Sprintf("Permission error on %s operation, check required Graph permissions", operation))

		var requiredPermissionsToReport []string
		switch operation {
		case "Read":
			requiredPermissionsToReport = requiredPermissions
		case "Create", "Update", "Delete":
			requiredPermissionsToReport = requiredPermissions
		default:
			requiredPermissionsToReport = []string{}
		}

		handlePermissionError(ctx, errorInfo, resp, operation, requiredPermissionsToReport)
		return

	case 404:
		if operation == "Read" {
			tflog.Warn(ctx, "Resource not found (404 Response), removing from state")
			removeResourceFromState(ctx, resp)
			return
		}
		addErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructErrorDetail(errorDesc.Detail, errorInfo.ErrorMessage))

	case 429:
		if operation == "Read" {
			tflog.Warn(ctx, "Rate limit exceeded on read operation, retrying")
			handleRateLimitError(ctx, errorInfo, resp)
			return
		}
		addErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructErrorDetail(errorDesc.Detail, errorInfo.ErrorMessage))

	case 503:
		if operation == "Read" {
			tflog.Warn(ctx, "Service Unavailable (503 Response), service is temporarily unavailable")
			handleServiceUnavailableError(ctx, errorInfo, resp)
			return
		}
		addErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructErrorDetail(errorDesc.Detail, errorInfo.ErrorMessage))

	default:
		addErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructErrorDetail(errorDesc.Detail, errorInfo.ErrorMessage))
	}
}

// Utility functions

// GraphError extracts and processes error information from Graph API errors
func GraphError(ctx context.Context, err error) GraphErrorInfo {
	errorInfo := GraphErrorInfo{
		AdditionalData: make(map[string]interface{}),
	}

	if err == nil {
		return errorInfo
	}

	errorInfo.ErrorMessage = err.Error()

	tflog.Debug(ctx, "Processing error", map[string]interface{}{
		"error_type": fmt.Sprintf("%T", err),
		"error":      err.Error(),
	})

	switch typedErr := err.(type) {
	case *url.Error:
		processURLError(ctx, typedErr, &errorInfo)
	case abstractions.ApiErrorable:
		processAPIError(ctx, typedErr, &errorInfo)
	default:
		// For unknown error types, set a sensible default
		errorInfo.StatusCode = 500 // Internal Server Error
		errorInfo.ErrorCode = "UnknownError"
	}

	logErrorDetails(ctx, &errorInfo)
	return errorInfo
}

// processURLError handles URL specific errors
func processURLError(ctx context.Context, urlErr *url.Error, errorInfo *GraphErrorInfo) {
	tflog.Debug(ctx, "URL error detected", map[string]interface{}{
		"url": urlErr.URL,
		"op":  urlErr.Op,
		"err": urlErr.Err.Error(),
	})

	// Handle different URL error cases
	switch {
	case strings.Contains(urlErr.Error(), "context deadline exceeded"):
		errorInfo.StatusCode = 504 // Gateway Timeout
		errorInfo.ErrorCode = "RequestTimeout"
		// Use a general error message, details will be added later from standardErrorDescriptions
	case strings.Contains(urlErr.Error(), "connection refused"):
		errorInfo.StatusCode = 503 // Service Unavailable
		errorInfo.ErrorCode = "ConnectionRefused"
	case strings.Contains(urlErr.Error(), "no such host"):
		errorInfo.StatusCode = 503 // Service Unavailable
		errorInfo.ErrorCode = "HostNotFound"
	default:
		errorInfo.StatusCode = 400 // Bad Request
		errorInfo.ErrorCode = "URLError"
	}

	// Store the original error message for context
	errorInfo.AdditionalData["original_error"] = urlErr.Error()

	// Set a consistent error message that will be enhanced with details from standardErrorDescriptions
	// when added to diagnostics
	errorInfo.ErrorMessage = urlErr.Error()
}

// processAPIError handles Microsoft Graph API specific errors
func processAPIError(ctx context.Context, apiErr abstractions.ApiErrorable, errorInfo *GraphErrorInfo) {
	errorInfo.StatusCode = apiErr.GetStatusCode()
	errorInfo.Headers = apiErr.GetResponseHeaders()

	extractHeadersFromError(apiErr, errorInfo)

	switch typedApiErr := apiErr.(type) {
	case *odataerrors.ODataError:
		processODataError(ctx, typedApiErr, errorInfo)
	case *abstractions.ApiError:
		// For API errors, use the error message
		if typedApiErr.Message != "" {
			errorInfo.ErrorMessage = typedApiErr.Message
		}
	}
}

// extractHeadersFromError extracts header information from the API error
func extractHeadersFromError(apiErr abstractions.ApiErrorable, errorInfo *GraphErrorInfo) {
	if headers := apiErr.GetResponseHeaders(); headers != nil {
		for _, key := range headers.ListKeys() {
			values := headers.Get(key)
			if len(values) > 0 {
				errorInfo.RequestDetails += fmt.Sprintf("%s: %v\n", key, values)
			}
		}
	}
}

// processODataError handles OData specific errors
func processODataError(ctx context.Context, odataErr *odataerrors.ODataError, errorInfo *GraphErrorInfo) {
	errorInfo.IsODataError = true

	// Get the main error object (inside the "error" property in JSON)
	if mainError := odataErr.GetErrorEscaped(); mainError != nil {
		// Extract information from the main error object
		extractErrorInfo(ctx, mainError, errorInfo, true)

		// Process inner errors recursively to find the most specific error code
		if innerError := mainError.GetInnerError(); innerError != nil {
			processInnerErrorRecursively(ctx, innerError, errorInfo)
		}
	}
}

// extractErrorInfo extracts common properties from an error object
func extractErrorInfo(ctx context.Context, errorObj interface{}, errorInfo *GraphErrorInfo, isTopLevel bool) {
	// Handle different error object types
	switch typedError := errorObj.(type) {
	case odataerrors.MainErrorable:
		// Extract code
		if code := typedError.GetCode(); code != nil && *code != "" {
			// Only overwrite the code if this is the top level, or if we don't have a code yet
			if isTopLevel || errorInfo.ErrorCode == "" {
				errorInfo.ErrorCode = *code
				tflog.Debug(ctx, "Found error code", map[string]interface{}{
					"code":  *code,
					"level": "main",
				})

				// Add description for known error codes
				if description, exists := commonODataErrorCodes[*code]; exists {
					errorInfo.AdditionalData["error_description"] = description
				}
			}
		}

		// Extract message
		if message := typedError.GetMessage(); message != nil && *message != "" {
			// For top-level errors, this becomes the primary error message
			if isTopLevel {
				errorInfo.ErrorMessage = *message
			} else {
				// For nested errors, add as additional context
				errorInfo.AdditionalData["inner_message"] = *message
			}

			tflog.Debug(ctx, "Found error message", map[string]interface{}{
				"message": *message,
				"level":   "main",
			})
		}

		// Process details array
		details := typedError.GetDetails()
		if len(details) > 0 {
			tflog.Debug(ctx, "Processing error details", map[string]interface{}{
				"detail_count": len(details),
			})

			var detailsInfo []map[string]string
			var detailMessages []string

			for i, detail := range details {
				detailInfo := make(map[string]string)

				if code := detail.GetCode(); code != nil && *code != "" {
					detailInfo["code"] = *code
				}

				if msg := detail.GetMessage(); msg != nil && *msg != "" {
					detailInfo["message"] = *msg
					detailMessages = append(detailMessages, *msg)
				}

				if target := detail.GetTarget(); target != nil && *target != "" {
					detailInfo["target"] = *target
				}

				if len(detailInfo) > 0 {
					detailsInfo = append(detailsInfo, detailInfo)
					tflog.Debug(ctx, "Processed detail", map[string]interface{}{
						"index": i,
						"info":  detailInfo,
					})
				}
			}

			if len(detailsInfo) > 0 {
				errorInfo.AdditionalData["details"] = detailsInfo
			}

			if len(detailMessages) > 0 {
				errorInfo.ErrorMessage += "\nDetails: " + strings.Join(detailMessages, "; ")
			}
		}
	}
}

// processInnerErrorRecursively processes inner errors to find the most specific error code
func processInnerErrorRecursively(ctx context.Context, innerError odataerrors.InnerErrorable, errorInfo *GraphErrorInfo) {
	if innerError == nil {
		return
	}

	tflog.Debug(ctx, "Processing inner error", map[string]interface{}{
		"has_inner_error": true,
	})

	// Extract request/client IDs and date
	if reqID := innerError.GetRequestId(); reqID != nil && *reqID != "" {
		errorInfo.AdditionalData["request_id"] = *reqID
	}

	if clientReqID := innerError.GetClientRequestId(); clientReqID != nil && *clientReqID != "" {
		errorInfo.AdditionalData["client_request_id"] = *clientReqID
	}

	if date := innerError.GetDate(); date != nil {
		errorInfo.AdditionalData["error_date"] = date.String()
	}

	// Extract more specific code if available
	if code := innerError.GetOdataType(); code != nil && *code != "" {
		// Inner error codes are often more specific, so they take precedence
		errorInfo.ErrorCode = *code
		tflog.Debug(ctx, "Found inner error code", map[string]interface{}{
			"code": *code,
		})

		// Add description for known error codes
		if description, exists := commonODataErrorCodes[*code]; exists {
			errorInfo.AdditionalData["inner_error_description"] = description
		}
	}

}

// getErrorDescription returns a standard error description based on the HTTP status code
func getErrorDescription(statusCode int) ErrorDescription {
	if desc, ok := standardErrorDescriptions[statusCode]; ok {
		return desc
	}
	return ErrorDescription{
		Summary: fmt.Sprintf("HTTP Error %d", statusCode),
		Detail:  "An unexpected error occurred. Please check the request and try again.",
	}
}

// constructErrorDetail combines standard and specific error messages
func constructErrorDetail(standardDetail, specificMessage string) string {
	if specificMessage != "" {
		return fmt.Sprintf("%s\nDetails: %s", standardDetail, specificMessage)
	}
	return standardDetail
}

// handlePermissionError processes permission-related errors
func handlePermissionError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}, operation string, requiredPermissions []string) {
	var permissionMsg string

	// Format the message based on number of permissions
	if len(requiredPermissions) == 1 {
		permissionMsg = fmt.Sprintf("%s operation requires permission: %s", operation, requiredPermissions[0])
	} else if len(requiredPermissions) > 1 {
		permissionMsg = fmt.Sprintf("%s operation requires one of the following permission options: %s", operation, strings.Join(requiredPermissions, ", "))
	} else {
		permissionMsg = fmt.Sprintf("%s operation: No specific permissions defined. Please check microsoft documentation.", operation)
	}

	errorDesc := getErrorDescription(errorInfo.StatusCode)
	detail := fmt.Sprintf("%s\n%s\nGraph API Error: %s",
		errorDesc.Detail,
		permissionMsg,
		errorInfo.ErrorMessage)

	addErrorToDiagnostics(ctx, resp, errorDesc.Summary, detail)
}

// handleRateLimitError processes rate limit errors and adds retry information to the error message
func handleRateLimitError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}) GraphErrorInfo {
	if headers := errorInfo.Headers; headers != nil {
		retryValues := headers.Get("Retry-After")
		if len(retryValues) > 0 {
			errorInfo.RetryAfter = retryValues[0]
		}
	}

	tflog.Warn(ctx, "Rate limit exceeded", map[string]interface{}{
		"retry_after": errorInfo.RetryAfter,
		"details":     errorInfo.ErrorMessage,
	})

	errorDesc := getErrorDescription(429)
	detail := constructErrorDetail(errorDesc.Detail, errorInfo.ErrorMessage)
	addErrorToDiagnostics(ctx, resp, errorDesc.Summary, detail)

	return errorInfo
}

// handleServiceUnavailableError processes 503 Service Unavailable errors
func handleServiceUnavailableError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}) {
	retryAfter := "unspecified"
	if headers := errorInfo.Headers; headers != nil {
		retryValues := headers.Get("Retry-After")
		if len(retryValues) > 0 {
			retryAfter = retryValues[0]
			errorInfo.RetryAfter = retryAfter
		}
	}

	tflog.Warn(ctx, "Service temporarily unavailable", map[string]interface{}{
		"retry_after": retryAfter,
		"details":     errorInfo.ErrorMessage,
	})

	errorDesc := getErrorDescription(503)
	detail := fmt.Sprintf(
		"%s\nThe service may be experiencing high load or undergoing maintenance.\nRetry-After: %s\nDetails: %s",
		errorDesc.Detail,
		retryAfter,
		errorInfo.ErrorMessage,
	)

	addErrorToDiagnostics(ctx, resp, errorDesc.Summary, detail)
}

// addErrorToDiagnostics adds an error to the response diagnostics
func addErrorToDiagnostics(ctx context.Context, resp interface{}, summary, detail string) {
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
		tflog.Error(ctx, "Unknown response type in addErrorToDiagnostics")
	}
}

// removeResourceFromState removes a resource from the state
func removeResourceFromState(ctx context.Context, resp interface{}) {
	switch r := resp.(type) {
	case *resource.ReadResponse:
		r.State.RemoveResource(ctx)
	default:
		tflog.Error(ctx, "Cannot remove resource from state for this response type")
	}
}

// logErrorDetails logs error details for debugging
func logErrorDetails(ctx context.Context, errorInfo *GraphErrorInfo) {
	details := map[string]interface{}{
		"status_code":    errorInfo.StatusCode,
		"is_odata_error": errorInfo.IsODataError,
		"error_message":  errorInfo.ErrorMessage,
	}

	if errorInfo.ErrorCode != "" {
		details["error_code"] = errorInfo.ErrorCode
	}

	if len(errorInfo.AdditionalData) > 0 {
		details["additional_data"] = errorInfo.AdditionalData
	}

	if errorInfo.RequestDetails != "" {
		details["request_details"] = errorInfo.RequestDetails
	}

	tflog.Debug(ctx, "Error details", details)
}
