package graphbetabrowsersite

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *BrowserSiteListResourceModel, remoteResource graphmodels.BrowserSiteListable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayName()))
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.PublishedDateTime = state.TimeToString(remoteResource.GetPublishedDateTime())
	data.Revision = types.StringValue(state.StringPtrToString(remoteResource.GetRevision()))
	data.Status = state.EnumPtrToTypeString(remoteResource.GetStatus())

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
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
		DisplayName: types.StringValue(state.StringPtrToString(identity.GetDisplayName())),
		ID:          types.StringValue(state.StringPtrToString(identity.GetId())),
	}
}
