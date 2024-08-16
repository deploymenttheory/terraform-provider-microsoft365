package crud

import (
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
