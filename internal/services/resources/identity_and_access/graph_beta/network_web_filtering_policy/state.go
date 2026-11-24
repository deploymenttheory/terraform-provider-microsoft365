package graphBetaNetworkWebFilteringPolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MapRemoteStateToTerraform(ctx context.Context, data *NetworkWebFilteringPolicyResourceModel, remoteResource *webFilteringPolicy) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.ID)
	data.Name = convert.GraphToFrameworkString(remoteResource.Name)
	data.Description = convert.GraphToFrameworkString(remoteResource.Description)

	// The default action is identified only by its OData type. An unknown type is kept as Graph returned it.
	data.DefaultAction = types.StringNull()
	if remoteResource.Settings != nil && remoteResource.Settings.DefaultAction != nil && remoteResource.Settings.DefaultAction.ODataType != "" {
		odataType := remoteResource.Settings.DefaultAction.ODataType
		switch value := strings.ToLower(odataType); {
		case strings.HasSuffix(value, "webfilteringactionallow"):
			data.DefaultAction = types.StringValue("allow")
		case strings.HasSuffix(value, "webfilteringactionblock"):
			data.DefaultAction = types.StringValue("block")
		default:
			data.DefaultAction = types.StringValue(odataType)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
