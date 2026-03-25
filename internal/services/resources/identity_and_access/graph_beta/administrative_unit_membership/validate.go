package graphBetaAdministrativeUnitMembership

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates that all member IDs exist as valid directory objects
// before attempting to add them to the administrative unit.
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, memberIDs []string, diagnostics *diag.Diagnostics) bool {
	if len(memberIDs) == 0 {
		return true
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating %d member IDs", len(memberIDs)))

	validMembers := make(map[string]bool)

	// Fetch all users, groups, and devices to validate member IDs
	users, err := fetchAllUsers(ctx, client)
	if err != nil {
		diagnostics.AddError(
			"Error fetching users for validation",
			fmt.Sprintf("Could not fetch users: %s", err.Error()),
		)
		return false
	}
	for _, user := range users {
		if user.GetId() != nil {
			validMembers[*user.GetId()] = true
		}
	}

	groups, err := fetchAllGroups(ctx, client)
	if err != nil {
		diagnostics.AddError(
			"Error fetching groups for validation",
			fmt.Sprintf("Could not fetch groups: %s", err.Error()),
		)
		return false
	}
	for _, group := range groups {
		if group.GetId() != nil {
			validMembers[*group.GetId()] = true
		}
	}

	devices, err := fetchAllDevices(ctx, client)
	if err != nil {
		diagnostics.AddError(
			"Error fetching devices for validation",
			fmt.Sprintf("Could not fetch devices: %s", err.Error()),
		)
		return false
	}
	for _, device := range devices {
		if device.GetId() != nil {
			validMembers[*device.GetId()] = true
		}
	}

	// For any remaining IDs, check directoryObjects directly
	for _, memberID := range memberIDs {
		if !validMembers[memberID] {
			exists, err := checkDirectoryObject(ctx, client, memberID)
			if err != nil {
				diagnostics.AddError(
					"Error validating directory object",
					fmt.Sprintf("Could not validate directory object %s: %s", memberID, err.Error()),
				)
				return false
			}
			if exists {
				validMembers[memberID] = true
			}
		}
	}

	// Validate each member ID
	var invalidIDs []string
	for _, memberID := range memberIDs {
		if !validMembers[memberID] {
			invalidIDs = append(invalidIDs, memberID)
		}
	}

	if len(invalidIDs) > 0 {
		diagnostics.AddError(
			"Invalid member IDs",
			fmt.Sprintf("The following member IDs do not exist as valid directory objects: %v", invalidIDs),
		)
		return false
	}

	tflog.Debug(ctx, fmt.Sprintf("All %d member IDs validated successfully", len(memberIDs)))
	return true
}

// fetchAllUsers retrieves all users from the directory using pagination
func fetchAllUsers(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) ([]graphmodels.Userable, error) {
	tflog.Trace(ctx, "Fetching all users for validation")

	usersResp, err := client.Users().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}

	users := make([]graphmodels.Userable, 0)

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Userable](
		usersResp,
		client.GetAdapter(),
		graphmodels.CreateUserCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(user graphmodels.Userable) bool {
		if user != nil {
			users = append(users, user)
		}
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	tflog.Trace(ctx, fmt.Sprintf("Fetched %d users", len(users)))
	return users, nil
}

// fetchAllGroups retrieves all groups from the directory using pagination
func fetchAllGroups(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) ([]graphmodels.Groupable, error) {
	tflog.Trace(ctx, "Fetching all groups for validation")

	groupsResp, err := client.Groups().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching groups: %w", err)
	}

	groups := make([]graphmodels.Groupable, 0)

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Groupable](
		groupsResp,
		client.GetAdapter(),
		graphmodels.CreateGroupCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(group graphmodels.Groupable) bool {
		if group != nil {
			groups = append(groups, group)
		}
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("error iterating groups: %w", err)
	}

	tflog.Trace(ctx, fmt.Sprintf("Fetched %d groups", len(groups)))
	return groups, nil
}

// fetchAllDevices retrieves all devices from the directory using pagination
func fetchAllDevices(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) ([]graphmodels.Deviceable, error) {
	tflog.Trace(ctx, "Fetching all devices for validation")

	devicesResp, err := client.Devices().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching devices: %w", err)
	}

	devices := make([]graphmodels.Deviceable, 0)

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Deviceable](
		devicesResp,
		client.GetAdapter(),
		graphmodels.CreateDeviceCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(device graphmodels.Deviceable) bool {
		if device != nil {
			devices = append(devices, device)
		}
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("error iterating devices: %w", err)
	}

	tflog.Trace(ctx, fmt.Sprintf("Fetched %d devices", len(devices)))
	return devices, nil
}

// checkDirectoryObject checks if a specific directory object exists by ID
func checkDirectoryObject(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, objectID string) (bool, error) {
	tflog.Trace(ctx, fmt.Sprintf("Checking directory object: %s", objectID))

	directoryObj, err := client.DirectoryObjects().ByDirectoryObjectId(objectID).Get(ctx, nil)
	if err != nil {
		errStr := err.Error()
		if errStr == "404" || errStr == "Not Found" || 
		   (len(errStr) > 0 && (errStr[0:3] == "404" || errStr == "error status code received from the API 404")) {
			tflog.Trace(ctx, fmt.Sprintf("Directory object %s not found (404)", objectID))
			return false, nil
		}
		return false, fmt.Errorf("error fetching directory object: %w", err)
	}

	exists := directoryObj != nil && directoryObj.GetId() != nil
	tflog.Trace(ctx, fmt.Sprintf("Directory object %s exists: %v", objectID, exists))
	return exists, nil
}
