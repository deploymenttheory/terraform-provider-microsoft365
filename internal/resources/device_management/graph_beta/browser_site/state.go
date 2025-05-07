package graphBetaBrowserSite

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *BrowserSiteResourceModel, remoteResource graphmodels.BrowserSiteable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.AllowRedirect = types.BoolPointerValue(remoteResource.GetAllowRedirect())
	data.Comment = types.StringPointerValue(remoteResource.GetComment())
	data.CompatibilityMode = state.EnumPtrToTypeString(remoteResource.GetCompatibilityMode())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.DeletedDateTime = state.TimeToString(remoteResource.GetDeletedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.MergeType = state.EnumPtrToTypeString(remoteResource.GetMergeType())
	data.Status = state.EnumPtrToTypeString(remoteResource.GetStatus())
	data.TargetEnvironment = state.EnumPtrToTypeString(remoteResource.GetTargetEnvironment())
	data.WebUrl = types.StringPointerValue(remoteResource.GetWebUrl())

	history := remoteResource.GetHistory()
	if len(history) == 0 {
		data.History = []BrowserSiteHistoryResourceModel{}
	} else {
		data.History = make([]BrowserSiteHistoryResourceModel, len(history))
		for i, historyItem := range history {
			data.History[i] = MapHistoryRemoteStateToTerraform(historyItem)
		}
	}

	if lastModifiedBy := remoteResource.GetLastModifiedBy(); lastModifiedBy != nil {
		data.LastModifiedBy = MapIdentitySetRemoteStateToTerraform(lastModifiedBy)
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}

func MapHistoryRemoteStateToTerraform(historyItem graphmodels.BrowserSiteHistoryable) BrowserSiteHistoryResourceModel {
	return BrowserSiteHistoryResourceModel{
		AllowRedirect:     types.BoolPointerValue(historyItem.GetAllowRedirect()),
		Comment:           types.StringPointerValue(historyItem.GetComment()),
		CompatibilityMode: state.EnumPtrToTypeString(historyItem.GetCompatibilityMode()),
		LastModifiedBy:    MapIdentitySetRemoteStateToTerraform(historyItem.GetLastModifiedBy()),
		MergeType:         state.EnumPtrToTypeString(historyItem.GetMergeType()),
		PublishedDateTime: state.TimeToString(historyItem.GetPublishedDateTime()),
		TargetEnvironment: state.EnumPtrToTypeString(historyItem.GetTargetEnvironment()),
	}
}

func MapIdentitySetRemoteStateToTerraform(identitySet graphmodels.IdentitySetable) sharedmodels.IdentitySetResourceModel {
	if identitySet == nil {
		return sharedmodels.IdentitySetResourceModel{}
	}

	return sharedmodels.IdentitySetResourceModel{
		Application: MapIdentityRemoteStateToTerraform(identitySet.GetApplication()),
		User:        MapIdentityRemoteStateToTerraform(identitySet.GetUser()),
		Device:      MapIdentityRemoteStateToTerraform(identitySet.GetDevice()),
		// TODO - not in SDK
		// Encrypted:                MapIdentityRemoteStateToTerraform(identitySet.GetEncrypted()),
		// OnPremises:               MapIdentityRemoteStateToTerraform(identitySet.GetOnPremises()),
		// Guest:                    MapIdentityRemoteStateToTerraform(identitySet.GetGuest()),
		// Phone:                    MapIdentityRemoteStateToTerraform(identitySet.GetPhone()),
		// ApplicationInstance:      MapIdentityRemoteStateToTerraform(identitySet.GetApplicationInstance()),
		// Conversation:             MapIdentityRemoteStateToTerraform(identitySet.GetConversation()),
		// ConversationIdentityType: MapIdentityRemoteStateToTerraform(identitySet.GetConversation()),
	}
}

func MapIdentityRemoteStateToTerraform(identity graphmodels.Identityable) sharedmodels.IdentityResourceModel {
	if identity == nil {
		return sharedmodels.IdentityResourceModel{}
	}

	return sharedmodels.IdentityResourceModel{
		DisplayName: types.StringPointerValue(identity.GetDisplayName()),
		ID:          types.StringPointerValue(identity.GetId()),
		// TODO - field missing from SDK
		//TenantID:    types.StringValue(state.StringPtrToString(identity.GetTenantId())),
	}
}
