package graphBetaBrowserSite

import (
	"context"
	"fmt"

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
	data.AllowRedirect = state.BoolPointerValue(remoteResource.GetAllowRedirect())
	data.Comment = types.StringPointerValue(remoteResource.GetComment())
	data.CompatibilityMode = state.EnumPtrToTypeString(remoteResource.GetCompatibilityMode())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.DeletedDateTime = state.TimeToString(remoteResource.GetDeletedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.MergeType = state.EnumPtrToTypeString(remoteResource.GetMergeType())
	data.Status = state.EnumPtrToTypeString(remoteResource.GetStatus())
	data.TargetEnvironment = state.EnumPtrToTypeString(remoteResource.GetTargetEnvironment())
	data.WebUrl = types.StringPointerValue(remoteResource.GetWebUrl())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
