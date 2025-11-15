package graphBetaFilteringPolicy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructResource converts the Terraform resource model to a plain map for JSON marshaling
// Returns a map[string]any that can be directly JSON marshaled by the HTTP client
func constructResource(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, data *FilteringPolicyResourceModel) (map[string]any, error) {

	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := make(map[string]any)

	convert.FrameworkToGraphString(data.Name, func(val *string) {
		if val != nil {
			requestBody["name"] = *val
		}
	})

	convert.FrameworkToGraphString(data.Description, func(val *string) {
		if val != nil {
			requestBody["description"] = *val
		}
	})

	convert.FrameworkToGraphString(data.Action, func(val *string) {
		if val != nil {
			requestBody["action"] = *val
		}
	})

	if debugJSON, err := json.MarshalIndent(requestBody, "", "    "); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), map[string]any{
			"json": "\n" + string(debugJSON),
		})
	} else {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
