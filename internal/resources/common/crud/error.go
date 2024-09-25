package crud

import (
	"fmt"
	"strings"

	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// IsNotFoundError checks if the given error is an OData error indicating that a resource was not found.
// The function first verifies if the error is not nil. Then, it attempts to cast the error to an ODataError
// type from the Microsoft Graph SDK. If the casting is successful, the function retrieves the main error
// details using the GetErrorEscaped method of the ODataError struct. It then checks if the error code or
// message contains indications of a "not found" error.
//
// Specifically, the function looks for the error codes "request_resourcenotfound" and "resourcenotfound"
// (case-insensitive), or a message containing the phrase "not found" (case-insensitive). If any of these
// conditions are met, the function returns true, indicating that the error is a "not found" error.
// Otherwise, it returns false.
//
// The ODataError struct is part of the Microsoft Graph SDK and includes various methods and properties
// to handle API errors. The main error details are encapsulated in a nested structure that provides
// additional context, such as error codes and descriptive messages.
//
// Usage:
//
//	if common.IsNotFoundError(err) {
//	    // Handle the "not found" error case
//	}
//
// Parameters:
//
//	err - The error to check.
//
// Returns:
//
//	bool - True if the error indicates that a resource was not found, otherwise false.
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	odataErr, ok := err.(*odataerrors.ODataError)
	if !ok {
		// If it's not an ODataError, check the error string
		return strings.Contains(strings.ToLower(err.Error()), "not found")
	}

	mainError := odataErr.GetErrorEscaped()
	if mainError != nil {
		if code := mainError.GetCode(); code != nil {
			switch strings.ToLower(*code) {
			case "request_resourcenotfound", "resourcenotfound", "notfound":
				return true
			}
		}

		if message := mainError.GetMessage(); message != nil {
			lowerMessage := strings.ToLower(*message)
			return strings.Contains(lowerMessage, "not found") ||
				strings.Contains(lowerMessage, "could not be found") ||
				strings.Contains(lowerMessage, "does not exist")
		}
	}

	return false
}

// PermissionError checks if the given error is related to insufficient permissions
// and returns a more informative error message if it is.
//
// This function is designed to work with the generic error messages returned by the Microsoft Graph API,
// which often do not provide detailed information about the specific permissions required.
//
// The function checks for common phrases indicating a permission issue, such as "permission",
// "access denied", or "unauthorized". If such a phrase is found, it constructs a more detailed
// error message that includes the operation being performed and the permissions required for that operation.
//
// For Read operations, the function will specify both Read and ReadWrite permissions as possible
// requirements. For all other operations, only the ReadWrite permission will be specified.
//
// Parameters:
//   - err: The original error returned by the API call.
//   - operation: A string describing the operation being performed (e.g., "Create", "Read", "Update", "Delete").
//   - requiredPermissions: A slice of strings listing the permissions required for the operation.
//     The first permission in this slice is assumed to be the ReadWrite permission.
//
// Returns:
//   - An error with a more detailed message if a permission issue is detected.
//   - The original error if no permission issue is detected.
//
// Usage:
//
//	err = crud.PermissionError(err, "Create", r.WritePermissions)
//	if err != nil {
//	    resp.Diagnostics.AddError("Error creating browser site list", err.Error())
//	    return
//	}
//
// For a Read operation, the usage would be:
//
//	err = crud.HandlePermissionError(err, "Read", r.WritePermissions)
//	// The function will automatically derive the Read permission from the WritePermissions
func PermissionError(err error, operation string, requiredPermissions []string) error {
	if err == nil {
		return nil
	}

	errorMsg := err.Error()
	lowerErrorMsg := strings.ToLower(errorMsg)

	if strings.Contains(lowerErrorMsg, "permission") ||
		strings.Contains(lowerErrorMsg, "access denied") ||
		strings.Contains(lowerErrorMsg, "unauthorized") {
		var permissionMsg string
		if strings.ToLower(operation) == "read" {
			readPerm := strings.Replace(requiredPermissions[0], "ReadWrite", "Read", 1)
			permissionMsg = fmt.Sprintf("Required permissions: %s or %s", readPerm, requiredPermissions[0])
		} else {
			permissionMsg = fmt.Sprintf("Required permissions: %s", strings.Join(requiredPermissions, ", "))
		}

		return fmt.Errorf("insufficient permissions for %s operation. %s. Original error: %s",
			operation, permissionMsg, errorMsg)
	}

	return err
}
