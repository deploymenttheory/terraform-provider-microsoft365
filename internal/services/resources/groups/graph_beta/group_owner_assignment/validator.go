package graphBetaGroupOwnerAssignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// ValidateGroupOwnerAssignment validates that the owner assignment is compatible with the target group
func ValidateGroupOwnerAssignment(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data GroupOwnerAssignmentResourceModel, isUpdate bool) error {
	tflog.Debug(ctx, "Starting group owner assignment validation")

	groupId := data.GroupID.ValueString()
	ownerId := data.OwnerID.ValueString()
	ownerObjectType := data.OwnerObjectType.ValueString()

	// Validate that the group exists and get its properties
	group, err := client.Groups().ByGroupId(groupId).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve group %s: %v", groupId, err)
	}

	if group == nil {
		return fmt.Errorf("group %s not found", groupId)
	}

	// Get group type information
	groupTypes := group.GetGroupTypes()
	mailEnabled := group.GetMailEnabled()
	securityEnabled := group.GetSecurityEnabled()

	var groupType string
	if groupTypes != nil && len(groupTypes) > 0 {
		for _, gt := range groupTypes {
			if gt == "Unified" {
				groupType = "Microsoft 365"
				break
			} else if gt == "DynamicMembership" {
				// This is handled separately, but note it
				continue
			}
		}
	}

	if groupType == "" {
		if mailEnabled != nil && securityEnabled != nil {
			if *mailEnabled && *securityEnabled {
				groupType = "Mail-enabled Security"
			} else if !*mailEnabled && *securityEnabled {
				groupType = "Security"
			} else if *mailEnabled && !*securityEnabled {
				groupType = "Distribution"
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Group %s is of type: %s", groupId, groupType))

	// Check if the group type supports owner management
	if groupType == "Mail-enabled Security" || groupType == "Distribution" {
		return fmt.Errorf("cannot manage owners for %s groups - they are read-only", groupType)
	}

	// Validate that the owner object exists and is of the correct type
	switch ownerObjectType {
	case "User":
		user, err := client.
			Users().
			ByUserId(ownerId).
			Get(ctx, nil)

		if err != nil {
			return fmt.Errorf("failed to retrieve user %s: %v", ownerId, err)
		}
		if user == nil {
			return fmt.Errorf("user %s not found", ownerId)
		}
		tflog.Debug(ctx, fmt.Sprintf("Validated user %s exists", ownerId))

	case "ServicePrincipal":
		sp, err := client.
			ServicePrincipals().
			ByServicePrincipalId(ownerId).
			Get(ctx, nil)

		if err != nil {
			return fmt.Errorf("failed to retrieve service principal %s: %v", ownerId, err)
		}

		if sp == nil {
			return fmt.Errorf("service principal %s not found", ownerId)
		}
		tflog.Debug(ctx, fmt.Sprintf("Validated service principal %s exists", ownerId))

	default:
		return fmt.Errorf("unsupported owner object type: %s", ownerObjectType)
	}

	// Additional validation: Check if trying to remove the last user owner
	if isUpdate {
		// During updates, we might be removing an owner, so we should check if this would leave the group without any user owners
		owners, err := client.
			Groups().
			ByGroupId(groupId).
			Owners().
			Get(ctx, nil)

		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Could not retrieve current owners for validation: %v", err))
			// Don't fail the validation, just warn
		} else if owners != nil && owners.GetValue() != nil {
			userOwnerCount := 0
			for _, owner := range owners.GetValue() {
				if owner.GetOdataType() != nil && *owner.GetOdataType() == "#microsoft.graph.user" {
					userOwnerCount++
				}
			}

			if userOwnerCount <= 1 && ownerObjectType == "User" {
				tflog.Warn(ctx, "Warning: This operation might remove the last user owner from the group, which is not allowed by Microsoft Graph")
			}
		}
	}

	tflog.Debug(ctx, "Group owner assignment validation completed successfully")
	return nil
}
