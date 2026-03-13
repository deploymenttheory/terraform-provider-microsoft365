package graphBetaApplicationsAgentIdentityBlueprint

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates the agent identity blueprint request by checking sponsors and owners
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *AgentIdentityBlueprintResourceModel) error {
	tflog.Debug(ctx, "Starting validation of agent identity blueprint request")

	if err := validateSponsorIsTypeUser(ctx, client, data.SponsorUserIds); err != nil {
		return fmt.Errorf("sponsor validation failed: %w", err)
	}

	if err := validateOwnerIsTypeUser(ctx, client, data.OwnerUserIds); err != nil {
		return fmt.Errorf("owner validation failed: %w", err)
	}

	tflog.Debug(ctx, "Successfully validated agent identity blueprint request")
	return nil
}

// validateSponsorIsTypeUser validates that all sponsor IDs are users
func validateSponsorIsTypeUser(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, sponsorUserIds types.Set) error {
	if sponsorUserIds.IsNull() || sponsorUserIds.IsUnknown() {
		return sentinels.ErrSponsorUserIDsNullOrUnknown
	}

	var sponsorIds []string
	diags := sponsorUserIds.ElementsAs(ctx, &sponsorIds, false)
	if diags.HasError() {
		return fmt.Errorf("failed to extract sponsor_user_ids: %v", diags.Errors()[0])
	}

	if len(sponsorIds) == 0 {
		return sentinels.ErrAtLeastOneSponsorRequired
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating %d sponsor IDs", len(sponsorIds)))

	for _, sponsorID := range sponsorIds {
		user, err := client.
			Users().
			ByUserId(sponsorID).
			Get(ctx, nil)

		if err != nil {
			return fmt.Errorf("sponsor ID %s is not a valid user: %w", sponsorID, err)
		}

		if user.GetId() == nil {
			return fmt.Errorf("%w for sponsor ID %s", sentinels.ErrSponsorUserObjectNull, sponsorID)
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully validated sponsor ID %s as user", sponsorID))
	}

	return nil
}

// validateOwnerIsTypeUser validates that all owner IDs are users
func validateOwnerIsTypeUser(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, ownerUserIds types.Set) error {
	if ownerUserIds.IsNull() || ownerUserIds.IsUnknown() {
		return sentinels.ErrOwnerUserIDsNullOrUnknown
	}

	var ownerIds []string
	diags := ownerUserIds.ElementsAs(ctx, &ownerIds, false)
	if diags.HasError() {
		return fmt.Errorf("failed to extract owner_user_ids: %v", diags.Errors()[0])
	}

	if len(ownerIds) == 0 {
		return sentinels.ErrAtLeastOneOwnerRequired
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating %d owner IDs", len(ownerIds)))

	for _, ownerID := range ownerIds {
		user, err := client.
			Users().
			ByUserId(ownerID).
			Get(ctx, nil)

		if err != nil {
			return fmt.Errorf("owner ID %s is not a valid user: %w", ownerID, err)
		}

		if user.GetId() == nil {
			return fmt.Errorf("%w for owner ID %s", sentinels.ErrOwnerUserObjectNull, ownerID)
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully validated owner ID %s as user", ownerID))
	}

	return nil
}
