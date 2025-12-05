package graphBetaApplicationsAgentIdentityBlueprint

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteOwnersToTerraform maps the fetched owners to Terraform state
// This function receives the owners collection from the API call made in crud.go
func MapRemoteOwnersToTerraform(ctx context.Context, data *AgentIdentityBlueprintResourceModel, owners graphmodels.DirectoryObjectCollectionResponseable) {
	if owners == nil {
		tflog.Debug(ctx, "No owners response received, preserving existing state")
		return
	}

	ownerObjects := owners.GetValue()
	if len(ownerObjects) == 0 {
		tflog.Debug(ctx, "No owners found for blueprint")
		return
	}

	ownerIds := make([]string, 0, len(ownerObjects))
	for _, owner := range ownerObjects {
		if owner != nil && owner.GetId() != nil {
			ownerIds = append(ownerIds, *owner.GetId())
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d owners to Terraform state", len(ownerIds)))

	data.OwnerUserIds = convert.GraphToFrameworkStringSet(ctx, ownerIds)
}
