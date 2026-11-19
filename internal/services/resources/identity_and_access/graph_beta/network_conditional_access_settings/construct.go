package graphBetaNetworkConditionalAccessSettings

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func constructResource(value types.String) (models.ConditionalAccessSettingsable, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, errMissingSignalingStatus
	}
	body := models.NewConditionalAccessSettings()
	if err := convert.FrameworkToGraphEnum(
		value,
		models.ParseStatus,
		body.SetSignalingStatus,
	); err != nil {
		return nil, fmt.Errorf("construct signaling status: %w", err)
	}
	return body, nil
}
