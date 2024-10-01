package errors

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// HandleGraphError processes Graph API errors and dispatches them to specific handlers based on the HTTP status code.
func HandleGraphError(ctx context.Context, err error, resp interface{}, operation string, requiredPermissions []string) {
	errorInfo := GraphError(ctx, err)

	switch errorInfo.StatusCode {
	case 400:
		handleBadRequestError(ctx, errorInfo, resp)
	case 401:
		handleUnauthorizedError(ctx, errorInfo, resp, operation, requiredPermissions)
	case 403:
		handleForbiddenError(ctx, errorInfo, resp, operation, requiredPermissions)
	case 404:
		handleNotFoundError(ctx, errorInfo, resp, operation)
	case 409:
		handleConflictError(ctx, errorInfo, resp)
	case 429:
		handleTooManyRequestsError(ctx, errorInfo, resp)
	case 500:
		handleInternalServerError(ctx, errorInfo, resp)
	default:
		handleGenericError(ctx, errorInfo, resp)
	}
}

// handleBadRequestError handles a 400 Bad Request error.
func handleBadRequestError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}) {
	addErrorToDiagnostics(ctx, resp, "Bad Request", errorInfo.ErrorMessage)
}

// handleUnauthorizedError handles a 401 Unauthorized error.
func handleUnauthorizedError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}, operation string, requiredPermissions []string) {
	constructPermissionError(ctx, errorInfo, resp, operation, requiredPermissions)
}

// handleForbiddenError handles a 403 Forbidden error.
func handleForbiddenError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}, operation string, requiredPermissions []string) {
	constructPermissionError(ctx, errorInfo, resp, operation, requiredPermissions)
}

// handleNotFoundError handles a 404 Not Found error.
func handleNotFoundError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}, operation string) {
	if operation == "Read" {
		tflog.Warn(ctx, "Resource not found, removing from state")
		removeResourceFromState(ctx, resp)
	} else {
		addErrorToDiagnostics(ctx, resp, "Not Found", errorInfo.ErrorMessage)
	}
}

// handleConflictError handles a 409 Conflict error.
func handleConflictError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}) {
	addErrorToDiagnostics(ctx, resp, "Conflict", errorInfo.ErrorMessage)
}

// handleTooManyRequestsError handles a 429 Too Many Requests error.
func handleTooManyRequestsError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}) {
	addErrorToDiagnostics(ctx, resp, "Too Many Requests", errorInfo.ErrorMessage)
}

// handleInternalServerError handles a 500 Internal Server Error.
func handleInternalServerError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}) {
	addErrorToDiagnostics(ctx, resp, "Internal Server Error", errorInfo.ErrorMessage)
}

// handleGenericError handles a generic error.
func handleGenericError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}) {
	addErrorToDiagnostics(ctx, resp, fmt.Sprintf("HTTP Error %d", errorInfo.StatusCode), errorInfo.ErrorMessage)
}

// addErrorToDiagnostics adds an error message to the response diagnostics.
func addErrorToDiagnostics(ctx context.Context, resp interface{}, summary string, detail string) {
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

// removeResourceFromState removes the resource from state in case of a 404 Not Found error during a Read operation.
func removeResourceFromState(ctx context.Context, resp interface{}) {
	switch r := resp.(type) {
	case *resource.ReadResponse:
		r.State.RemoveResource(ctx)
	default:
		tflog.Error(ctx, "Cannot remove resource from state for this response type")
	}
}

// constructPermissionError constructs a more informative error message for permission errors.
func constructPermissionError(ctx context.Context, errorInfo GraphErrorInfo, resp interface{}, operation string, requiredPermissions []string) {
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

	addErrorToDiagnostics(ctx, resp, "Permission Error", fmt.Sprintf("Insufficient permissions for %s operation. %s Original error: %s",
		operation, permissionMsg, errorInfo.ErrorMessage))
}
