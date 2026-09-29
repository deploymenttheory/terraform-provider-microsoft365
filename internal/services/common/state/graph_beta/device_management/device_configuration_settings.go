package devicemanagement

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	normalizeImportedCertificateKeyUsages(response)
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
		var settings map[string]any
		if err := json.Unmarshal([]byte(prior.ValueString()), &settings); err != nil {
			return types.StringNull(), fmt.Errorf("read configured device configuration settings: %w", err)
		}
		normalizeImportedCertificateKeyUsages(settings)
		content, err := json.Marshal(settings)
		if err != nil {
			return types.StringNull(), fmt.Errorf("serialize configured device configuration settings: %w", err)
		}
		configured, err := normalize.JSONAlphabetically(string(content))
		if err != nil {
			return types.StringNull(), fmt.Errorf(
				"normalize configured device configuration settings: %w",
				err,
			)
		}
		// Keep the configured representation only when the complete settings are
		// semantically equivalent. A known string must not be rewritten during apply.
		if configured == normalized {
			return prior, nil
		}
	}
	return types.StringValue(normalized), nil
}

func normalizeImportedCertificateKeyUsages(settings map[string]any) {
	switch settings["@odata.type"] {
	case "#microsoft.graph.androidDeviceOwnerImportedPFXCertificateProfile",
		"#microsoft.graph.androidForWorkImportedPFXCertificateProfile":
		// Graph appends Any Purpose on every write, even when it is already
		// present. Identical key usages have the same certificate semantics.
		usages, ok := settings["extendedKeyUsages"].([]any)
		if !ok {
			return
		}
		unique := make([]any, 0, len(usages))
		for _, usage := range usages {
			duplicate := false
			for _, existing := range unique {
				if reflect.DeepEqual(usage, existing) {
					duplicate = true
					break
				}
			}
			if !duplicate {
				unique = append(unique, usage)
			}
		}
		settings["extendedKeyUsages"] = unique
	}
}
