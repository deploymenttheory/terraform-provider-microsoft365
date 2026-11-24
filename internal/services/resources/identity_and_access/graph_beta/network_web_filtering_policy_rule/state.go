package graphBetaNetworkWebFilteringPolicyRule

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// customHeaderObjectType is the Terraform object type of one custom_headers element.
var customHeaderObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
}}

func MapRemoteStateToTerraform(ctx context.Context, data *NetworkWebFilteringPolicyRuleResourceModel, remoteResource *webFilteringRule) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	// HTTP methods and session types are comma-separated strings; Terraform keeps them as lower-case sets.
	splitCommaValues := func(value *string) []string {
		if value == nil || strings.TrimSpace(*value) == "" {
			return nil
		}
		parts := strings.Split(*value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, strings.ToLower(trimmed))
			}
		}
		return result
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.ID)
	data.Name = convert.GraphToFrameworkString(remoteResource.Name)
	data.Description = convert.GraphToFrameworkString(remoteResource.Description)
	data.Priority = convert.GraphToFrameworkInt64(remoteResource.Priority)
	data.Status = convert.GraphToFrameworkString(remoteResource.Settings.Status)

	// The action is identified only by its OData type. An unknown type is kept as Graph returned it. Graph rejects
	// writes of any header modification other than headerModificationAdd, so other kinds are skipped rather than
	// mapped to an add. Headers keep Graph order so state matches the configured list.
	data.Action = types.StringNull()
	data.CustomHeaders = types.ListNull(customHeaderObjectType)
	if action := remoteResource.Action; action != nil {
		switch value := strings.ToLower(action.ODataType); {
		case action.ODataType == "":
		case strings.HasSuffix(value, "webfilteringactionallow"):
			data.Action = types.StringValue("allow")
		case strings.HasSuffix(value, "webfilteringactionblock"):
			data.Action = types.StringValue("block")
		default:
			data.Action = types.StringValue(action.ODataType)
		}
		if action.HeaderSettings != nil {
			headers := make([]attr.Value, 0, len(action.HeaderSettings.Modifications))
			for _, modification := range action.HeaderSettings.Modifications {
				if modification.ODataType != headerModificationAddODataType {
					continue
				}
				headers = append(headers, types.ObjectValueMust(customHeaderObjectType.AttrTypes, map[string]attr.Value{
					"name":  types.StringValue(modification.HeaderName),
					"value": types.StringValue(modification.HeaderValue),
				}))
			}
			if len(headers) > 0 {
				data.CustomHeaders = types.ListValueMust(customHeaderObjectType, headers)
			}
		}
	}

	// Graph returns URL/FQDN destinations as a values array under
	// webFilteringUrlDestination. Keep Terraform state in the same API-oriented
	// set(string) shape even though the portal renders those values in one
	// comma-delimited text box.
	var urlsOrFqdns, webCategories, httpMethods, sessionTypes []string
	if destinations := remoteResource.MatchingConditions.Destinations; destinations != nil {
		for _, target := range destinations.Targets {
			switch target.ODataType {
			case webFilteringURLDestinationODataType:
				urlsOrFqdns = append(urlsOrFqdns, target.Values...)
			case webFilteringWebCategoryDestinationODataType:
				webCategories = append(webCategories, target.Values...)
			}
		}
		httpMethods = splitCommaValues(destinations.HTTPRequestMethod)
	}
	if sources := remoteResource.MatchingConditions.Sources; sources != nil {
		sessionTypes = splitCommaValues(sources.SessionType)
	}
	data.UrlsOrFqdns = convert.GraphToFrameworkStringSet(ctx, urlsOrFqdns)
	data.WebCategories = convert.GraphToFrameworkStringSet(ctx, webCategories)
	data.HTTPMethods = convert.GraphToFrameworkStringSet(ctx, httpMethods)
	data.SessionTypes = convert.GraphToFrameworkStringSet(ctx, sessionTypes)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
