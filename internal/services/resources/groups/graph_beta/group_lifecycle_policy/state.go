package graphBetaGroupLifecyclePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of a GroupLifecyclePolicy resource to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *GroupLifecyclePolicyResourceModel, remoteResource graphmodels.GroupLifecyclePolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.AlternateNotificationEmails = convert.GraphToFrameworkString(remoteResource.GetAlternateNotificationEmails())
	data.GroupLifetimeInDays = convert.GraphToFrameworkInt32(remoteResource.GetGroupLifetimeInDays())
	data.ManagedGroupTypes = convert.GraphToFrameworkString(remoteResource.GetManagedGroupTypes())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state %s with id %s", ResourceName, data.ID.ValueString()))
}
