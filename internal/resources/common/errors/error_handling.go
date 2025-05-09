package errors

import (
	"context"
	"fmt"
	"strings"

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

// standardErrorDescriptions provides consistent error messaging across the provider
var standardErrorDescriptions = map[int]ErrorDescription{
	400: {
		Summary: "Bad Request - 400",
		Detail:  "The request was invalid or malformed. Please check the request parameters and try again.",
	},
	401: {
		Summary: "Unauthorized - 401",
		Detail:  "Authentication failed. Please check your Entra ID credentials and permissions.",
	},
	403: {
		Summary: "Forbidden - 403",
		Detail:  "Your credentials lack sufficient authorisation to perform this operation. Grant the required Microsoft Graph permissions to your Entra ID authentication method.",
	},
	404: {
		Summary: "Not Found - 404",
		Detail:  "The requested resource was not found.",
	},
	409: {
		Summary: "Conflict - 409",
		Detail:  "The operation failed due to a conflicts with the current state of the target resource. this might be due to multiple clients modifying the same resource simultaneously,the requested resource may not be in the state that was expected, or the request itself may create a conflict if it is completed.",
	},
	429: {
		Summary: "Too Many Requests - 429",
		Detail:  "Request throttled by Microsoft Graph API rate limits. Please try again later.",
	},
	500: {
		Summary: "Internal Server Error - 500",
		Detail:  "Microsoft Graph API encountered an internal error.. Please try again later.",
	},
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

	if apiErr, ok := err.(abstractions.ApiErrorable); ok {
		errorInfo.StatusCode = apiErr.GetStatusCode()
		errorInfo.Headers = apiErr.GetResponseHeaders()

		if headers := apiErr.GetResponseHeaders(); headers != nil {
			for _, key := range headers.ListKeys() {
				values := headers.Get(key)
				if len(values) > 0 {
					errorInfo.RequestDetails += fmt.Sprintf("%s: %v\n", key, values)
				}
			}
		}

		if odataErr, ok := err.(*odataerrors.ODataError); ok {
			errorInfo.IsODataError = true
			if mainError := odataErr.GetErrorEscaped(); mainError != nil {
				if code := mainError.GetCode(); code != nil {
					errorInfo.ErrorCode = *code
				}
				if message := mainError.GetMessage(); message != nil && *message != "" {
					errorInfo.ErrorMessage = *message
				}

				details := mainError.GetDetails()
				if len(details) > 0 {
					var detailMessages []string
					for _, detail := range details {
						if msg := detail.GetMessage(); msg != nil && *msg != "" {
							detailMessages = append(detailMessages, *msg)
						}
					}
					if len(detailMessages) > 0 {
						errorInfo.ErrorMessage += "\nDetails: " + strings.Join(detailMessages, "; ")
					}
				}

				if innerError := mainError.GetInnerError(); innerError != nil {
					if reqID := innerError.GetRequestId(); reqID != nil {
						errorInfo.AdditionalData["request_id"] = *reqID
					}
					if clientReqID := innerError.GetClientRequestId(); clientReqID != nil {
						errorInfo.AdditionalData["client_request_id"] = *clientReqID
					}
					if date := innerError.GetDate(); date != nil {
						errorInfo.AdditionalData["error_date"] = date.String()
					}
				}
			}
		} else if apiBaseErr, ok := apiErr.(*abstractions.ApiError); ok {
			if apiBaseErr.Message != "" {
				errorInfo.ErrorMessage = apiBaseErr.Message
			}
		}
	}

	// If after all processing we still don't have an error message, use the original error
	if errorInfo.ErrorMessage == "" {
		errorInfo.ErrorMessage = err.Error()
	}

	logErrorDetails(ctx, &errorInfo)
	return errorInfo
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
