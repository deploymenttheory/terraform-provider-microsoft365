package graphBetaNetworkTLSInspectionPolicy

import (
	"context"
	"fmt"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(
	_ context.Context,
	data *NetworkTLSInspectionPolicyResourceModel,
	remote models.TlsInspectionPolicyable,
) error {
	if remote == nil || remote.GetId() == nil || *remote.GetId() == "" || remote.GetName() == nil ||
		remote.GetSettings() == nil || remote.GetVersion() == nil || remote.GetLastModifiedDateTime() == nil {
		return fmt.Errorf("%w: missing required fields", errInvalidResponse)
	}
	defaultAction, ok := additionalString(remote.GetSettings().GetAdditionalData(), "defaultAction")
	if !ok {
		return fmt.Errorf("%w: settings.defaultAction is missing or invalid", errInvalidResponse)
	}
	data.ID = convert.GraphToFrameworkString(remote.GetId())
	data.Name = convert.GraphToFrameworkString(remote.GetName())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.DefaultAction = convert.GraphToFrameworkString(defaultAction)
	data.Version = convert.GraphToFrameworkString(remote.GetVersion())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remote.GetLastModifiedDateTime())
	return nil
}

func additionalString(values map[string]any, key string) (*string, bool) {
	if values == nil {
		return nil, false
	}
	switch value := values[key].(type) {
	case string:
		return &value, true
	case *string:
		return value, value != nil
	default:
		return nil, false
	}
}
