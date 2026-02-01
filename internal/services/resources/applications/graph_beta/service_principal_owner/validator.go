package graphBetaServicePrincipalOwner

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// ValidateServicePrincipalOwner validates that the owner assignment is compatible with the target service principal
func ValidateServicePrincipalOwner(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data ServicePrincipalOwnerResourceModel, isUpdate bool) error {
	tflog.Debug(ctx, "Starting service principal owner assignment validation")

	servicePrincipalId := data.ServicePrincipalID.ValueString()
	ownerId := data.OwnerID.ValueString()
	ownerObjectType := data.OwnerObjectType.ValueString()

	// Validate that the service principal exists
	servicePrincipal, err := client.ServicePrincipals().ByServicePrincipalId(servicePrincipalId).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve service principal %s: %v", servicePrincipalId, err)
	}

	if servicePrincipal == nil {
		return fmt.Errorf("service principal %s not found", servicePrincipalId)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validated service principal %s exists", servicePrincipalId))

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
		// During updates, we might be removing an owner, so we should check if this would leave the service principal without any user owners
		owners, err := client.
			ServicePrincipals().
			ByServicePrincipalId(servicePrincipalId).
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
				tflog.Warn(ctx, "Warning: This operation might remove the last user owner from the service principal, which is not recommended")
			}
		}
	}

	tflog.Debug(ctx, "Service principal owner assignment validation completed successfully")
	return nil
}
