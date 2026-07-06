package graphBetaNetworkWebFilteringPolicyRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MapRemoteStateToTerraform(ctx context.Context, data *NetworkWebFilteringPolicyRuleResourceModel, remoteResource *webFilteringPolicyRuleResponse) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.id)
	data.Name = convert.GraphToFrameworkString(remoteResource.name)
	data.Description = convert.GraphToFrameworkString(remoteResource.description)
	data.Priority = convert.GraphToFrameworkInt64(remoteResource.priority)
	data.Action = convert.GraphToFrameworkString(remoteResource.action)
	data.Status = convert.GraphToFrameworkString(remoteResource.status)
	// Graph returns URL/FQDN destinations as a values array under
	// webFilteringUrlDestination. Keep Terraform state in the same API-oriented
	// set(string) shape even though the portal renders those values in one
	// comma-delimited text box.
	data.UrlsOrFqdns = convert.GraphToFrameworkStringSet(ctx, remoteResource.urlsOrFqdns)
	data.WebCategories = convert.GraphToFrameworkStringSet(ctx, remoteResource.webCategories)
	data.HTTPMethods = convert.GraphToFrameworkStringSet(ctx, remoteResource.httpMethods)
	data.SessionTypes = convert.GraphToFrameworkStringSet(ctx, remoteResource.sessionTypes)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
