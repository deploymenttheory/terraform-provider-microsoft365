package errors

import (
	"context"
	"fmt"
	"strconv"
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
}

// standardErrorDescriptions provides consistent error messaging across the provider
var standardErrorDescriptions = map[int]ErrorDescription{
	400: {
		Summary: "Bad Request - 400",
		Detail:  "The request was invalid or malformed. Please check the request parameters and try again.",
	},
	401: {
		Summary: "Unauthorized - 401",
		Detail:  "Authentication failed. Please check your credentials and permissions.",
	},
	403: {
		Summary: "Forbidden - 403",
		Detail:  "Your credentials do not have permission to perform this action.",
	},
	404: {
		Summary: "Not Found - 404",
		Detail:  "The requested resource was not found.",
	},
	409: {
		Summary: "Conflict - 409",
		Detail:  "The request conflicts with the current state of the server.",
	},
	429: {
		Summary: "Too Many Requests - 429",
		Detail:  "Too many requests. Please try again later.",
	},
	500: {
		Summary: "Internal Server Error - 500",
		Detail:  "An internal server error occurred. Please try again later.",
	},
}

// HandleGraphError processes Graph API errors and dispatches them appropriately
func HandleGraphError(ctx context.Context, err error, resp interface{}, operation string, requiredPermissions []string) {
	errorInfo := GraphError(ctx, err)
	errorDesc := getErrorDescription(errorInfo.StatusCode)

	// Log the handling attempt
	tflog.Debug(ctx, "Handling Graph error", map[string]interface{}{
		"status_code": errorInfo.StatusCode,
		"operation":   operation,
	})

	// Handle special cases first
	switch errorInfo.StatusCode {
	case 404:
		if operation == "Read" {
			tflog.Warn(ctx, "Resource not found, removing from state")
			removeResourceFromState(ctx, resp)
			return
		}
		addErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructErrorDetail(errorDesc.Detail, errorInfo.ErrorMessage))

	case 401, 403:
		handlePermissionError(ctx, errorInfo, resp, operation, requiredPermissions)

	default:
		// Handle all other cases
		addErrorToDiagnostics(ctx, resp, errorDesc.Summary,
			constructErrorDetail(errorDesc.Detail, errorInfo.ErrorMessage))
	}
}

// Utility functions

// parseStatusCode extracts the status code from an error message
func parseStatusCode(errMsg string) int {
	if strings.Contains(errMsg, "status code") {
		parts := strings.Split(errMsg, "status code")
		if len(parts) > 1 {
			remaining := strings.Trim(parts[1], " :")
			for _, word := range strings.Fields(remaining) {
				if code, err := strconv.Atoi(word); err == nil {
					return code
				}
			}
		}
	}
	return 0
}

// GraphError extracts and processes error information from Graph API errors
func GraphError(ctx context.Context, err error) GraphErrorInfo {
	errorInfo := GraphErrorInfo{
		AdditionalData: make(map[string]interface{}),
	}

	if err == nil {
		return errorInfo
	}

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
				if message := mainError.GetMessage(); message != nil {
					errorInfo.ErrorMessage = *message
				}
			}
			errorInfo.AdditionalData = odataErr.GetAdditionalData()
		} else if apiBaseErr, ok := apiErr.(*abstractions.ApiError); ok {
			errorInfo.ErrorMessage = apiBaseErr.Message
		}
	} else {
		// Handle non-API errors
		errorInfo.ErrorMessage = err.Error()
		errorInfo.StatusCode = parseStatusCode(err.Error())
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
	if len(requiredPermissions) > 0 {
		if strings.ToLower(operation) == "read" {
			readPerm := strings.Replace(requiredPermissions[0], "ReadWrite", "Read", 1)
			permissionMsg = fmt.Sprintf("Required permissions: %s or %s", readPerm, requiredPermissions[0])
		} else {
			permissionMsg = fmt.Sprintf("Required permissions: %s", strings.Join(requiredPermissions, ", "))
		}
	} else {
		permissionMsg = "No specific permissions provided."
	}

	errorDesc := getErrorDescription(errorInfo.StatusCode)
	detail := fmt.Sprintf("%s\n%s\nOriginal error: %s",
		errorDesc.Detail,
		permissionMsg,
		errorInfo.ErrorMessage)

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
