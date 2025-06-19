package graphBetaBrowserSiteList

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *BrowserSiteListResourceModel, remoteResource graphmodels.BrowserSiteListable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.PublishedDateTime = convert.GraphToFrameworkTime(remoteResource.GetPublishedDateTime())
	data.Revision = convert.GraphToFrameworkString(remoteResource.GetRevision())
	data.Status = convert.GraphToFrameworkEnum(remoteResource.GetStatus())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}

func MapIdentitySetRemoteStateToTerraform(identitySet graphmodels.IdentitySetable) sharedmodels.IdentitySetResourceModel {
	if identitySet == nil {
		return sharedmodels.IdentitySetResourceModel{}
	}

	return sharedmodels.IdentitySetResourceModel{
		Application: MapIdentityRemoteStateToTerraform(identitySet.GetApplication()),
		User:        MapIdentityRemoteStateToTerraform(identitySet.GetUser()),
		Device:      MapIdentityRemoteStateToTerraform(identitySet.GetDevice()),
	}
}

func MapIdentityRemoteStateToTerraform(identity graphmodels.Identityable) sharedmodels.IdentityResourceModel {
	if identity == nil {
		return sharedmodels.IdentityResourceModel{}
	}

	return sharedmodels.IdentityResourceModel{
		DisplayName: convert.GraphToFrameworkString(identity.GetDisplayName()),
		ID:          convert.GraphToFrameworkString(identity.GetId()),
	}
}
