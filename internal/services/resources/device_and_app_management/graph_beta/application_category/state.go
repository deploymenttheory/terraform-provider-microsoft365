package graphBetaApplicationCategory

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of an ApplicationCategoryResourceModel to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *ApplicationCategoryResourceModel, remoteResource graphmodels.MobileAppCategoryable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": types.StringPointerValue(remoteResource.GetId()),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())

	// Handle last modified date time
	if lastModified := remoteResource.GetLastModifiedDateTime(); lastModified != nil {
		data.LastModifiedDateTime = types.StringValue(lastModified.Format(time.RFC3339))
	} else {
		data.LastModifiedDateTime = types.StringNull()
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
