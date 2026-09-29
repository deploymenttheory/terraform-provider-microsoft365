package sharedConstructors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
)

var errInvalidDeviceConfigurationSettings = errors.New("invalid device configuration settings")

// ConstructDeviceConfigurationSettings preserves the complete template JSON in an SDK request.
// AdditionalData preserves explicit nulls and derived-type fields without requiring an HCL
// schema or a provider-maintained type switch for every device configuration template.
func ConstructDeviceConfigurationSettings(
	settings types.String,
) (graphmodels.DeviceConfigurationable, error) {
	var content map[string]any
	if err := json.Unmarshal([]byte(settings.ValueString()), &content); err != nil {
		return nil, fmt.Errorf("decode device configuration settings: %w", err)
	}
	if content == nil {
		return nil, fmt.Errorf(
			"%w: settings must be a JSON object",
			errInvalidDeviceConfigurationSettings,
		)
	}
	odataType, ok := content["@odata.type"].(string)
	if !ok || !strings.HasPrefix(odataType, "#microsoft.graph.") ||
		len(odataType) == len("#microsoft.graph.") {
		return nil, fmt.Errorf(
			"%w: include the root @odata.type, for example #microsoft.graph.iosGeneralDeviceConfiguration",
			errInvalidDeviceConfigurationSettings,
		)
	}
	for _, field := range []string{"displayName", "description", "roleScopeTagIds", "id", "createdDateTime", "lastModifiedDateTime", "version", "supportsScopeTags", "assignments", "@odata.context", "@odata.etag"} {
		if _, exists := content[field]; exists {
			return nil, fmt.Errorf(
				"%w: %s belongs to the resource metadata or assignments, not settings",
				errInvalidDeviceConfigurationSettings,
				field,
			)
		}
	}
	// Graph's XML OMA field is Edm.Binary on write, but its plaintext endpoint
	// returns XML text. Keep configuration in cleartext and encode only this typed
	// field at the request boundary; ordinary strings and binary values are unchanged.
	if odataType == "#microsoft.graph.windows10CustomConfiguration" {
		entries, _ := content["omaSettings"].([]any)
		for _, entry := range entries {
			setting, ok := entry.(map[string]any)
			if !ok || setting["@odata.type"] != "#microsoft.graph.omaSettingStringXml" {
				continue
			}
			if value, ok := setting["value"].(string); ok {
				setting["value"] = helpers.ByteStringToBase64([]byte(value))
			}
		}
	}
	requestBody := graphmodels.NewDeviceConfiguration()
	requestBody.SetOdataType(&odataType)
	delete(content, "@odata.type")
	requestBody.SetAdditionalData(content)
	return requestBody, nil
}
