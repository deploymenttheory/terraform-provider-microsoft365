package graphBetaApplicationsAgentIdentityBlueprint

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteOwnersToTerraform maps the fetched owners to Terraform state.
// It filters the API response to only include owners that are explicitly configured
// in Terraform, ignoring auto-added owners like the app registration caller.
func MapRemoteOwnersToTerraform(ctx context.Context, data *AgentIdentityBlueprintResourceModel, owners graphmodels.DirectoryObjectCollectionResponseable) {
	if owners == nil {
		tflog.Debug(ctx, "No owners response received, preserving existing state")
		return
	}

	ownerObjects := owners.GetValue()
	if len(ownerObjects) == 0 {
		tflog.Debug(ctx, "No owners found for blueprint")
		data.OwnerUserIds = types.SetNull(types.StringType)
		return
	}

	var configuredOwnerIds []string
	if !data.OwnerUserIds.IsNull() && !data.OwnerUserIds.IsUnknown() {
		diags := data.OwnerUserIds.ElementsAs(ctx, &configuredOwnerIds, false)
		if diags.HasError() {
			tflog.Warn(ctx, "Failed to extract configured owner IDs, falling back to all owners")
			configuredOwnerIds = nil
		}
	}

	configuredSet := make(map[string]bool)
	for _, id := range configuredOwnerIds {
		configuredSet[id] = true
	}

	filteredOwnerIds := make([]string, 0)
	for _, owner := range ownerObjects {
		if owner != nil && owner.GetId() != nil {
			ownerId := *owner.GetId()
			if len(configuredSet) == 0 || configuredSet[ownerId] {
				filteredOwnerIds = append(filteredOwnerIds, ownerId)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d owners to Terraform state (filtered from %d total)", len(filteredOwnerIds), len(ownerObjects)))

	if len(filteredOwnerIds) > 0 {
		data.OwnerUserIds = convert.GraphToFrameworkStringSet(ctx, filteredOwnerIds)
	} else {
		data.OwnerUserIds = types.SetNull(types.StringType)
	}
}
