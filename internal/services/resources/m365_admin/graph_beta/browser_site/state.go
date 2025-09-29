package graphBetaBrowserSite

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *BrowserSiteResourceModel, remoteResource graphmodels.BrowserSiteable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.AllowRedirect = convert.GraphToFrameworkBool(remoteResource.GetAllowRedirect())
	data.Comment = convert.GraphToFrameworkString(remoteResource.GetComment())
	data.CompatibilityMode = convert.GraphToFrameworkEnum(remoteResource.GetCompatibilityMode())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.DeletedDateTime = convert.GraphToFrameworkTime(remoteResource.GetDeletedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.MergeType = convert.GraphToFrameworkEnum(remoteResource.GetMergeType())
	data.Status = convert.GraphToFrameworkEnum(remoteResource.GetStatus())
	data.TargetEnvironment = convert.GraphToFrameworkEnum(remoteResource.GetTargetEnvironment())
	data.WebUrl = convert.GraphToFrameworkString(remoteResource.GetWebUrl())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
