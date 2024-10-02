package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CustomError represents the structure of the unexpected error
type CustomError struct {
	Version              int    `json:"_version"`
	Message              string `json:"Message"`
	CustomApiErrorPhrase string `json:"CustomApiErrorPhrase"`
	RetryAfter           any    `json:"RetryAfter"`
	ErrorSourceService   string `json:"ErrorSourceService"`
	HttpHeaders          string `json:"HttpHeaders"`
}

func main() {
	// Replace these with your actual values
	tenantID := os.Getenv("M365_TENANT_ID")
	clientID := os.Getenv("M365_CLIENT_ID")
	clientSecret := os.Getenv("M365_CLIENT_SECRET")
	resourceID := "1236b4bb-7b18-4615-ab5f-9634f409a93c" // Replace with your actual resource ID

	// Create a credential using client secret
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		log.Fatalf("Error creating credential: %v", err)
	}

	// Create a new Graph client
	client, err := msgraphbetasdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		log.Fatalf("Error creating Graph client: %v", err)
	}

	// Attempt to get the resource
	result, err := client.DeviceManagement().AssignmentFilters().
		ByDeviceAndAppManagementAssignmentFilterId(resourceID).
		Get(context.Background(), nil)

	if err != nil {
		handleError(err, resourceID)
		return
	}

	// If we get here, the resource exists
	fmt.Printf("Resource with ID %s exists\n", resourceID)
	fmt.Printf("Resource details: %+v\n", result)
}

func handleError(err error, resourceID string) {
	fmt.Printf("Error occurred while checking resource %s:\n", resourceID)
	fmt.Printf("Error type: %T\n", err)

	if odataErr, ok := err.(*odataerrors.ODataError); ok {
		handleODataError(odataErr, resourceID)
	} else {
		fmt.Printf("Unexpected error: %v\n", err)
	}
}

func handleODataError(odataErr *odataerrors.ODataError, resourceID string) {
	errMsg := odataErr.Error()
	fmt.Printf("OData Error detected\n")

	// Try to get the status code
	statusCode := getStatusCode(odataErr)
	fmt.Printf("HTTP Status Code: %d\n", statusCode)

	var customErr CustomError
	if err := json.Unmarshal([]byte(errMsg), &customErr); err == nil {
		fmt.Printf("Parsed Custom Error:\n")
		fmt.Printf("Version: %d\n", customErr.Version)
		fmt.Printf("Message: %s\n", customErr.Message)
		fmt.Printf("CustomApiErrorPhrase: %s\n", customErr.CustomApiErrorPhrase)
		fmt.Printf("RetryAfter: %v\n", customErr.RetryAfter)
		fmt.Printf("ErrorSourceService: %s\n", customErr.ErrorSourceService)
		fmt.Printf("HttpHeaders: %s\n", customErr.HttpHeaders)
	} else {
		fmt.Printf("Could not parse error as CustomError: %v\n", err)
		fmt.Printf("Raw error message: %s\n", errMsg)
	}

	// Interpret the status code
	switch statusCode {
	case http.StatusNotFound:
		fmt.Printf("Resource with ID %s does not exist.\n", resourceID)
	case http.StatusForbidden:
		fmt.Printf("Access to resource with ID %s is forbidden. Check your permissions.\n", resourceID)
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized access to resource with ID %s. Check your credentials.\n", resourceID)
	default:
		fmt.Printf("Unable to determine the status of the resource. Please check the error details.\n")
	}
}

func getStatusCode(err error) int {
	if respErr, ok := err.(interface{ ResponseStatusCode() int }); ok {
		return respErr.ResponseStatusCode()
	}
	return 0
}
