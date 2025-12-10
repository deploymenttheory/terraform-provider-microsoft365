package graphBetaApplicationsAgentIdentityBlueprint

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapSponsorIdsToTerraform maps sponsor IDs from raw JSON response to Terraform state.
// It filters the API response to only include sponsors that are explicitly configured
// in Terraform, ignoring any auto-added sponsors.
func MapSponsorIdsToTerraform(ctx context.Context, data *AgentIdentityBlueprintResourceModel, sponsorResponse json.RawMessage) {
	if len(sponsorResponse) == 0 {
		tflog.Debug(ctx, "No sponsor response received")
		return
	}

	var sponsorData struct {
		Value []struct {
			ID string `json:"id"`
		} `json:"value"`
	}

	if err := json.Unmarshal(sponsorResponse, &sponsorData); err != nil {
		tflog.Warn(ctx, "Failed to unmarshal sponsors response", map[string]any{"error": err.Error()})
		return
	}

	if len(sponsorData.Value) == 0 {
		tflog.Debug(ctx, "No sponsors found for blueprint")
		data.SponsorUserIds = types.SetNull(types.StringType)
		return
	}

	var configuredSponsorIds []string
	if !data.SponsorUserIds.IsNull() && !data.SponsorUserIds.IsUnknown() {
		diags := data.SponsorUserIds.ElementsAs(ctx, &configuredSponsorIds, false)
		if diags.HasError() {
			tflog.Warn(ctx, "Failed to extract configured sponsor IDs, falling back to all sponsors")
			configuredSponsorIds = nil
		}
	}

	configuredSet := make(map[string]bool)
	for _, id := range configuredSponsorIds {
		configuredSet[id] = true
	}

	filteredSponsorIds := make([]string, 0)
	for _, s := range sponsorData.Value {
		if len(configuredSet) == 0 || configuredSet[s.ID] {
			filteredSponsorIds = append(filteredSponsorIds, s.ID)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d sponsors to Terraform state (filtered from %d total)", len(filteredSponsorIds), len(sponsorData.Value)))

	if len(filteredSponsorIds) > 0 {
		data.SponsorUserIds = convert.GraphToFrameworkStringSet(ctx, filteredSponsorIds)
	} else {
		data.SponsorUserIds = types.SetNull(types.StringType)
	}
}
