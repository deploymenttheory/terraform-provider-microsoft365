package graphBetaAgentUser

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapSponsorsToTerraform maps sponsors from SDK response to Terraform state.
// Uses the response from GET /users/{id}/sponsors
func MapSponsorsToTerraform(ctx context.Context, data *AgentUserResourceModel, sponsorsResponse graphmodels.DirectoryObjectCollectionResponseable) {
	if sponsorsResponse == nil {
		tflog.Debug(ctx, "Sponsors response is nil, setting null sponsor_ids")
		data.SponsorIds = types.SetNull(types.StringType)
		return
	}

	sponsors := sponsorsResponse.GetValue()
	if len(sponsors) == 0 {
		tflog.Debug(ctx, "No sponsors found for agent user")
		data.SponsorIds = types.SetNull(types.StringType)
		return
	}

	sponsorIds := make([]string, 0, len(sponsors))
	for _, sponsor := range sponsors {
		if sponsor.GetId() != nil {
			sponsorIds = append(sponsorIds, *sponsor.GetId())
		}
	}

	if len(sponsorIds) == 0 {
		tflog.Debug(ctx, "No valid sponsor IDs found, setting null sponsor_ids")
		data.SponsorIds = types.SetNull(types.StringType)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapped %d sponsors to Terraform state", len(sponsorIds)))
	data.SponsorIds = convert.GraphToFrameworkStringSet(ctx, sponsorIds)
}
