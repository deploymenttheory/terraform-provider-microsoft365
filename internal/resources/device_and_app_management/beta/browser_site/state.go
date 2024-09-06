package graphbetabrowsersite

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/models/graph_beta"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *BrowserSiteResourceModel, remoteResource models.BrowserSiteable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.AllowRedirect = types.BoolValue(state.BoolPtrToBool(remoteResource.GetAllowRedirect()))
	data.Comment = types.StringValue(state.StringPtrToString(remoteResource.GetComment()))
	data.CompatibilityMode = state.EnumPtrToTypeString(remoteResource.GetCompatibilityMode())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.DeletedDateTime = state.TimeToString(remoteResource.GetDeletedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.MergeType = state.EnumPtrToTypeString(remoteResource.GetMergeType())
	data.Status = state.EnumPtrToTypeString(remoteResource.GetStatus())
	data.TargetEnvironment = state.EnumPtrToTypeString(remoteResource.GetTargetEnvironment())
	data.WebUrl = types.StringValue(state.StringPtrToString(remoteResource.GetWebUrl()))

	// Handle History
	history := remoteResource.GetHistory()
	if len(history) == 0 {
		data.History = []BrowserSiteHistoryModel{}
	} else {
		data.History = make([]BrowserSiteHistoryModel, len(history))
		for i, historyItem := range history {
			data.History[i] = MapHistoryRemoteStateToTerraform(historyItem)
		}
	}

	// Handle LastModifiedBy
	if lastModifiedBy := remoteResource.GetLastModifiedBy(); lastModifiedBy != nil {
		data.LastModifiedBy = MapIdentitySetRemoteStateToTerraform(lastModifiedBy)
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}

func MapHistoryRemoteStateToTerraform(historyItem models.BrowserSiteHistoryable) BrowserSiteHistoryModel {
	return BrowserSiteHistoryModel{
		AllowRedirect:     types.BoolValue(state.BoolPtrToBool(historyItem.GetAllowRedirect())),
		Comment:           types.StringValue(state.StringPtrToString(historyItem.GetComment())),
		CompatibilityMode: state.EnumPtrToTypeString(historyItem.GetCompatibilityMode()),
		LastModifiedBy:    MapIdentitySetRemoteStateToTerraform(historyItem.GetLastModifiedBy()),
		MergeType:         state.EnumPtrToTypeString(historyItem.GetMergeType()),
		PublishedDateTime: state.TimeToString(historyItem.GetPublishedDateTime()),
		TargetEnvironment: state.EnumPtrToTypeString(historyItem.GetTargetEnvironment()),
	}
}

func MapIdentitySetRemoteStateToTerraform(identitySet models.IdentitySetable) sharedmodels.IdentitySetModel {
	if identitySet == nil {
		return sharedmodels.IdentitySetModel{}
	}

	return sharedmodels.IdentitySetModel{
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

func MapIdentityRemoteStateToTerraform(identity models.Identityable) sharedmodels.IdentityModel {
	if identity == nil {
		return sharedmodels.IdentityModel{}
	}

	return sharedmodels.IdentityModel{
		DisplayName: types.StringValue(state.StringPtrToString(identity.GetDisplayName())),
		ID:          types.StringValue(state.StringPtrToString(identity.GetId())),
		// TODO - field missing from SDK
		//TenantID:    types.StringValue(state.StringPtrToString(identity.GetTenantId())),
	}
}
