package graphBetaGroupPolicyMultiTextValue

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// mapRemoteStateToTerraform maps the Graph API response to the Terraform resource model
func mapRemoteStateToTerraform(ctx context.Context, data *GroupPolicyMultiTextValueResourceModel, remoteResource interface {
	GetId() *string
	GetValues() []string
	GetCreatedDateTime() *time.Time
	GetLastModifiedDateTime() *time.Time
}) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform for %s", ResourceName))

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())

	// Map the values array
	values := remoteResource.GetValues()
	if values != nil {
		listValue, diags := types.ListValueFrom(ctx, types.StringType, values)
		if !diags.HasError() {
			data.Values = listValue
		} else {
			tflog.Error(ctx, "Failed to convert values array to Terraform list")
			data.Values = types.ListNull(types.StringType)
		}
	} else {
		data.Values = types.ListNull(types.StringType)
	}

	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform for %s", ResourceName))
}
