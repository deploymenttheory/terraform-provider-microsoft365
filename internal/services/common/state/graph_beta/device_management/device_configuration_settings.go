package devicemanagement

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/normalize"
)

// StateDeviceConfigurationSettings stores the complete template settings, excluding
// metadata owned by the resource schema and read-only OData response annotations.
func StateDeviceConfigurationSettings(
	prior types.String,
	response map[string]any,
) (types.String, error) {
	// Graph does not disclose Wi-Fi keys. An explicit false is authoritative when
	// the platform returns preSharedKeyIsSet; otherwise preserve the configured key.
	if key, exists := response["preSharedKey"]; exists && key == nil &&
		response["preSharedKeyIsSet"] != false &&
		!prior.IsNull() &&
		!prior.IsUnknown() {
		var configured map[string]any
		if err := json.Unmarshal([]byte(prior.ValueString()), &configured); err != nil {
			return types.StringNull(), fmt.Errorf("read configured Wi-Fi key: %w", err)
		}
		if value, ok := configured["preSharedKey"].(string); ok {
			response["preSharedKey"] = value
		}
	}
	delete(response, "preSharedKeyIsSet")
	if response["@odata.type"] == "#microsoft.graph.windowsUpdateForBusinessConfiguration" {
		for _, field := range []string{"qualityUpdatesPauseExpiryDateTime", "featureUpdatesPauseExpiryDateTime", "qualityUpdatesRollbackStartDateTime", "featureUpdatesRollbackStartDateTime", "qualityUpdatesPauseStartDate", "featureUpdatesPauseStartDate", "qualityUpdatesWillBeRolledBack", "featureUpdatesWillBeRolledBack"} {
			delete(response, field)
		}
	}
	for _, field := range []string{"id", "displayName", "description", "roleScopeTagIds", "createdDateTime", "lastModifiedDateTime", "version", "supportsScopeTags", "assignments"} {
		delete(response, field)
	}
	for field := range response {
		if strings.HasPrefix(field, "@odata.") && field != "@odata.type" {
			delete(response, field)
		}
	}
	content, err := json.Marshal(response)
	if err != nil {
		return types.StringNull(), fmt.Errorf("serialize device configuration settings: %w", err)
	}
	normalized, err := normalize.JSONAlphabetically(string(content))
	if err != nil {
		return types.StringNull(), fmt.Errorf("normalize device configuration settings: %w", err)
	}
	if !prior.IsNull() && !prior.IsUnknown() {
		configured, err := normalize.JSONAlphabetically(prior.ValueString())
		if err != nil {
			return types.StringNull(), fmt.Errorf(
				"normalize configured device configuration settings: %w",
				err,
			)
		}
		// Keep the configured representation when only formatting differs. A known
		// configured string must not be rewritten during apply.
		if configured == normalized {
			return prior, nil
		}
	}
	return types.StringValue(normalized), nil
}
