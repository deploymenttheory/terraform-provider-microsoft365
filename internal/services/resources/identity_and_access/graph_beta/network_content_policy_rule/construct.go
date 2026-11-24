package graphBetaNetworkContentPolicyRule

import (
	"context"
	"errors"
	"fmt"
	"strings"

	commonattr "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/attr"
)

var (
	errReadDestinations           = errors.New("failed to read destinations")
	errUnsupportedDestinationType = errors.New("unsupported destination type")
)

const (
	fileRuleODataType                         = "#microsoft.graph.networkaccess.fileRule"
	filePolicyWebCategoryDestinationODataType = "#microsoft.graph.networkaccess.filePolicyWebCategoryDestination"
	filePolicyFQDNDestinationODataType        = "#microsoft.graph.networkaccess.filePolicyFqdnDestination"
	filePolicyURLDestinationODataType         = "#microsoft.graph.networkaccess.filePolicyUrlDestination"
)

// constructResource builds the complete file rule body observed from the portal. Create and Update both send it.
// Activities and session types are sent as comma-separated strings, content types always as arrays (empty when unset),
// and null values are left out of the body.
func constructResource(ctx context.Context, data *NetworkContentPolicyRuleResourceModel) (map[string]any, error) {
	settings := map[string]any{}
	if !data.Status.IsNull() {
		settings["status"] = data.Status.ValueString()
	}
	body := map[string]any{"@odata.type": fileRuleODataType, "settings": settings}
	if !data.Name.IsNull() {
		body["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() {
		body["description"] = data.Description.ValueString()
	}
	if !data.Action.IsNull() {
		body["action"] = data.Action.ValueString()
	}
	if !data.Priority.IsNull() {
		body["priority"] = data.Priority.ValueInt64()
	}

	fileAttributes := map[string]any{"contentTypes": []string{}, "textContentTypes": []string{}}
	if activities := commonattr.StringSetElements(data.Activities); len(activities) > 0 {
		fileAttributes["activities"] = strings.Join(activities, ",")
	}
	if contentTypes := commonattr.StringSetElements(data.ContentTypes); contentTypes != nil {
		fileAttributes["contentTypes"] = contentTypes
	}
	if textContentTypes := commonattr.StringSetElements(data.TextContentTypes); textContentTypes != nil {
		fileAttributes["textContentTypes"] = textContentTypes
	}
	sources := map[string]any{}
	if sessionTypes := commonattr.StringSetElements(data.SessionTypes); len(sessionTypes) > 0 {
		sources["sessionType"] = strings.Join(sessionTypes, ",")
	}
	conditions := map[string]any{"fileAttributes": fileAttributes, "sources": sources}

	if !data.Destinations.IsNull() && !data.Destinations.IsUnknown() {
		var models []ContentPolicyRuleDestinationModel
		if diags := data.Destinations.ElementsAs(ctx, &models, false); diags.HasError() {
			return nil, fmt.Errorf("%w: %s", errReadDestinations, diags.Errors()[0].Detail())
		}
		destinations := make([]map[string]any, 0, len(models))
		for _, model := range models {
			destination := map[string]any{}
			switch model.Type.ValueString() {
			case destinationTypeWebCategory:
				destination["@odata.type"] = filePolicyWebCategoryDestinationODataType
			case destinationTypeFQDN:
				destination["@odata.type"] = filePolicyFQDNDestinationODataType
			case destinationTypeURL:
				destination["@odata.type"] = filePolicyURLDestinationODataType
			default:
				return nil, fmt.Errorf("%w %q", errUnsupportedDestinationType, model.Type.ValueString())
			}
			if values := commonattr.StringSetElements(model.Values); values != nil {
				destination["values"] = values
			}
			destinations = append(destinations, destination)
		}
		conditions["destinations"] = destinations
	}
	body["matchingConditions"] = conditions
	return body, nil
}
