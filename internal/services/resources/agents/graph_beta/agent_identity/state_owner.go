package graphBetaAgentIdentity

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteOwnersToTerraform maps the fetched owners to Terraform state
// This function receives the owners collection from the API call made in crud.go
func MapRemoteOwnersToTerraform(ctx context.Context, data *AgentIdentityResourceModel, owners graphmodels.DirectoryObjectCollectionResponseable) {
	if owners == nil {
		tflog.Debug(ctx, "No owners response received, setting empty owner_ids")
		data.OwnerIds = types.SetValueMust(types.StringType, []attr.Value{})
		return
	}

	ownerObjects := owners.GetValue()
	if len(ownerObjects) == 0 {
		tflog.Debug(ctx, "No owners found for agent identity, setting empty owner_ids")
		data.OwnerIds = types.SetValueMust(types.StringType, []attr.Value{})
		return
	}

	ownerIds := make([]string, 0, len(ownerObjects))
	for _, owner := range ownerObjects {
		if owner != nil && owner.GetId() != nil {
			ownerIds = append(ownerIds, *owner.GetId())
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d owners to Terraform state", len(ownerIds)))

	data.OwnerIds = convert.GraphToFrameworkStringSet(ctx, ownerIds)
}
