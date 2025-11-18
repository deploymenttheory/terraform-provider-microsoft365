package graphBetaNamedLocation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructResourceForDeletion builds the PATCH request body for preparing a trusted IP named location for deletion.
// Microsoft Graph API requires trusted IP locations to be set to untrusted before deletion.
// The request must include displayName, isTrusted=false, and at least one IP address.
func constructResourceForDeletion(ctx context.Context) (map[string]any, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s deletion patch body", ResourceName))

	requestBody := map[string]any{
		"@odata.type": "#microsoft.graph.ipNamedLocation",
		"displayName": "for_deletion",
		"isTrusted":   false,
		"ipRanges": []map[string]any{
			{
				"@odata.type": "#microsoft.graph.iPv4CidrRange",
				"cidrAddress": "0.0.0.0/32",
			},
		},
	}

	// Debug logging using plain JSON marshal
	if debugJSON, err := json.MarshalIndent(requestBody, "", "    "); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s deletion", ResourceName), map[string]any{
			"json": "\n" + string(debugJSON),
		})
	} else {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s deletion patch body", ResourceName))

	return requestBody, nil
}
