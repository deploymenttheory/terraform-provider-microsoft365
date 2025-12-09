package graphBetaAgentIdentity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapSponsorIdsToTerraform maps sponsor IDs from raw JSON response to Terraform state
// This function receives the raw JSON response from the custom API call in crud.go
// and handles all unmarshaling and mapping logic
func MapSponsorIdsToTerraform(ctx context.Context, data *AgentIdentityResourceModel, sponsorResponse json.RawMessage) {
	if len(sponsorResponse) == 0 {
		tflog.Debug(ctx, "No sponsor response received, setting empty sponsor_ids")
		data.SponsorIds = types.SetValueMust(types.StringType, []attr.Value{})
		return
	}

	var sponsorData struct {
		Value []struct {
			ID string `json:"id"`
		} `json:"value"`
	}

	if err := json.Unmarshal(sponsorResponse, &sponsorData); err != nil {
		tflog.Warn(ctx, "Failed to unmarshal sponsors response", map[string]any{"error": err.Error()})
		data.SponsorIds = types.SetValueMust(types.StringType, []attr.Value{})
		return
	}

	if len(sponsorData.Value) == 0 {
		tflog.Debug(ctx, "No sponsors found for agent identity, setting empty sponsor_ids")
		data.SponsorIds = types.SetValueMust(types.StringType, []attr.Value{})
		return
	}

	sponsorIds := make([]string, len(sponsorData.Value))
	for i, s := range sponsorData.Value {
		sponsorIds[i] = s.ID
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d sponsors to Terraform state", len(sponsorIds)))

	data.SponsorIds = convert.GraphToFrameworkStringSet(ctx, sponsorIds)
}
