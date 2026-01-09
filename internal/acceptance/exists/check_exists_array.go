package exists

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// GraphClientArrayMembershipCheckFunc is a function type that retrieves a parent resource using the Graph SDK client.
// The function should return the parent resource that contains the array to be searched.
type GraphClientArrayMembershipCheckFunc func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, parentID string) (any, error)

// CheckResourceExistsByArrayMembership checks if a specific value exists in an array field of a parent resource using GraphClient.
//
// Example: For checking if a license SKU exists in a user's assignedLicenses array:
// - Callback fetches the user by ID
// - Helper searches assignedLicenses array for matching skuId
//
// Parameters:
//   - ctx: Context for the request
//   - state: Terraform instance state containing parent ID and search value
//   - parentIDAttribute: The attribute name containing the parent resource ID (e.g., "user_id", "blueprint_id")
//   - arrayFieldName: The JSON field name of the array to search (e.g., "assignedLicenses", "keyCredentials")
//   - searchFieldName: The JSON field name within array items to match (e.g., "skuId", "keyId")
//   - searchValueAttribute: The attribute name containing the value to search for (e.g., "sku_id", "key_id")
//   - checkFunc: A function that retrieves the parent resource using the SDK
//
// Returns:
//   - *bool: Pointer to boolean indicating if the value exists in the array
//   - error: Error if the check fails or required attributes are missing
//
// Example usage:
//
//	return exists.CheckResourceExistsByArrayMembership(ctx, state, "user_id", "assignedLicenses", "skuId", "sku_id", func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, parentID string) (any, error) {
//	    return client.Users().ByUserId(parentID).Get(ctx, &users.UserItemRequestBuilderGetRequestConfiguration{
//	        QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{Select: []string{"id", "assignedLicenses"}},
//	    })
//	})
func CheckResourceExistsByArrayMembership(
	ctx context.Context,
	state *terraform.InstanceState,
	parentIDAttribute string,
	arrayFieldName string,
	searchFieldName string,
	searchValueAttribute string,
	checkFunc GraphClientArrayMembershipCheckFunc,
) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	parentID := state.Attributes[parentIDAttribute]
	searchValue := state.Attributes[searchValueAttribute]

	if parentID == "" {
		return nil, fmt.Errorf("%s not found in state", parentIDAttribute)
	}
	if searchValue == "" {
		return nil, fmt.Errorf("%s not found in state", searchValueAttribute)
	}

	resource, err := checkFunc(graphClient, ctx, parentID)
	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		if errorInfo.StatusCode == 404 ||
			errorInfo.ErrorCode == "ResourceNotFound" ||
			errorInfo.ErrorCode == "Request_ResourceNotFound" ||
			errorInfo.ErrorCode == "ItemNotFound" {
			exists := false
			return &exists, nil
		}
		return nil, err
	}

	// Convert resource to JSON for easier navigation
	resourceJSON, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}

	var resourceMap map[string]any
	if err := json.Unmarshal(resourceJSON, &resourceMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	arrayRaw, ok := resourceMap[arrayFieldName]
	if !ok {
		exists := false
		return &exists, nil
	}

	array, ok := arrayRaw.([]any)
	if !ok || len(array) == 0 {
		exists := false
		return &exists, nil
	}

	for _, item := range array {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		if fieldValue, ok := itemMap[searchFieldName].(string); ok && fieldValue == searchValue {
			exists := true
			return &exists, nil
		}
	}

	exists := false
	return &exists, nil
}

// CheckResourceExistsByStringArrayMembership checks if a specific string value exists in a string array field of a parent resource using GraphClient.
//
// Example: For checking if a specific URI exists in an application's identifierUris array:
// - Callback fetches the application by ID
// - Helper searches identifierUris array for matching string
//
// Parameters:
//   - ctx: Context for the request
//   - state: Terraform instance state containing parent ID and search value
//   - parentIDAttribute: The attribute name containing the parent resource ID (e.g., "blueprint_id")
//   - arrayFieldName: The JSON field name of the string array to search (e.g., "identifierUris")
//   - searchValueAttribute: The attribute name containing the value to search for (e.g., "identifier_uri")
//   - checkFunc: A function that retrieves the parent resource using the SDK
//
// Returns:
//   - *bool: Pointer to boolean indicating if the value exists in the array
//   - error: Error if the check fails or required attributes are missing
//
// Example usage:
//
//	return exists.CheckResourceExistsByStringArrayMembership(ctx, state, "blueprint_id", "identifierUris", "identifier_uri", func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, parentID string) (any, error) {
//	    return client.Applications().ByApplicationId(parentID).Get(ctx, &applications.ApplicationItemRequestBuilderGetRequestConfiguration{
//	        QueryParameters: &applications.ApplicationItemRequestBuilderGetQueryParameters{Select: []string{"id", "identifierUris"}},
//	    })
//	})
func CheckResourceExistsByStringArrayMembership(
	ctx context.Context,
	state *terraform.InstanceState,
	parentIDAttribute string,
	arrayFieldName string,
	searchValueAttribute string,
	checkFunc GraphClientArrayMembershipCheckFunc,
) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	parentID := state.Attributes[parentIDAttribute]
	searchValue := state.Attributes[searchValueAttribute]

	if parentID == "" {
		return nil, fmt.Errorf("%s not found in state", parentIDAttribute)
	}
	if searchValue == "" {
		return nil, fmt.Errorf("%s not found in state", searchValueAttribute)
	}

	// Call the check function to get the parent resource
	resource, err := checkFunc(graphClient, ctx, parentID)
	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		if errorInfo.StatusCode == 404 ||
			errorInfo.ErrorCode == "ResourceNotFound" ||
			errorInfo.ErrorCode == "Request_ResourceNotFound" ||
			errorInfo.ErrorCode == "ItemNotFound" {
			exists := false
			return &exists, nil
		}
		return nil, err
	}

	// Convert resource to JSON for easier navigation
	resourceJSON, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}

	var resourceMap map[string]any
	if err := json.Unmarshal(resourceJSON, &resourceMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	arrayRaw, ok := resourceMap[arrayFieldName]
	if !ok {
		exists := false
		return &exists, nil
	}

	array, ok := arrayRaw.([]any)
	if !ok || len(array) == 0 {
		exists := false
		return &exists, nil
	}

	for _, item := range array {
		if strValue, ok := item.(string); ok && strValue == searchValue {
			exists := true
			return &exists, nil
		}
	}

	exists := false
	return &exists, nil
}
