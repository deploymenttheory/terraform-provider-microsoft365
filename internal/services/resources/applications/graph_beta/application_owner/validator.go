package graphBetaApplicationOwner

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// ValidateApplicationOwner validates that the owner assignment is compatible with the target application
func ValidateApplicationOwner(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data ApplicationOwnerResourceModel, isUpdate bool) error {
	tflog.Debug(ctx, "Starting application owner assignment validation")

	applicationId := data.ApplicationID.ValueString()
	ownerId := data.OwnerID.ValueString()
	ownerObjectType := data.OwnerObjectType.ValueString()

	// Validate that the application exists
	application, err := client.Applications().ByApplicationId(applicationId).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve application %s: %v", applicationId, err)
	}

	if application == nil {
		return fmt.Errorf("application %s not found", applicationId)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validated application %s exists", applicationId))

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
		// During updates, we might be removing an owner, so we should check if this would leave the application without any user owners
		owners, err := client.
			Applications().
			ByApplicationId(applicationId).
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
				tflog.Warn(ctx, "Warning: This operation might remove the last user owner from the application, which is not recommended")
			}
		}
	}

	tflog.Debug(ctx, "Application owner assignment validation completed successfully")
	return nil
}
