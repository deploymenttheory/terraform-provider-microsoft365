package graphBetaAgentsAgentCollection

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the API response to the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AgentCollectionResourceModel, agentCollection graphmodels.AgentCollectionable) {
	if agentCollection == nil {
		return
	}

	tflog.Debug(ctx, "Mapping agent collection response to Terraform state")

	// Map basic properties
	data.ID = convert.GraphToFrameworkString(agentCollection.GetId())
	data.DisplayName = convert.GraphToFrameworkString(agentCollection.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(agentCollection.GetDescription())
	data.ManagedBy = convert.GraphToFrameworkString(agentCollection.GetManagedBy())
	data.OriginatingStore = convert.GraphToFrameworkString(agentCollection.GetOriginatingStore())
	data.CreatedBy = convert.GraphToFrameworkString(agentCollection.GetCreatedBy())

	// Map timestamps
	data.CreatedDateTime = convert.GraphToFrameworkTime(agentCollection.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(agentCollection.GetLastModifiedDateTime())

	// Map owner IDs
	data.OwnerIds = convert.GraphToFrameworkStringSet(ctx, agentCollection.GetOwnerIds())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping agent collection response to Terraform state for ID: %s", data.ID.ValueString()))
}
