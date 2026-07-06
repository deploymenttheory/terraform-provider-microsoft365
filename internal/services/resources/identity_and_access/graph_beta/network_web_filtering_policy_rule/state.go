package graphBetaNetworkWebFilteringPolicyRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	data.CustomHeaders = graphCustomHeadersToFramework(ctx, remoteResource.customHeaders)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

func graphCustomHeadersToFramework(ctx context.Context, headers []customHeaderResponse) types.List {
	objectType := customHeaderObjectType()
	if len(headers) == 0 {
		return types.ListNull(objectType)
	}

	values := make([]customHeaderModel, 0, len(headers))
	for _, header := range headers {
		values = append(values, customHeaderModel{
			HeaderName:  convert.GraphToFrameworkString(header.headerName),
			HeaderValue: convert.GraphToFrameworkString(header.headerValue),
		})
	}

	list, diags := types.ListValueFrom(ctx, objectType, values)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert custom headers to types.List", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.ListNull(objectType)
	}

	return list
}

func customHeaderObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"header_name":  types.StringType,
			"header_value": types.StringType,
		},
	}
}
